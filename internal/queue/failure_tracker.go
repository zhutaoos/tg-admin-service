package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type FailureTracker struct {
	rdb *redis.Client
	cfg *Config
}

type failureState struct {
	FirstFailMs     int64 // 首次失败时间（毫秒时间戳）
	LastFailMs      int64 // 最近一次失败时间
	Attempts        int64 // 累计失败次数
	CooldownUntilMs int64 // 冷却截止时间（毫秒时间戳）
	Blocked         bool  // 是否已封禁（需人工干预或等待窗口）
}

// NewFailureTracker 初始化失败追踪器，配合配置生成退避策略。
func NewFailureTracker(rdb *redis.Client, cfg *Config) *FailureTracker {
	return &FailureTracker{rdb: rdb, cfg: cfg}
}

// Allow 判断某 bot 是否已经脱离冷却/封禁窗口。
func (ft *FailureTracker) Allow(ctx context.Context, bot string, chatID int64, now time.Time) (allowed bool, waitUntilMs int64, blocked bool, err error) {
	if ft == nil {
		return true, 0, false, nil
	}
	state, err := ft.loadState(ctx, bot, chatID)
	if err != nil {
		if err == redis.Nil {
			return true, 0, false, nil
		}
		return true, 0, false, err
	}
	if state.Blocked {
		return false, state.CooldownUntilMs, true, nil
	}
	nowMs := now.UnixMilli()
	if state.CooldownUntilMs > nowMs {
		return false, state.CooldownUntilMs, false, nil
	}
	return true, 0, false, nil
}

// ReportFailure 记录一次发送失败并返回下次可尝试的时间点。
func (ft *FailureTracker) ReportFailure(ctx context.Context, bot string, chatID int64, now time.Time) (nextRetryMs int64, blocked bool, err error) {
	if ft == nil {
		return 0, false, nil
	}
	state, err := ft.loadState(ctx, bot, chatID)
	if err != nil && err != redis.Nil {
		return 0, false, err
	}
	nowMs := now.UnixMilli()
	if state.FirstFailMs == 0 {
		state.FirstFailMs = nowMs
	}
	state.LastFailMs = nowMs
	state.Attempts++

	maxDuration := ft.cfg.FailureMaxDurationMs
	if maxDuration <= 0 {
		maxDuration = 86_400_000
	}

	backoff := ft.backoffFor(int(state.Attempts))
	if backoff <= 0 {
		backoff = 60_000
	}
	deadline := state.FirstFailMs + maxDuration
	nextRetry := nowMs + backoff
	if deadline > 0 && nextRetry > deadline {
		nextRetry = deadline
	}
	state.CooldownUntilMs = nextRetry

	seqLen := len(ft.cfg.FailureBackoffMs)
	if (deadline > 0 && nowMs >= deadline) || (seqLen > 0 && int(state.Attempts) >= seqLen) {
		state.Blocked = true
	}
	if state.Blocked {
		if deadline > 0 {
			state.CooldownUntilMs = deadline
		} else {
			state.CooldownUntilMs = nowMs
		}
	}

	if err := ft.saveState(ctx, bot, chatID, state); err != nil {
		return state.CooldownUntilMs, state.Blocked, err
	}
	return state.CooldownUntilMs, state.Blocked, nil
}

// ReportSuccess 成功发送后清除失败状态。
func (ft *FailureTracker) ReportSuccess(ctx context.Context, bot string, chatID int64) {
	if ft == nil {
		return
	}
	key := keyFailureState(bot, chatID)
	_ = ft.rdb.Del(ctx, key).Err()
}

// loadState 从 Redis 读取失败状态。
func (ft *FailureTracker) loadState(ctx context.Context, bot string, chatID int64) (failureState, error) {
	key := keyFailureState(bot, chatID)
	vals, err := ft.rdb.HMGet(ctx, key, "first_fail_ms", "last_fail_ms", "attempts", "cooldown_until_ms", "blocked").Result()
	if err != nil {
		return failureState{}, err
	}
	empty := true
	for _, v := range vals {
		if v != nil {
			empty = false
			break
		}
	}
	if empty {
		return failureState{}, redis.Nil
	}
	state := failureState{}
	if len(vals) >= 1 {
		state.FirstFailMs = toInt64(vals[0])
	}
	if len(vals) >= 2 {
		state.LastFailMs = toInt64(vals[1])
	}
	if len(vals) >= 3 {
		state.Attempts = toInt64(vals[2])
	}
	if len(vals) >= 4 {
		state.CooldownUntilMs = toInt64(vals[3])
	}
	if len(vals) >= 5 {
		state.Blocked = toBool(vals[4])
	}
	return state, nil
}

// saveState 将失败状态持久化到 Redis。
func (ft *FailureTracker) saveState(ctx context.Context, bot string, chatID int64, state failureState) error {
	key := keyFailureState(bot, chatID)
	data := map[string]any{
		"first_fail_ms":     state.FirstFailMs,
		"last_fail_ms":      state.LastFailMs,
		"attempts":          state.Attempts,
		"cooldown_until_ms": state.CooldownUntilMs,
		"blocked":           boolToInt(state.Blocked),
	}
	if err := ft.rdb.HSet(ctx, key, data).Err(); err != nil {
		return err
	}
	ttlMs := ft.cfg.FailureMaxDurationMs + 3_600_000
	if ttlMs <= 0 {
		ttlMs = 172_800_000
	}
	_ = ft.rdb.PExpire(ctx, key, time.Duration(ttlMs)*time.Millisecond).Err()
	return nil
}

// backoffFor 返回当前重试次数对应的退避毫秒数。
func (ft *FailureTracker) backoffFor(attempts int) int64 {
	if attempts <= 0 {
		return 60_000
	}
	seq := ft.cfg.FailureBackoffMs
	if len(seq) == 0 {
		return 60_000
	}
	idx := attempts - 1
	if idx >= len(seq) {
		idx = len(seq) - 1
	}
	return seq[idx]
}

// toInt64 辅助将 Redis 返回值转换为整数。
func toInt64(v any) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int64:
		return x
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	case []byte:
		n, _ := strconv.ParseInt(string(x), 10, 64)
		return n
	default:
		return 0
	}
}

// toBool 辅助将 Redis 返回值转换为布尔。
func toBool(v any) bool {
	switch x := v.(type) {
	case int:
		return x != 0
	case int64:
		return x != 0
	case string:
		return x == "1" || x == "true"
	case []byte:
		s := string(x)
		return s == "1" || s == "true"
	default:
		return false
	}
}

// boolToInt 将布尔转换为 0/1，便于写入 Redis。
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
