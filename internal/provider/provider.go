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
	botregistry "app/internal/provider/botregistry"
	telegram "app/internal/provider/telegram"
	"app/internal/queue"

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
		NewQueueConfig,
		queue.NewLimiter,
		queue.NewFailureTracker,
		queue.NewRunnerManager,
		queue.NewProducer,
		telegram.NewClient,
		botregistry.NewRegistry,
		// 适配为 queue.Worker 所需接口类型
		AsTelegramProvider,
		AsBotRegistry,
		// Worker 由 RunnerManager 在运行期按 chatID 动态创建
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
		// 启动队列 Runner 管理器（按 chatID 动态拉起 runner）
		StartQueueRunners,
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
