package queue

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

// Limiter 提供简化版的限流：
// - 全局：基于固定窗口计数（每秒），避免复杂Lua，在早期迭代可用
// - 每群：基于 next_allowed_ms 键控制最小间隔
type Limiter struct {
    rdb    *redis.Client
    cfg    *Config
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

// CheckPerChatGap 检查每群最小间隔：若未到时间返回需等待的毫秒数
func (l *Limiter) CheckPerChatGap(ctx context.Context, bot string, chatID int64, now time.Time) (allow bool, waitMs int64, err error) {
    key := keyChatNextAllowed(bot, chatID)
    s, err := l.rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return true, 0, nil
    }
    if err != nil {
        return true, 0, nil
    }
    var ts int64
    _, _ = fmtSscanfInt64(s, &ts)
    nowMs := now.UnixMilli()
    if ts <= nowMs {
        return true, 0, nil
    }
    return false, ts - nowMs, nil
}

// SetPerChatGap 成功发送后设置下一次允许时间
func (l *Limiter) SetPerChatGap(ctx context.Context, bot string, chatID int64, nextAllowedMs int64) {
    key := keyChatNextAllowed(bot, chatID)
    ttl := time.Duration(10) * time.Minute
    _ = l.rdb.Set(ctx, key, nextAllowedMs, ttl).Err()
}

// 辅助：fmt.Sscanf 的轻量封装，避免引入fmt在热路径
func fmtSscanfInt64(s string, p *int64) (n int, err error) {
    var x int64
    for i := 0; i < len(s); i++ {
        c := s[i]
        if c < '0' || c > '9' { break }
        x = x*10 + int64(c-'0')
        n++
    }
    *p = x
    return n, nil
}

