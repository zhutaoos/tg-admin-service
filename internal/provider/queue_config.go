package provider

import (
    "app/internal/config"
    "app/internal/queue"
)

// NewQueueConfig 从配置文件读取 [queue] 段，覆盖默认参数
func NewQueueConfig(conf *config.Config) *queue.Config {
    cfg := queue.DefaultConfig()

    if v := config.Get[string](conf, "queue", "shard"); v != "" {
        cfg.Shard = v
    }
    if v := config.Get[int](conf, "queue", "global_rate_per_sec"); v > 0 {
        cfg.GlobalRatePerSec = v
    }
    if v := config.Get[int](conf, "queue", "per_chat_min_gap_ms"); v > 0 {
        cfg.PerChatMinGapMs = int64(v)
    }
    if v := config.Get[int](conf, "queue", "mover_batch"); v > 0 {
        cfg.MoverBatch = v
    }
    if v := config.Get[int](conf, "queue", "mover_interval_ms"); v > 0 {
        cfg.MoverIntervalMs = v
    }
    if v := config.Get[int](conf, "queue", "horizon_sec"); v > 0 {
        cfg.HorizonSec = v
    }
    if v := config.Get[int](conf, "queue", "stream_max_len"); v > 0 {
        cfg.StreamMaxLen = int64(v)
    }
    return cfg
}

