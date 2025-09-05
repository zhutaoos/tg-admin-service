package job

import (
    "app/internal/config"
    "app/internal/model"
    toolsCron "app/tools/cron"
    "app/tools/logger"
    "context"
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
    "sync"
    "time"

    "github.com/robfig/cron/v3"

    "github.com/hibiken/asynq"
    "go.uber.org/fx"
    "gorm.io/gorm"
)

// JobHandler 任务处理器接口
type JobHandler interface {
	Process(ctx context.Context, payload []byte) error
	TaskType() string
}

// TaskConfig 任务配置
type TaskConfig struct {
	RedisAddr   string
	Concurrency int
}

type JobService struct {
    handlers     map[string]JobHandler
    handlersLock sync.RWMutex
    client       *asynq.Client
    server       *asynq.Server
    scheduler    *asynq.Scheduler
    mux          *asynq.ServeMux
    config       *TaskConfig
    redisConf    *config.RedisConf
    db           *gorm.DB
}

func NewJobService(db *gorm.DB, redisConf *config.RedisConf, lc fx.Lifecycle) *JobService {
	// 从配置中读取 Redis 信息
	redisAddr := fmt.Sprintf("%s:%s", redisConf.Ip, redisConf.Port)

	taskConfig := &TaskConfig{
		RedisAddr:   redisAddr,
		Concurrency: 10, // 默认并发数
	}

    ts := &JobService{
        handlers: make(map[string]JobHandler),
        config:   taskConfig,
        db:       db,
    }

	// 初始化 asynq client
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Username: redisConf.Username,
		Password: redisConf.Password,
		DB:       redisConf.Db,
		PoolSize: redisConf.MaxTotal,
	}

	// 创建客户端并测试连接
	ts.client = asynq.NewClient(redisOpt)
	ts.redisConf = redisConf

	// 初始化调度器，设置时区为上海时间
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.System("无法加载Asia/Shanghai时区，使用备选方案", "error", err)
		location = time.FixedZone("Asia/Shanghai", 8*60*60) // 备选方案
	}

	// 记录时区配置信息
	now := time.Now()
	nowInLocation := now.In(location)
	logger.System("时区配置详情",
		"系统时间", now.Format("2006-01-02 15:04:05"),
		"系统时区", now.Location().String(),
		"调度器时区", location.String(),
		"调度器时区时间", nowInLocation.Format("2006-01-02 15:04:05"),
		"时区偏移", location.String())
	schedulerOpt := &asynq.SchedulerOpts{
		Location: location,
	}
	ts.scheduler = asynq.NewScheduler(redisOpt, schedulerOpt)

	// FX 生命周期管理
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 异步启动 Worker (避免阻塞)
			go func() {
				logger.System("正在启动 Worker...")
				if err := ts.StartWorker(); err != nil {
					logger.System("任务服务启动失败", "error", err)
				}
			}()

			// 等待 Worker 启动
			time.Sleep(200 * time.Millisecond)

			// 异步启动 Scheduler
			go func() {
				logger.System("正在启动调度器...")

				if err := ts.scheduler.Start(); err != nil {
					logger.System("调度器启动失败", "error", err)
				} else {
					logger.System("调度器启动成功", "当前时间", time.Now().Format("2006-01-02 15:04:05"))
				}
			}()

			// 等待调度器启动
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.System("任务服务停止中...")

			// 先停止调度器，避免新任务入队
			if ts.scheduler != nil {
				ts.scheduler.Shutdown()
				logger.System("调度器已停止")
			}

			// 再停止 Worker，处理完剩余任务
			ts.Stop()

			logger.System("任务服务已完全停止")
			return nil
		},
	})

	return ts
}

