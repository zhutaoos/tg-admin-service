package router

import (
	"github.com/gin-gonic/gin"
)

type IndexRoute struct {
	group *gin.RouterGroup
}

// NewIndexRoute 创建首页路由
func NewIndexRoute() *IndexRoute {
	return &IndexRoute{}
}

// InitRoute 初始化首页路由
func (ir *IndexRoute) InitRoute(engine *gin.Engine) {
	ir.group = engine.Group("api/index")

	// 健康检查
	ir.group.GET("health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "service is running",
		})
	})

	// 其他首页相关路由可以在这里添加
}
