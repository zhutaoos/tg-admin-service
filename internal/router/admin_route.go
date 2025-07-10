package router

import (
	adminApi "app/internal/controller/admin"
	"app/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AdminRoute struct {
	group *gin.RouterGroup
}

func (r *AdminRoute) initRoute() {
	// 当前实现：每个路由单独配置JWT中间件
	r.group.POST("login", adminApi.AdminLogin)
	r.group.POST("getAdminList", middleware.CheckJwt(), adminApi.GetAdminList)
	r.group.POST("delAdmin", middleware.CheckJwt(), adminApi.DelAdmin)

	// 使用全局JWT中间件后的简化版本：
	// （需要在 router.go 中启用 JwtMiddlewareWithWhitelist 全局中间件）
	/*
		r.group.POST("login", adminApi.AdminLogin)           // 在白名单中，自动跳过鉴权
		r.group.POST("getAdminList", adminApi.GetAdminList)  // 自动进行JWT鉴权
		r.group.POST("delAdmin", adminApi.DelAdmin)          // 自动进行JWT鉴权
	*/
}