// StartWorker 启动任务工作进程
func (ts *JobService) StartWorker() error {
	concurrency := ts.config.Concurrency
	if concurrency <= 0 {
		concurrency = 10
	}

	// 从配置中读取 Redis 信息
	redisAddr := fmt.Sprintf("%s:%s", ts.redisConf.Ip, ts.redisConf.Port)

	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Username: ts.redisConf.Username,
		Password: ts.redisConf.Password,
		DB:       ts.redisConf.Db,
		PoolSize: ts.redisConf.MaxTotal,
	}

	ts.server = asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: concurrency,
	})

	// 构建ServeMux并注册当前已知的任务类型
	ts.mux = asynq.NewServeMux()
	ts.handlersLock.RLock()
	for taskType := range ts.handlers {
		ts.mux.HandleFunc(taskType, ts.processTask)
		logger.System("Worker注册任务类型", "taskType", taskType)
	}
	ts.handlersLock.RUnlock()

	logger.System("启动 asynq worker", "concurrency", concurrency, "redisAddr", redisAddr)

	// 这是阻塞调用，会一直运行直到服务停止
	err := ts.server.Start(ts.mux)
	if err != nil {
		logger.System("asynq worker 启动失败", "error", err)
		return err
	}
	return nil
}

// Stop 停止任务服务
func (ts *JobService) Stop() {
	if ts.server != nil {
		ts.server.Stop()
		ts.server.Shutdown()
		logger.System("Asynq worker stopped")
	}
	if ts.client != nil {
		ts.client.Close()
		logger.System("Asynq client closed")
	}
}

// RegisterHandler 注册任务处理器
func (ts *JobService) RegisterHandler(handler JobHandler) {
	ts.handlersLock.Lock()
	defer ts.handlersLock.Unlock()

	if handler == nil {
		panic("handler cannot be nil")
	}

	taskType := handler.TaskType()
	if taskType == "" {
		panic("task type cannot be empty")
	}

	// 检查是否重复注册
	if _, exists := ts.handlers[taskType]; exists {
		logger.System("警告: 任务类型 %s 已经注册，将覆盖原有 Handler", taskType)
	}

	ts.handlers[taskType] = handler
	logger.System("注册任务处理器成功", "taskType", taskType)

	// 若Worker已创建mux，则动态注册到ServeMux
	if ts.mux != nil {
		ts.mux.HandleFunc(taskType, ts.processTask)
		logger.System("ServeMux已动态注册任务类型", "taskType", taskType)
	}
}

// GetHandler 获取任务处理器
func (ts *JobService) GetHandler(taskType string) (JobHandler, bool) {
	ts.handlersLock.RLock()
	defer ts.handlersLock.RUnlock()

	handler, ok := ts.handlers[taskType]
	return handler, ok
}

// helper: new asynq Inspector with current redis config
func (ts *JobService) newInspector() *asynq.Inspector {
    redisAddr := fmt.Sprintf("%s:%s", ts.redisConf.Ip, ts.redisConf.Port)
    redisOpt := asynq.RedisClientOpt{
        Addr:     redisAddr,
        Username: ts.redisConf.Username,
        Password: ts.redisConf.Password,
        DB:       ts.redisConf.Db,
        PoolSize: ts.redisConf.MaxTotal,
    }
    return asynq.NewInspector(redisOpt)
}

