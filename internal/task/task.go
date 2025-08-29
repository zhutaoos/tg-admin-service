package task

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	redis  *redis.Client
	log    *zap.Logger
	client *asynq.Client
}

func NewTaskHandler(redis *redis.Client, log *zap.Logger, client *asynq.Client, registry *TaskHandlerRegistry) *CouponTaskHandler {
	handler := &TaskHandler{
		redis:  redis,
		log:    log,
		client: client,
	}

	// 注册处理器到registry
	registry.Register(func(mux *asynq.ServeMux) {
		mux.HandleFunc(CouponOpen, handler.handleCouponOpen)
		mux.HandleFunc(CouponExpire, handler.handleCouponExpire)
	})

	return handler
}

func (h *TaskHandler) EnqueueTask(taskInfo any, taskId string, processIn time.Duration) error {
	taskJSON, err := json.Marshal(taskInfo)
	if err != nil {
		h.log.Error("序列化任务失败", zap.Error(err))
		return err
	}

	task := asynq.NewTask(taskId, taskJSON)
	info, err := h.client.Enqueue(task, asynq.ProcessIn(processIn))
	if err != nil {
		h.log.Error("投递任务失败", zap.Error(err))
		return err
	}

	h.log.Info("✅ 已投递任务 taskId", zap.Any("info", info), zap.String("taskId", taskId))
	return nil
}
