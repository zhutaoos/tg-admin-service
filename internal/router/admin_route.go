package router

import (
	adminApi "app/internal/controller/admin"

	"github.com/gin-gonic/gin"
)

type AdminRoute struct {
	group *gin.RouterGroup
}

func (r *AdminRoute) initRoute() {
	// 当前实现：每个路由单独配置JWT中间件
	r.group.POST("login", adminApi.AdminLogin)
	r.group.POST("initPwd", adminApi.InitPwd)
}