// PurgeQueuesByDBTaskID 尝试从所有asynq队列中移除与DB任务关联的任务
// - pending/scheduled/retry/archived/completed: 直接删除
// - active: 发送取消信号（最佳努力）
// 返回删除数量与取消中的数量
func (ts *JobService) PurgeQueuesByDBTaskID(dbTaskID uint64) (int, int, error) {
    inspector := ts.newInspector()
    defer inspector.Close()
    if inspector == nil {
        return 0, 0, fmt.Errorf("inspector init failed")
    }

    removed := 0
    canceled := 0
    queue := "default"

    // helper: 删除匹配任务
    deleteMatches := func(tasks []*asynq.TaskInfo) {
        for _, ti := range tasks {
            if ti == nil || ti.Type != BotMsgType || len(ti.Payload) == 0 {
                continue
            }
            if parseTaskIDFromPayload(ti.Payload) == dbTaskID || strings.Contains(string(ti.Payload), fmt.Sprintf("任务ID：%d", dbTaskID)) {
                if err := inspector.DeleteTask(ti.Queue, ti.ID); err == nil {
                    removed++
                } else {
                    logger.Error("删除队列任务失败", "error", err, "queue", ti.Queue, "id", ti.ID)
                }
            }
        }
    }

    // helper: 分页扫描方法
    scan := func(list func(string, ...asynq.ListOption) ([]*asynq.TaskInfo, error)) {
        page := 1
        for {
            items, err := list(queue, asynq.Page(page), asynq.PageSize(100))
            if err != nil {
                logger.Error("扫描队列失败", "error", err)
                return
            }
            if len(items) == 0 {
                return
            }
            deleteMatches(items)
            if len(items) < 100 {
                return
            }
            page++
        }
    }

    // pending
    scan(inspector.ListPendingTasks)
    // scheduled
    scan(inspector.ListScheduledTasks)
    // retry
    scan(inspector.ListRetryTasks)
    // archived
    scan(inspector.ListArchivedTasks)
    // completed
    scan(inspector.ListCompletedTasks)

    // active（无法直接删除，仅发送取消信号）
    {
        page := 1
        for {
            items, err := inspector.ListActiveTasks(queue, asynq.Page(page), asynq.PageSize(100))
            if err != nil {
                logger.Error("扫描active队列失败", "error", err)
                break
            }
            if len(items) == 0 {
                break
            }
            for _, ti := range items {
                if ti == nil || ti.Type != BotMsgType || len(ti.Payload) == 0 {
                    continue
                }
                if parseTaskIDFromPayload(ti.Payload) == dbTaskID || strings.Contains(string(ti.Payload), fmt.Sprintf("任务ID：%d", dbTaskID)) {
                    if err := inspector.CancelProcessing(ti.ID); err == nil {
                        canceled++
                    } else {
                        logger.Error("取消active任务失败", "error", err, "id", ti.ID)
                    }
                }
            }
            if len(items) < 100 {
                break
            }
            page++
        }
    }

    return removed, canceled, nil
}

// DeleteScheduledByDBTaskID 删除一次性定时任务（Scheduled队列）
func (ts *JobService) DeleteScheduledByDBTaskID(dbTaskID uint64) error {
    inspector := ts.newInspector()
    defer inspector.Close()
    // 优先按固定TaskID删除（新版本使用 TaskID("schedule:<id>")）
    taskID := fmt.Sprintf("schedule:%d", dbTaskID)
    if err := inspector.DeleteTask("default", taskID); err == nil {
        logger.System("已按TaskID删除一次性定时任务", "taskID", taskID)
        return nil
    }
    // 兼容旧版本：遍历Scheduled任务，按payload中的 任务ID 匹配
    tasks, err := inspector.ListScheduledTasks("default")
    if err != nil {
        return err
    }
    needle := fmt.Sprintf("任务ID：%d", dbTaskID)
    for _, t := range tasks {
        if t.Type == BotMsgType && strings.Contains(string(t.Payload), needle) {
            if err := inspector.DeleteTask("default", t.ID); err == nil {
                logger.System("已按payload匹配删除一次性定时任务", "deleted_id", t.ID)
                return nil
            }
        }
    }
    return fmt.Errorf("未找到待删除的一次性定时任务: %s", taskID)
}

// UnregisterCronByTask 精确卸载与DB任务关联的cron条目
func (ts *JobService) UnregisterCronByTask(cronExpr string, dbTaskID uint64) (int, error) {
    if ts.scheduler == nil {
        return 0, fmt.Errorf("scheduler not initialized")
    }
    inspector := ts.newInspector()
    defer inspector.Close()
    entries, err := inspector.SchedulerEntries()
    if err != nil {
        return 0, err
    }
    match := 0
    needle := fmt.Sprintf("任务ID：%d", dbTaskID)
    for _, e := range entries {
        if e.Spec == cronExpr && e.Task != nil && e.Task.Type() == BotMsgType {
            // 通过结构化task_id或payload包含中文提示精确定位
            matched := false
            var p BotMsgPayload
            if json.Unmarshal(e.Task.Payload(), &p) == nil {
                if p.TaskID == dbTaskID {
                    matched = true
                }
            }
            if !matched && strings.Contains(string(e.Task.Payload()), needle) {
                matched = true
            }
            if matched {
                if err := ts.scheduler.Unregister(e.ID); err != nil {
                    logger.Error("卸载cron条目失败", "error", err, "entryID", e.ID)
                } else {
                    logger.System("已卸载cron条目", "entryID", e.ID, "spec", e.Spec)
                    match++
                }
            }
        }
    }
    if match == 0 {
        return 0, fmt.Errorf("未找到匹配的cron条目")
    }
    return match, nil
}

