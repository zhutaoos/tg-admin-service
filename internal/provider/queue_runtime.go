package provider

import (
	botregistry "app/internal/provider/botregistry"
	telegram "app/internal/provider/telegram"
	"app/internal/queue"
	"app/tools/logger"
	"context"

	"go.uber.org/fx"
)

// AsTelegramProvider 适配 telegram.Client 为 queue.TelegramProvider
func AsTelegramProvider(c *telegram.Client) queue.TelegramProvider { return c }

// AsBotRegistry 适配 botregistry.Registry 为 queue.BotRegistry
func AsBotRegistry(r *botregistry.Registry) queue.BotRegistry { return r }

// StartQueueRunners 启动 RunnerManager，用于按 chatID 动态拉起 mover/worker
func StartQueueRunners(
	lc fx.Lifecycle,
	manager *queue.RunnerManager,
) {
	var cancel context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			runCtx, c := context.WithCancel(context.Background())
			cancel = c
			go manager.Run(runCtx)
			logger.System("队列 Runner 管理器已启动")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			manager.Shutdown()
			logger.System("队列 Runner 管理器已停止")
			return nil
		},
	})
}
