package router

import (
	admin_controller "app/internal/controller/admin"

	"github.com/gin-gonic/gin"
)

type AdminRoute struct {
	group           *gin.RouterGroup
	adminController *admin_controller.AdminController
}

// NewAdminRoute 创建管理员路由
func NewAdminRoute(adminController *admin_controller.AdminController) *AdminRoute {
	return &AdminRoute{
		adminController: adminController,
	}
}

// InitRoute 初始化管理员路由
func (ar *AdminRoute) InitRoute(engine *gin.Engine) {
	ar.group = engine.Group("/admin")

	// 管理员相关路由
	ar.group.POST("login", ar.adminController.AdminLogin)
	ar.group.POST("initPwd", ar.adminController.InitPwd)
	ar.group.GET("profile", ar.adminController.Profile)
}