// processTask 统一任务处理函数
func (ts *JobService) processTask(ctx context.Context, task *asynq.Task) error {
    taskType := task.Type()
    startTime := time.Now()

    handler, ok := ts.GetHandler(taskType)
    if !ok {
        logger.Error("没有找到任务处理器", "taskType", taskType)
        return fmt.Errorf("no handler registered for task type: %s", taskType)
    }

    payload := task.Payload()
    logger.System("开始处理任务", "taskType", taskType, "payload", string(payload), "开始时间", startTime.Format("2006-01-02 15:04:05"))

    // 从payload解析 DB 任务ID
    dbTaskID := parseTaskIDFromPayload(payload)

    // 过期检查
    var expireAt *time.Time
    var requiresExpire bool
    var dbTask model.Task
    if t, ok := parseExpireTimeFromPayload(payload); ok {
        expireAt = t
    }
    if dbTaskID > 0 && ts.db != nil {
        // 读取任务类型与DB到期时间
        if err := ts.db.Select("trigger_type", "expire_time", "cron_expression").Where("id = ? AND is_delete = 0", dbTaskID).First(&dbTask).Error; err == nil {
            // 仅周期任务需要强制过期检查；定时执行任务不需要到期日期
            requiresExpire = dbTask.TriggerType == model.TriggerTypeCron
            if expireAt == nil && dbTask.ExpireTime != nil {
                expireAt = dbTask.ExpireTime
            }
        }
    }
    if dbTaskID > 0 {
        // 仅在需要时（cron）进行过期校验
        if requiresExpire {
            if expireAt == nil {
                ts.markExpiredAndCleanup(dbTaskID, "缺少ExpireTime")
                return nil
            }
            if !time.Now().Before(*expireAt) { // now >= expireAt
                ts.markExpiredAndCleanup(dbTaskID, "任务已到期")
                return nil
            }
        }
        // 标记执行中
        ts.updateTaskExecuting(dbTaskID)
    }

    err := handler.Process(ctx, payload)
    duration := time.Since(startTime)
    if err != nil {
        logger.System("任务处理失败", "taskType", taskType, "error", err, "耗时", duration.String())
        if dbTaskID > 0 {
            ts.updateTaskOnFailure(dbTaskID, err)
        }
        return err
    }
    logger.System("任务处理成功", "taskType", taskType, "耗时", duration.String())
    if dbTaskID > 0 {
        ts.updateTaskOnSuccess(dbTaskID)
    }
    return nil
}

// parseTaskIDFromPayload: 优先解析结构化 taskId 字段；兼容旧版 task_id；最后回退从文本中提取中文“任务ID”后的数字
func parseTaskIDFromPayload(payload []byte) uint64 {
    // 新版：taskId
    {
        var obj struct{ TaskID uint64 `json:"taskId"` }
        if err := json.Unmarshal(payload, &obj); err == nil && obj.TaskID > 0 {
            return obj.TaskID
        }
    }
    // 兼容旧版：task_id
    {
        var obj struct{ TaskID uint64 `json:"task_id"` }
        if err := json.Unmarshal(payload, &obj); err == nil && obj.TaskID > 0 {
            return obj.TaskID
        }
    }
    s := string(payload)
    re := regexp.MustCompile(`(?m)(任务ID[：:](\s*))(\d+)`)
    m := re.FindStringSubmatch(s)
    if len(m) == 4 {
        var id uint64
        _, _ = fmt.Sscanf(m[3], "%d", &id)
        return id
    }
    return 0
}

