package provider

import (
	"app/internal/config"
	admin_controller "app/internal/controller/admin"
	evaluate_controller "app/internal/controller/evaluate"
	user_controller "app/internal/controller/user"
	"app/internal/router"
	"app/internal/service"
)

// NewAdminRoute 创建管理员路由Provider
func NewAdminRoute(adminController *admin_controller.AdminController) *router.AdminRoute {
	return router.NewAdminRoute(adminController)
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

// NewRouter 创建主路由Provider
func NewRouter(
	adminRoute *router.AdminRoute,
	userRoute *router.UserRoute,
	indexRoute *router.IndexRoute,
	evaluateRoute *router.EvaluateRoute,
	conf *config.Config,
	tokenService service.TokenService,
	adminService service.AdminService,
) *router.Router {
	return router.NewRouter(adminRoute, userRoute, indexRoute, evaluateRoute, conf, tokenService, adminService)
}
