package provider

import (
    "app/internal/queue"
    botregistry "app/internal/provider/botregistry"
    telegram "app/internal/provider/telegram"
    "app/tools/logger"
    "context"

    "github.com/redis/go-redis/v9"
    "go.uber.org/fx"
)

// AsTelegramProvider 适配 telegram.Client 为 queue.TelegramProvider
func AsTelegramProvider(c *telegram.Client) queue.TelegramProvider { return c }

// AsBotRegistry 适配 botregistry.Registry 为 queue.BotRegistry
func AsBotRegistry(r *botregistry.Registry) queue.BotRegistry { return r }

// StartQueueRunners 启动每个分片的 Mover 与 Worker，并在启动前预创建消费组
func StartQueueRunners(
    lc fx.Lifecycle,
    rdb *redis.Client,
    cfg *queue.Config,
    limiter *queue.Limiter,
    tg queue.TelegramProvider,
    registry queue.BotRegistry,
) {
    var cancels []context.CancelFunc
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            // 预先确保所有分片的消费组已创建
            ensureQueueGroups(ctx, rdb, cfg)
            logger.System("启动队列分片 Runner", "shards", cfg.ShardCount)
            // 为每个分片启动一对 mover/worker
            for i := 0; i < cfg.ShardCount; i++ {
                shard := cfg.ShardName(i)
                m := queue.NewMover(rdb, cfg, shard)
                w := queue.NewWorker(rdb, cfg, limiter, tg, registry, shard)
                mctx, mcancel := context.WithCancel(context.Background())
                wctx, wcancel := context.WithCancel(context.Background())
                cancels = append(cancels, mcancel, wcancel)
                go func() { _ = m.Run(mctx) }()
                go func() { _ = w.Run(wctx) }()
                logger.System("分片 Runner 已启动", "shard", shard)
            }
            return nil
        },
        OnStop: func(ctx context.Context) error {
            for _, c := range cancels {
                if c != nil {
                    c()
                }
            }
            logger.System("队列分片 Runner 已停止")
            return nil
        },
    })
}

// ensureQueueGroups 预先为所有分片创建对应的 Stream 与消费组（幂等）
func ensureQueueGroups(ctx context.Context, rdb *redis.Client, cfg *queue.Config) {
    p := queue.NewProducer(rdb, cfg)
    for i := 0; i < cfg.ShardCount; i++ {
        shard := cfg.ShardName(i)
        _ = p.EnsureGroupFor(ctx, shard)
    }
    logger.System("预创建分片消费组完成", "shards", cfg.ShardCount)
}
