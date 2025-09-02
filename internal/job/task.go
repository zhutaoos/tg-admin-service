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
	ts.client = asynq.NewClient(redisOpt)
	ts.redisConf = redisConf

	// 初始化调度器，设置时区为上海时间
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("Asia/Shanghai", 8*60*60) // 备选方案
	}
	schedulerOpt := &asynq.SchedulerOpts{
		Location: location,
	}
	ts.scheduler = asynq.NewScheduler(redisOpt, schedulerOpt)

	// FX 生命周期管理
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 先启动 Worker
			if err := ts.StartWorker(); err != nil {
				logger.System("任务服务启动失败", "error", err)
				return err
			}

			// 等待一小段时间确保 Worker 完全启动
			time.Sleep(100 * time.Millisecond)

			// 再启动 Scheduler
			if err := ts.scheduler.Start(); err != nil {
				logger.System("调度器启动失败", "error", err)
				return err
			}

			logger.System("任务服务和调度器启动成功", "当前时间", time.Now().Format("2006-01-02 15:04:05"), "时区", time.Now().Location().String())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.System("任务服务停止")
			ts.scheduler.Shutdown()
			ts.Stop()
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

	logger.System("Starting asynq worker with concurrency: %d", concurrency)
	return ts.server.Start(mux)
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
	ts.handlers[taskType] = handler
	logger.System("Registered task handler: %s", taskType)
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

	handler, ok := ts.GetHandler(taskType)
	if !ok {
		logger.Error("no handler registered for task type", "taskType", taskType)
		return fmt.Errorf("no handler registered for task type: %s", taskType)
	}

	logger.System("Processing task: %s, payload: %s", taskType, string(task.Payload()))
	return handler.Process(ctx, task.Payload())
}

// EnqueueTask 添加任务到队列
func (ts *JobService) EnqueueTask(taskType string, payload any) (*asynq.TaskInfo, error) {
	if ts.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}

	task := asynq.NewTask(taskType, data)
	return ts.client.Enqueue(task)
}

// ScheduleTask 计划任务
func (ts *JobService) ScheduleTask(taskType string, payload any, processAt time.Time) (*asynq.TaskInfo, error) {
	if ts.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}

	task := asynq.NewTask(taskType, data)
	return ts.client.Enqueue(task, asynq.ProcessAt(processAt))
}

// AddCronTask 添加周期性任务
func (ts *JobService) AddCronTask(cronExpr, taskType string, payload any) (string, error) {
	if ts.scheduler == nil {
		return "", fmt.Errorf("scheduler not initialized")
	}

	// 转换cron表达式格式
	convertedCronExpr, err := convertCronExpr(cronExpr)
	if err != nil {
		return "", fmt.Errorf("cron表达式格式错误: %w", err)
	}

	// 分析cron表达式，计算下次执行时间
	nextRunTime, err := calculateNextRunTime(convertedCronExpr)
	if err != nil {
		logger.System("警告: 无法计算下次执行时间", "error", err)
	} else {
		logger.System("计算的下次执行时间", "nextRun", nextRunTime.Format("2006-01-02 15:04:05"))
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload failed: %w", err)
	}

	task := asynq.NewTask(taskType, data)
	entryID, err := ts.scheduler.Register(convertedCronExpr, task)
	if err != nil {
		return "", fmt.Errorf("register periodic task failed: %w", err)
	}

	logger.System("注册周期任务成功", "原始表达式", cronExpr, "转换后表达式", convertedCronExpr, "taskType", taskType, "entryID", entryID)

	// 验证 Handler 是否已注册
	if _, ok := ts.GetHandler(taskType); !ok {
		logger.System("警告: 任务类型 %s 没有对应的 Handler", taskType)
	}

	return entryID, nil
}

// convertCronExpr 转换cron表达式格式
// 输入：6字段格式 "秒 分 时 日 月 周"
// 输出：5字段格式 "分 时 日 月 周"
func convertCronExpr(cronExpr string) (string, error) {
	fields := strings.Fields(cronExpr)

	// 如果已经是5字段，直接返回
	if len(fields) == 5 {
		return cronExpr, nil
	}

	// 如果是6字段，去掉第一个字段（秒）
	if len(fields) == 6 {
		// 检查秒字段是否为0，如果不是0，给出警告
		if fields[0] != "0" {
			logger.System("警告: cron表达式中的秒字段 '%s' 将被忽略，asynq只支持分钟级精度", fields[0])
		}
		return strings.Join(fields[1:], " "), nil
	}

	return "", fmt.Errorf("不支持的cron表达式格式，期望5字段或6字段，实际%d字段: %s", len(fields), cronExpr)
}

// calculateNextRunTime 计算cron表达式的下次执行时间（简单实现用于调试）
func calculateNextRunTime(cronExpr string) (time.Time, error) {
	fields := strings.Fields(cronExpr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("invalid cron expression")
	}

	now := time.Now()
	// 这里只做简单的时间计算，主要用于调试
	// 实际的cron解析比较复杂，asynq内部会处理

	// 如果指定了具体的月份和日期
	if fields[2] != "*" && fields[3] != "*" {
		logger.System("检测到具体日期的cron表达式", "日", fields[2], "月", fields[3])
	}

	// 返回一个示例时间用于日志显示
	return now.Add(time.Minute), nil
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

// AddTestCronTask 添加测试用的cron任务（每分钟执行）
func (ts *JobService) AddTestCronTask() (string, error) {
	testPayload := map[string]string{
		"msg_type": "test",
		"content":  fmt.Sprintf("测试任务执行 - %s", time.Now().Format("2006-01-02 15:04:05")),
	}

	// 每分钟执行一次的cron表达式
	return ts.AddCronTask("* * * * *", "bot_msg", testPayload)
}

// TriggerImmediateTask 立即触发一个测试任务
func (ts *JobService) TriggerImmediateTask() (*asynq.TaskInfo, error) {
	testPayload := map[string]string{
		"msg_type": "immediate_test",
		"content":  fmt.Sprintf("立即执行测试 - %s", time.Now().Format("2006-01-02 15:04:05")),
	}
	
	logger.System("触发立即执行任务")
	return ts.EnqueueTask("bot_msg", testPayload)
}
