package provider

import (
    "app/internal/config"
    "app/internal/queue"
)

// NewQueueConfig 从配置文件读取 [queue] 段，覆盖默认参数
func NewQueueConfig(conf *config.Config) *queue.Config {
    cfg := queue.DefaultConfig()

    // 使用 IsSet 区分“未设置”和“设置为0/空字符串”
    if config.IsSet(conf, "queue", "shard_count") {
        cfg.ShardCount = config.Get[int](conf, "queue", "shard_count")
    }
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
    if config.IsSet(conf, "queue", "worker_concurrency") {
        cfg.WorkerConcurrency = config.Get[int](conf, "queue", "worker_concurrency")
    }
    return cfg
}
