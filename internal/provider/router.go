package provider

import (
	"app/internal/config"
	admin_controller "app/internal/controller/admin"
	user_controller "app/internal/controller/user"
	"app/internal/repository"
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

// NewRouter 创建主路由Provider
func NewRouter(
	adminRoute *router.AdminRoute,
	userRoute *router.UserRoute,
	indexRoute *router.IndexRoute,
	conf *config.Config,
	service service.Service,
	repo repository.Repository,
) *router.Router {
	return router.NewRouter(adminRoute, userRoute, indexRoute, conf, service.Token(), repo.Admin())
}
