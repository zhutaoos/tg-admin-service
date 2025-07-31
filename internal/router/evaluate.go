package router

import (
	"github.com/gin-gonic/gin"
)

type EvaluateRoute struct {
	group *gin.RouterGroup
}

// NewEvaluateRoute 创建首页路由
func NewEvaluateRoute() *EvaluateRoute {
	return &EvaluateRoute{}
}

// InitRoute 初始化首页路由
func (ir *EvaluateRoute) InitRoute(engine *gin.Engine) {
	ir.group = engine.Group("api/evaluate")

	//ir.group.GET("/list", controller.EvaluateController.GetList)

	// 其他首页相关路由可以在这里添加
}
