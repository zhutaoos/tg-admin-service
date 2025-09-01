package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

var (
	handlers     = make(map[string]TaskHandler)
	handlersLock sync.RWMutex
	client       *asynq.Client
	server       *asynq.Server
)

// TaskHandler 任务处理器接口
type TaskHandler interface {
	Process(ctx context.Context, payload []byte) error
	TaskType() string
}

// RegisterHandler 注册任务处理器
func RegisterHandler(handler TaskHandler) {
	handlersLock.Lock()
	defer handlersLock.Unlock()
	
	if handler == nil {
		panic("handler cannot be nil")
	}
	
	taskType := handler.TaskType()
	if taskType == "" {
		panic("task type cannot be empty")
	}
	
	handlers[taskType] = handler
	log.Printf("Registered task handler: %s", taskType)
}

// GetHandler 获取任务处理器
func GetHandler(taskType string) (TaskHandler, bool) {
	handlersLock.RLock()
	defer handlersLock.RUnlock()
	
	handler, ok := handlers[taskType]
	return handler, ok
}

// StartWorker 启动任务工作进程
func StartWorker(redisAddr string, concurrency int) error {
	if concurrency <= 0 {
		concurrency = 10
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	server = asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: concurrency,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("*", processTask)

	log.Printf("Starting asynq worker with concurrency: %d", concurrency)
	return server.Start(mux)
}

// StopWorker 停止任务工作进程
func StopWorker() {
	if server != nil {
		server.Stop()
		server.Shutdown()
		log.Println("Asynq worker stopped")
	}
}

// InitClient 初始化客户端
func InitClient(redisAddr string) {
	if client != nil {
		return
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}
	client = asynq.NewClient(redisOpt)
	log.Printf("Asynq client initialized")
}

// CloseClient 关闭客户端
func CloseClient() {
	if client != nil {
		client.Close()
		client = nil
		log.Println("Asynq client closed")
	}
}

// EnqueueTask 添加任务到队列
func EnqueueTask(taskType string, payload any) (*asynq.TaskInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}

	task := asynq.NewTask(taskType, data)
	return client.Enqueue(task)
}

// ScheduleTask 计划任务
func ScheduleTask(taskType string, payload any, processAt time.Time) (*asynq.TaskInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}

	task := asynq.NewTask(taskType, data)
	return client.Enqueue(task, asynq.ProcessAt(processAt))
}

// processTask 统一任务处理函数
func processTask(ctx context.Context, task *asynq.Task) error {
	taskType := task.Type()
	
	handler, ok := GetHandler(taskType)
	if !ok {
		return fmt.Errorf("no handler registered for task type: %s", taskType)
	}

	log.Printf("Processing task: %s, payload: %s", taskType, string(task.Payload()))
	return handler.Process(ctx, task.Payload())
}
