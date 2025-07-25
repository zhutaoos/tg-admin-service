package provider

import "go.uber.org/fx"

// InfrastructureModule 基础设施层Module
var InfrastructureModule = fx.Options(
	fx.Provide(
		NewConfig,
		NewDatabaseConfig,
		NewRedisConfig,
		NewConfigWatcher,
		NewDatabase,
		NewRedis,
	),
)

// RepositoryModule 仓储层Module
var RepositoryModule = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

// ServiceModule 服务层Module
var ServiceModule = fx.Options(
	fx.Provide(
		NewService,
	),
)

// ControllerModule 控制器Module
var ControllerModule = fx.Options(
	fx.Provide(
		NewAdminController,
		NewUserController,
	),
)

// RouterModule 路由Module
var RouterModule = fx.Options(
	fx.Provide(
		NewAdminRoute,
		NewUserRoute,
		NewIndexRoute,
		NewRouter,
	),
)

// AllModules 所有模块
var AllModules = fx.Options(
	InfrastructureModule,
	RepositoryModule,
	ServiceModule,
	ControllerModule,
	RouterModule,
)
