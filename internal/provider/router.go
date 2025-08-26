package provider

import (
	"app/internal/config"
	admin_controller "app/internal/controller/admin"
	group_controller "app/internal/controller/group"
	controller "app/internal/controller/bot"
	evaluate_controller "app/internal/controller/evaluate"
	file_controller "app/internal/controller/file"
	message_controller "app/internal/controller/message"
	task_controller "app/internal/controller/task"
	user_controller "app/internal/controller/user"
	"app/internal/router"
	"app/internal/service"
)

// NewAdminRoute 创建管理员路由Provider
func NewAdminRoute(adminController *admin_controller.AdminController) *router.AdminRoute {
	return router.NewAdminRoute(adminController)
}

// NewGroupRoute 创建群组路由Provider
func NewGroupRoute(groupController *group_controller.GroupController) *router.GroupRoute {
	return router.NewGroupRoute(groupController)
}

// NewUserRoute 创建用户路由Provider
func NewUserRoute(userController *user_controller.UserController) *router.UserRoute {
	return router.NewUserRoute(userController)
}

// NewIndexRoute 创建首页路由Provider
func NewIndexRoute() *router.IndexRoute {
	return router.NewIndexRoute()
}

// NewEvaluateRoute 创建评价路由Provider
func NewEvaluateRoute(evaluateController *evaluate_controller.EvaluateController) *router.EvaluateRoute {
	return router.NewEvaluateRoute(evaluateController)
}

// NewBotRoute 创建机器人管理路由Provider
func NewBotRoute(
	botController *controller.BotController,
) *router.BotRoute {
	return router.NewBotRoute(botController)
}

// NewMessageRoute 创建消息路由Provider
func NewMessageRoute(messageController *message_controller.MessageController) *router.MessageRoute {
	return router.NewMessageRoute(messageController)
}

// NewFileRoute 创建文件路由Provider
func NewFileRoute(fileController *file_controller.FileController) *router.FileRoute {
	return router.NewFileRoute(fileController)
}

// NewTaskRoute 创建任务路由Provider
func NewTaskRoute(taskController *task_controller.TaskController) *router.TaskRoute {
	return router.NewTaskRoute(taskController)
}

// NewRouter 创建主路由Provider
func NewRouter(
	adminRoute *router.AdminRoute,
	groupRoute *router.GroupRoute,
	userRoute *router.UserRoute,
	indexRoute *router.IndexRoute,
	evaluateRoute *router.EvaluateRoute,
	botRoute *router.BotRoute,
	messageRoute *router.MessageRoute,
	fileRoute *router.FileRoute,
	taskRoute *router.TaskRoute,
	conf *config.Config,
	tokenService service.TokenService,
	adminService service.AdminService,
) *router.Router {
	return router.NewRouter(adminRoute, groupRoute, userRoute, indexRoute, evaluateRoute, botRoute, messageRoute, fileRoute, taskRoute, conf, tokenService, adminService)
}
