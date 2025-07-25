package router

import (
	user_controller "app/internal/controller/user"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	group          *gin.RouterGroup
	userController *user_controller.UserController
}

// NewUserRoute 创建用户路由
func NewUserRoute(userController *user_controller.UserController) *UserRoute {
	return &UserRoute{
		userController: userController,
	}
}

// InitRoute 初始化用户路由
func (ur *UserRoute) InitRoute(engine *gin.Engine) {
	ur.group = engine.Group("api/user")

	// 用户相关路由
	ur.group.POST("list", ur.userController.UserList)
	ur.group.GET(":id", ur.userController.GetUserInfo)
}