// parseExpireTimeFromPayload 解析 payload 中的 expireTime 字段
func parseExpireTimeFromPayload(payload []byte) (*time.Time, bool) {
    var obj struct{ ExpireTime string `json:"expireTime"` }
    if err := json.Unmarshal(payload, &obj); err != nil {
        return nil, false
    }
    if strings.TrimSpace(obj.ExpireTime) == "" {
        return nil, false
    }
    if t, err := time.ParseInLocation("2006-01-02 15:04:05", obj.ExpireTime, time.Local); err == nil {
        return &t, true
    }
    return nil, false
}

// 标记任务为过期失败，并做清理（cron卸载）
func (ts *JobService) markExpiredAndCleanup(taskID uint64, msg string) {
    if ts.db == nil {
        return
    }
    now := time.Now()
    // 先读取任务类型，以确定过期后的状态
    var t model.Task
    _ = ts.db.Select("trigger_type", "cron_expression").Where("id = ? AND is_delete = 0", taskID).First(&t).Error

    updates := map[string]interface{}{
        "next_execute_at":  nil,
        "last_executed_at": &now,
        "update_time":      now,
    }
    if t.TriggerType == model.TriggerTypeCron {
        // 周期任务到期：标记已完成
        updates["status"] = 2
        updates["error_message"] = ""
    } else {
        // 其他情况（如一次性任务或数据异常）：按失败处理
        updates["status"] = 3
        updates["error_message"] = msg
    }
    _ = ts.db.Model(&model.Task{}).
        Where("id = ? AND is_delete = 0", taskID).
        Updates(updates).Error

    // 周期任务到期需要卸载后续调度
    if t.TriggerType == model.TriggerTypeCron && t.CronExpression != "" {
        _, _ = ts.UnregisterCronByTask(t.CronExpression, taskID)
    }
}

func (ts *JobService) updateTaskExecuting(taskID uint64) {
    if ts.db == nil {
        return
    }
    now := time.Now()
    _ = ts.db.Model(&model.Task{}).
        Where("id = ? AND is_delete = 0", taskID).
        Updates(map[string]interface{}{
            "status":            1,
            "last_executed_at": &now,
            "update_time":      now,
        }).Error
}

