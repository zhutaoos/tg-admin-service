package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter 提供简化版的限流，仅基于每个 Bot 的固定窗口计数（每秒）。
type Limiter struct {
	rdb *redis.Client
	cfg *Config
}

func NewLimiter(rdb *redis.Client, cfg *Config) *Limiter { return &Limiter{rdb: rdb, cfg: cfg} }

// TryAcquireGlobal 基于固定秒窗口计数：<=rate 允许，否则需要等待到下一个秒
func (l *Limiter) TryAcquireGlobal(ctx context.Context, bot string, now time.Time) (allow bool, waitMs int64, err error) {
	sec := now.Unix()
	key := keyBotFixedWindow(bot, sec)
	// 计数+1，并设置2秒过期
	n, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return true, 0, nil // 容错：Redis异常时默认放行，避免阻塞
	}
	if n == 1 {
		_ = l.rdb.Expire(ctx, key, 2*time.Second).Err()
	}
	if int(n) <= l.cfg.GlobalRatePerSec {
		return true, 0, nil
	}
	// 需要等到下一秒
	next := (sec + 1) * 1000
	nowMs := now.UnixMilli()
	if next <= nowMs {
		return true, 0, nil
	}
	return false, next - nowMs, nil
}
