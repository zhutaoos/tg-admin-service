package provider

import (
	"app/internal/controller/admin"
	"app/internal/controller/bot"
	"app/internal/controller/evaluate"
	"app/internal/controller/file"
	"app/internal/controller/group"
	"app/internal/controller/message"
	"app/internal/controller/task"
	"app/internal/controller/user"
	"app/internal/job"
    "app/internal/queue"
    botregistry "app/internal/provider/botregistry"
    telegram "app/internal/provider/telegram"
    "context"

	"go.uber.org/fx"
    "github.com/redis/go-redis/v9"
)

// InfrastructureModule 基础设施层Module
var InfrastructureModule = fx.Options(
    fx.Provide(
        NewConfig,
        NewDatabaseConfig,
        NewRedisConfig,
        NewConfigWatcher,
        NewDatabase,
        NewRedis,
        job.NewJobService, // 异步任务服务 (*asynqTask.TaskService)
        // 队列与发送提供者
        NewQueueConfig,
        queue.NewLimiter,
        queue.NewProducer,
        telegram.NewClient,
        botregistry.NewRegistry,
        // 适配为 queue.Worker 所需接口类型
        func(c *telegram.Client) queue.TelegramProvider { return c },
        func(r *botregistry.Registry) queue.BotRegistry { return r },
        // queue.NewWorker 不再通过 Provide 注入，改为在 Invoke 中按分片动态创建
		// Service层Provider
		NewUserService,
		NewAdminService,
		NewTokenService,
		NewEvaluateService,
		NewBotService,
		NewGroupService,
		NewMessageService,
		NewFileService,
		NewTaskService,
    ),
    fx.Invoke(
        job.NewBotMsgHandler, // 注册Bot消息处理器
        job.NewTaskRestorer,  // 启动时恢复任务
        // 启动多分片 Mover 与 Worker（按 groupID 分片）
        func(
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
                    }
                    return nil
                },
                OnStop: func(ctx context.Context) error {
                    for _, c := range cancels { if c != nil { c() } }
                    return nil
                },
            })
        },
    ),
)

// ControllerModule 控制器Module
var ControllerModule = fx.Options(
	fx.Provide(
		admin.NewAdminController,
		admin.NewGroupController,
		group.NewGroupController,
		user.NewUserController,
		evaluate.NewEvaluateController,
		bot.NewBotController,
		message.NewMessageController,
		file.NewFileController,
		task.NewTaskController,
	),
)

// RouterModule 路由Module
var RouterModule = fx.Options(
	fx.Provide(
		NewAdminRoute,
		NewGroupRoute,
		NewUserRoute,
		NewIndexRoute,
		NewEvaluateRoute,
		NewBotRoute,
		NewMessageRoute,
		NewFileRoute,
		NewTaskRoute,
		NewRouter,
	),
)

// AllModules 所有模块
var AllModules = fx.Options(
	InfrastructureModule,
	ControllerModule,
	RouterModule,
)