func (ts *JobService) updateTaskOnSuccess(taskID uint64) {
    if ts.db == nil {
        return
    }
    var t model.Task
    if err := ts.db.Where("id = ? AND is_delete = 0", taskID).First(&t).Error; err != nil {
        return
    }
    now := time.Now()
    updates := map[string]interface{}{
        // 默认成功：一次性任务在分支中设为完成；周期任务保持执行中
        "execute_count":    t.ExecuteCount + 1,
        "retry_count":      0,
        "error_message":    "",
        "last_executed_at": &now,
        "update_time":      now,
    }
    if t.TriggerType == model.TriggerTypeSchedule {
        // 一次性任务：成功后置为完成
        updates["status"] = 2
        updates["next_execute_at"] = nil
    } else if t.TriggerType == model.TriggerTypeCron && t.CronExpression != "" {
        // 周期任务：保持执行中，计算下一次执行时间
        updates["status"] = 1
        cu := toolsCron.NewCronUtils()
        if next, err := cu.CalculateNextExecution(t.CronExpression, now); err == nil {
            updates["next_execute_at"] = next
        }
    }
    _ = ts.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

func (ts *JobService) updateTaskOnFailure(taskID uint64, execErr error) {
    if ts.db == nil {
        return
    }
    var t model.Task
    if err := ts.db.Where("id = ? AND is_delete = 0", taskID).First(&t).Error; err != nil {
        return
    }
    now := time.Now()
    newRetry := t.RetryCount + 1
    updates := map[string]interface{}{
        "status":        3,
        "retry_count":   newRetry,
        "error_message": fmt.Sprintf("%v", execErr),
        "last_executed_at": &now,
        "update_time":   now,
    }
    // 一次性定时任务：不再设置下一次执行时间
    if t.TriggerType == model.TriggerTypeSchedule {
        updates["next_execute_at"] = nil
    } else {
        // 周期任务失败：按退避计算下一次执行时间
        backoff := computeBackoff(newRetry)
        next := now.Add(backoff)
        updates["next_execute_at"] = &next
    }
    _ = ts.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

func computeBackoff(retry int) time.Duration {
    if retry <= 0 {
        return time.Minute
    }
    base := time.Minute
    // 指数退避，最大 32 分钟
    var shift uint
    if retry-1 < 0 {
        shift = 0
    } else if retry-1 > 5 {
        shift = 5
    } else {
        shift = uint(retry - 1)
    }
    d := base * time.Duration(1<<shift)
    max := time.Hour
    if d > max {
        d = max
    }
    return d
}

// EnqueueTask 添加任务到队列
func (ts *JobService) EnqueueTask(taskType string, payload string) (*asynq.TaskInfo, error) {
	if ts.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	task := asynq.NewTask(taskType, []byte(payload))
	return ts.client.Enqueue(task)
}

// ScheduleTask 计划任务
func (ts *JobService) ScheduleTask(taskType string, payload string, processAt time.Time) (*asynq.TaskInfo, error) {
    if ts.client == nil {
        return nil, fmt.Errorf("client not initialized")
    }

    task := asynq.NewTask(taskType, []byte(payload))
    return ts.client.Enqueue(task, asynq.ProcessAt(processAt))
}

// ScheduleTaskWithID 计划一次性任务并指定固定 TaskID（用于去重）
func (ts *JobService) ScheduleTaskWithID(taskType string, payload string, processAt time.Time, taskID string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
    if ts.client == nil {
        return nil, fmt.Errorf("client not initialized")
    }

    task := asynq.NewTask(taskType, []byte(payload))
    options := []asynq.Option{asynq.ProcessAt(processAt)}
    if strings.TrimSpace(taskID) != "" {
        options = append(options, asynq.TaskID(taskID))
    }
    options = append(options, opts...)
    return ts.client.Enqueue(task, options...)
}

// AddCronTask 添加周期性任务
func (ts *JobService) AddCronTask(cronExpr, taskType string, payload string, opts ...asynq.Option) (string, error) {
	if ts.scheduler == nil {
		return "", fmt.Errorf("scheduler not initialized")
	}

	// 严格只支持5字段（分钟 小时 日 月 周），不再自动兼容6字段
	fields := strings.Fields(cronExpr)
	if len(fields) != 5 {
		return "", fmt.Errorf("cron表达式格式错误: 仅支持标准5字段，实际%d字段: %s", len(fields), cronExpr)
	}

	// 使用robfig/cron进行5字段语法校验，提前拦截包含 '?' 等Quartz语法
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := p.Parse(cronExpr); err != nil {
		logger.System("cron表达式解析失败", "error", err, "cronExpr", cronExpr)
		return "", fmt.Errorf("无效的cron表达式: %v", err)
	}

    task := asynq.NewTask(taskType, []byte(payload))
    entryID, err := ts.scheduler.Register(cronExpr, task, opts...)
	if err != nil {
		logger.System("注册周期任务失败", "error", err, "cronExpr", cronExpr, "taskType", taskType)
		return "", fmt.Errorf("register periodic task failed: %w", err)
	}

	logger.System("注册周期任务成功", "cronExpr", cronExpr, "taskType", taskType, "entryID", entryID)

	// 验证 Handler 是否已注册
	if _, ok := ts.GetHandler(taskType); !ok {
		logger.System("错误: 任务类型 %s 没有对应的 Handler，cron 任务将无法执行", taskType)
		// 移除刚注册的任务
		ts.scheduler.Unregister(entryID)
		return "", fmt.Errorf("no handler registered for task type: %s", taskType)
	}

	return entryID, nil
}

// RemoveCronTask 移除周期性任务
func (ts *JobService) RemoveCronTask(entryID string) error {
	if ts.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	err := ts.scheduler.Unregister(entryID)
	if err != nil {
		return fmt.Errorf("unregister periodic task failed: %w", err)
	}

	logger.System("移除周期任务成功", "entryID", entryID)
	return nil
}

// CreateJSONPayload 创建JSON格式的payload字符串
func CreateJSONPayload(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal payload to JSON failed: %w", err)
	}
	return string(jsonData), nil
}
