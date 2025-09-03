package job

import (
	"app/internal/config"
	"app/tools/logger"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
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
	config       *TaskConfig
	redisConf    *config.RedisConf
}

func NewJobService(redisConf *config.RedisConf, lc fx.Lifecycle) *JobService {
	// 从配置中读取 Redis 信息
	redisAddr := fmt.Sprintf("%s:%s", redisConf.Ip, redisConf.Port)

	taskConfig := &TaskConfig{
		RedisAddr:   redisAddr,
		Concurrency: 10, // 默认并发数
	}

	ts := &JobService{
		handlers: make(map[string]JobHandler),
		config:   taskConfig,
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

			logger.System("任务服务启动成功", "当前时间", time.Now().Format("2006-01-02 15:04:05"), "时区", time.Now().Location().String())
			logger.System("已注册的任务处理器数量", "count", len(ts.handlers))
			for taskType := range ts.handlers {
				logger.System("已注册任务类型", "taskType", taskType)
			}
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

	mux := asynq.NewServeMux()
	mux.HandleFunc("*", ts.processTask)

	logger.System("启动 asynq worker", "concurrency", concurrency, "redisAddr", redisAddr)

	// 这是阻塞调用，会一直运行直到服务停止
	err := ts.server.Start(mux)
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
}

// GetHandler 获取任务处理器
func (ts *JobService) GetHandler(taskType string) (JobHandler, bool) {
	ts.handlersLock.RLock()
	defer ts.handlersLock.RUnlock()

	handler, ok := ts.handlers[taskType]
	return handler, ok
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

	logger.System("开始处理任务", "taskType", taskType, "payload", string(task.Payload()), "开始时间", startTime.Format("2006-01-02 15:04:05"))

	err := handler.Process(ctx, task.Payload())
	duration := time.Since(startTime)

	if err != nil {
		logger.System("任务处理失败", "taskType", taskType, "error", err, "耗时", duration.String())
		return err
	}

	logger.System("任务处理成功", "taskType", taskType, "耗时", duration.String())
	return nil
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

// AddCronTask 添加周期性任务
func (ts *JobService) AddCronTask(cronExpr, taskType string, payload string) (string, error) {
    if ts.scheduler == nil {
        return "", fmt.Errorf("scheduler not initialized")
    }

    logger.System("准备注册周期任务",
        "cronExpr", cronExpr,
        "taskType", taskType,
        "当前时间", time.Now().Format("2006-01-02 15:04:05"),
        "调度器时区", "Asia/Shanghai")

    // 兼容6字段（含秒）与5字段表达式：asynq/robfig scheduler使用5字段
    // 若为6字段，自动去掉秒字段
    fields := strings.Fields(cronExpr)
    switch len(fields) {
    case 5:
        // ok
    case 6:
        normalized := strings.Join(fields[1:], " ")
        logger.System("检测到6字段Cron，自动转换为5字段", "原始", cronExpr, "转换后", normalized)
        cronExpr = normalized
    default:
        return "", fmt.Errorf("cron表达式格式错误: 期望5字段(或6字段含秒)，实际%d字段: %s", len(fields), cronExpr)
    }

    task := asynq.NewTask(taskType, []byte(payload))
    logger.System("正在注册周期任务到调度器", "cronExpr", cronExpr, "taskType", taskType)
	entryID, err := ts.scheduler.Register(cronExpr, task)
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

// ListHandlers 列出所有已注册的处理器（调试用）
func (ts *JobService) ListHandlers() map[string]JobHandler {
	ts.handlersLock.RLock()
	defer ts.handlersLock.RUnlock()

	handlers := make(map[string]JobHandler)
	for k, v := range ts.handlers {
		handlers[k] = v
	}

	logger.System("已注册的任务处理器", "count", len(handlers))
	for taskType := range handlers {
		logger.System("- 任务类型: %s", taskType)
	}

	return handlers
}

// AddFrequentTestCronTask 添加频繁执行的测试cron任务（每分钟执行，用于快速验证）
func (ts *JobService) AddFrequentTestCronTask() (string, error) {
	testPayload := fmt.Sprintf(`{"msg_type":"frequent_test","content":"频繁测试任务 - %s"}`, time.Now().Format("2006-01-02 15:04:05"))

	// 使用5字段cron格式：每分钟执行
	entryID, err := ts.AddCronTask("* * * * *", "bot_msg", testPayload)
	if err != nil {
		logger.System("添加频繁测试任务失败", "error", err)
		return "", err
	}

	logger.System("成功添加频繁测试任务", "entryID", entryID, "cronExpr", "* * * * *", "下次执行时间", "每分钟执行")
	return entryID, nil
}

// GetSchedulerEntries 获取调度器中的所有条目（调试用）
func (ts *JobService) GetSchedulerEntries() {
	if ts.scheduler == nil {
		logger.System("Scheduler未初始化")
		return
	}

	logger.System("检查Scheduler状态", "scheduler_type", fmt.Sprintf("%T", ts.scheduler))

	// 创建Redis客户端来直接检查
	redisAddr := fmt.Sprintf("%s:%s", ts.redisConf.Ip, ts.redisConf.Port)
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Username: ts.redisConf.Username,
		Password: ts.redisConf.Password,
		DB:       ts.redisConf.Db,
		PoolSize: ts.redisConf.MaxTotal,
	}

	// 使用asynq的Inspector来检查
	inspector := asynq.NewInspector(redisOpt)
	defer inspector.Close()

	// 获取调度器条目 (scheduled entries)
	entries, err := inspector.SchedulerEntries()
	if err != nil {
		logger.System("获取调度器条目失败", "error", err)
		return
	}

	logger.System("调度器条目数量", "count", len(entries))
	for i, entry := range entries {
		logger.System("调度器条目",
			"index", i,
			"id", entry.ID,
			"spec", entry.Spec,
			"task_type", entry.Task.Type(),
			"next_enqueue", entry.Next.Format("2006-01-02 15:04:05"),
			"prev_enqueue", entry.Prev.Format("2006-01-02 15:04:05"),
		)
	}
}

// TestCronValidation 测试 cron 表达式验证功能
func (ts *JobService) TestCronValidation() {
	logger.System("=== 开始测试 cron 表达式验证 ===")

	validExpressions := []string{
		"0 * * * *",   // 每小时
		"30 14 * * *", // 每天14:30
		"0 0 * * *",   // 每天0点
		"*/5 * * * *", // 每5分钟
		"0 */2 * * *", // 每2小时
		"15 10 * * 1", // 每周一10:15
		"0 12 1 * *",  // 每月1号12点
		"* * * * *",   // 每分钟
	}

	invalidExpressions := []string{
		"0 0 */1 * * *", // 6字段格式（错误）
		"* * * *",       // 4字段格式（错误）
		"* * * * * *",   // 6字段格式（错误）
		"",              // 空字符串
		"invalid",       // 无效格式
	}

	logger.System("测试有效表达式:")
	for _, expr := range validExpressions {
		fields := strings.Fields(expr)
		status := "✅ 有效"
		if len(fields) != 5 {
			status = "❌ 无效"
		}
		logger.System("验证测试", "状态", status, "表达式", expr, "字段数", len(fields))
	}

	logger.System("测试无效表达式:")
	for _, expr := range invalidExpressions {
		fields := strings.Fields(expr)
		status := "❌ 无效"
		if len(fields) == 5 {
			status = "✅ 有效"
		}
		logger.System("验证测试", "状态", status, "表达式", expr, "字段数", len(fields))
	}

	logger.System("=== cron 表达式验证测试完成 ===")
}

// AddTestHourlyTask 添加每小时测试任务，验证修复后的转换
func (ts *JobService) AddTestHourlyTask() (string, error) {
	testPayload := fmt.Sprintf(`{"msg_type":"hourly_test","content":"每小时测试任务 - %s"}`, time.Now().Format("2006-01-02 15:04:05"))

	// 使用5字段表达式：每小时执行
	cronExpr := "0 * * * *"
	logger.System("准备添加每小时测试任务", "cronExpr", cronExpr)

	entryID, err := ts.AddCronTask(cronExpr, "bot_msg", testPayload)
	if err != nil {
		logger.System("添加每小时测试任务失败", "error", err)
		return "", err
	}

	logger.System("成功添加每小时测试任务", "entryID", entryID, "cronExpr", cronExpr)
	return entryID, nil
}

// ComprehensiveTest 综合测试所有功能
func (ts *JobService) ComprehensiveTest() {
	logger.System("=== 开始综合测试 ===")

	// 1. 测试 cron 验证
	ts.TestCronValidation()

	// 2. 添加测试任务
	logger.System("添加每分钟测试任务")
	minuteEntryID, err := ts.AddFrequentTestCronTask()
	if err != nil {
		logger.System("添加每分钟测试任务失败", "error", err)
	} else {
		logger.System("每分钟测试任务添加成功", "entryID", minuteEntryID)
	}

	// 3. 添加每小时测试任务
	logger.System("添加每小时测试任务")
	hourEntryID, err := ts.AddTestHourlyTask()
	if err != nil {
		logger.System("添加每小时测试任务失败", "error", err)
	} else {
		logger.System("每小时测试任务添加成功", "entryID", hourEntryID)
	}

	// 4. 立即触发一个测试任务
	logger.System("触发立即执行测试")
	testPayload := fmt.Sprintf(`{"msg_type":"immediate_test","content":"立即执行测试 - %s"}`, time.Now().Format("2006-01-02 15:04:05"))

	taskInfo, err := ts.EnqueueTask("bot_msg", testPayload)
	if err != nil {
		logger.System("立即执行测试失败", "error", err)
	} else {
		logger.System("立即执行测试成功", "taskID", taskInfo.ID, "queue", taskInfo.Queue)
	}

	// 5. 检查调度器条目
	time.Sleep(100 * time.Millisecond) // 等待任务注册
	ts.GetSchedulerEntries()

	logger.System("=== 综合测试完成 ===")
}
