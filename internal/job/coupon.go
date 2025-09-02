package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	CouponOpenJob   = "coupon:open"
	CouponExpireJob = "coupon:expire"
)

// CouponPayload 优惠券任务载荷
type CouponPayload struct {
	CouponID int64     `json:"coupon_id"`
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	Count    string    `json:"count"`
	ExpireAt time.Time `json:"expire_at"`
}

// CouponOpenHandler 优惠券开奖处理器
type CouponOpenHandler struct {
	log *zap.Logger
}

func NewCouponOpenHandler(log *zap.Logger) *CouponOpenHandler {
	return &CouponOpenHandler{
		log: log,
	}
}

func (h *CouponOpenHandler) TaskType() string {
	return CouponOpenJob
}

func (h *CouponOpenHandler) Process(ctx context.Context, payload []byte) error {
	var coupon CouponPayload
	if err := json.Unmarshal(payload, &coupon); err != nil {
		h.log.Error("优惠券开奖任务反序列化失败", zap.Error(err))
		return err
	}

	h.log.Info("开始处理优惠券开奖", zap.Any("coupon", coupon))

	// 这里实现你的开奖逻辑
	// ...

	h.log.Info("优惠券开奖处理完成", zap.Int64("couponID", coupon.CouponID))
	return nil
}

// 辅助方法：手动投递过期任务（需要在业务层调用）
func EnqueueCouponExpireTask(taskService *JobService, coupon CouponPayload, expireAt time.Time) error {
	_, err := taskService.ScheduleTask(CouponExpireJob, coupon, expireAt)
	return err
}

// CouponExpireHandler 优惠券过期处理器
type CouponExpireHandler struct {
	redis *redis.Client
	log   *zap.Logger
}

func NewCouponExpireHandler(redis *redis.Client, log *zap.Logger) *CouponExpireHandler {
	return &CouponExpireHandler{
		redis: redis,
		log:   log,
	}
}

func (h *CouponExpireHandler) TaskType() string {
	return CouponExpireJob
}

func (h *CouponExpireHandler) Process(ctx context.Context, payload []byte) error {
	var coupon CouponPayload
	if err := json.Unmarshal(payload, &coupon); err != nil {
		h.log.Error("优惠券过期任务反序列化失败", zap.Error(err))
		return err
	}

	h.log.Info("开始处理优惠券过期", zap.Any("coupon", coupon))

	// 删除相关Redis数据
	sellerCouponKey := fmt.Sprintf("seller:coupon:%d:%d", coupon.GroupID, coupon.UserID)
	if err := h.redis.ZRem(ctx, sellerCouponKey, coupon.CouponID).Err(); err != nil {
		h.log.Error("删除商家优惠券失败", zap.Error(err))
		return err
	}

	h.log.Info("优惠券过期处理完成", zap.Int64("couponID", coupon.CouponID))
	return nil
}
