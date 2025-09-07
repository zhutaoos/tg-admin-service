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
        func() *queue.Config { return queue.DefaultConfig() },
        queue.NewLimiter,
        queue.NewProducer,
        queue.NewMover,
        queue.NewWorker,
        telegram.NewClient,
        botregistry.NewRegistry,
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
        // 启动Mover与Worker
        func(lc fx.Lifecycle, mover *queue.Mover, worker *queue.Worker) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    go func() { _ = mover.Run(context.Background()) }()
                    go func() { _ = worker.Run(context.Background()) }()
                    return nil
                },
                OnStop: func(ctx context.Context) error { return nil },
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
