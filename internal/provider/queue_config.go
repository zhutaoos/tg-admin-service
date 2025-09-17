package provider

import (
	"app/internal/config"
	"app/internal/queue"
	"strconv"
	"strings"
)

// NewQueueConfig 从配置文件读取 [queue] 段，覆盖默认参数
func NewQueueConfig(conf *config.Config) *queue.Config {
	cfg := queue.DefaultConfig()

	// 使用 IsSet 区分“未设置”和“设置为0/空字符串”
	if config.IsSet(conf, "queue", "global_rate_per_sec") {
		cfg.GlobalRatePerSec = config.Get[int](conf, "queue", "global_rate_per_sec")
	}
	if config.IsSet(conf, "queue", "per_chat_min_gap_ms") {
		cfg.PerChatMinGapMs = int64(config.Get[int](conf, "queue", "per_chat_min_gap_ms"))
	}
	if config.IsSet(conf, "queue", "mover_batch") {
		cfg.MoverBatch = config.Get[int](conf, "queue", "mover_batch")
	}
	if config.IsSet(conf, "queue", "mover_interval_ms") {
		cfg.MoverIntervalMs = config.Get[int](conf, "queue", "mover_interval_ms")
	}
	if config.IsSet(conf, "queue", "horizon_sec") {
		cfg.HorizonSec = config.Get[int](conf, "queue", "horizon_sec")
	}
	if config.IsSet(conf, "queue", "stream_max_len") {
		cfg.StreamMaxLen = int64(config.Get[int](conf, "queue", "stream_max_len"))
	}
	if config.IsSet(conf, "queue", "max_chat_runners") {
		cfg.MaxChatRunners = config.Get[int](conf, "queue", "max_chat_runners")
	}
	if config.IsSet(conf, "queue", "idle_runner_ttl_ms") {
		cfg.IdleRunnerTTLMs = config.Get[int](conf, "queue", "idle_runner_ttl_ms")
	}
	if config.IsSet(conf, "queue", "candidate_cache_ttl_ms") {
		cfg.CandidateCacheTTLMs = int64(config.Get[int](conf, "queue", "candidate_cache_ttl_ms"))
	}
	if config.IsSet(conf, "queue", "failure_max_duration_ms") {
		cfg.FailureMaxDurationMs = int64(config.Get[int](conf, "queue", "failure_max_duration_ms"))
	}
	if config.IsSet(conf, "queue", "failure_backoff_ms") {
		if parsed := parseInt64List(config.Get[string](conf, "queue", "failure_backoff_ms")); len(parsed) > 0 {
			cfg.FailureBackoffMs = parsed
		}
	}
	return cfg
}

func parseInt64List(raw string) []int64 {
	parts := strings.Split(raw, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, v)
	}
	return result
}
