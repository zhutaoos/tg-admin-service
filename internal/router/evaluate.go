package router

import (
	evaluate_controller "app/internal/controller/evaluate"

	"github.com/gin-gonic/gin"
)

type EvaluateRoute struct {
	group              *gin.RouterGroup
	evaluateController *evaluate_controller.EvaluateController
}

// NewEvaluateRoute 创建首页路由
func NewEvaluateRoute(evaluateController *evaluate_controller.EvaluateController) *EvaluateRoute {
	return &EvaluateRoute{
		evaluateController: evaluateController,
	}
}

// InitRoute 初始化首页路由
func (ir *EvaluateRoute) InitRoute(engine *gin.Engine) {
	ir.group = engine.Group("api/evaluate")

	ir.group.POST("/list", ir.evaluateController.GetEvaluateList)

	// 其他首页相关路由可以在这里添加
}
