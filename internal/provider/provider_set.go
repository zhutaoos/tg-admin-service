package provider

import (
	"app/internal/controller/bot"
	admin_controller "app/internal/controller/admin"
	evaluate_controller "app/internal/controller/evaluate"
	user_controller "app/internal/controller/user"

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
		// Service层Provider
		NewUserService,
		NewAdminService,
		NewTokenService,
		NewEvaluateService,
		NewBotService,
	),
)

// ControllerModule 控制器Module
var ControllerModule = fx.Options(
	fx.Provide(
		admin_controller.NewAdminController,
		user_controller.NewUserController,
		evaluate_controller.NewEvaluateController,
		controller.NewBotController,
	),
)

// RouterModule 路由Module
var RouterModule = fx.Options(
	fx.Provide(
		NewAdminRoute,
		NewUserRoute,
		NewIndexRoute,
		NewEvaluateRoute,
		NewBotRoute,
		NewRouter,
	),
)

// AllModules 所有模块
var AllModules = fx.Options(
	InfrastructureModule,
	ControllerModule,
	RouterModule,
)
