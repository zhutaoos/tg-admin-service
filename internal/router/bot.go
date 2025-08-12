package router

import (
	controller "app/internal/controller/bot"

	"github.com/gin-gonic/gin"
)

// BotRoute 机器人管理路由
type BotRoute struct {
	group         *gin.RouterGroup
	botController *controller.BotController
}

// NewBotRoute 创建机器人管理路由
func NewBotRoute(
	botController *controller.BotController,
) *BotRoute {
	return &BotRoute{
		botController: botController,
	}
}

// InitRoute 初始化机器人管理路由
func (r *BotRoute) InitRoute(engine *gin.Engine) {
	// 机器人管理路由组
	r.group = engine.Group("/api/bot")
	r.group.POST("/config/create", r.botController.CreateBotConfig)
	r.group.POST("/config/update", r.botController.UpdateBotConfig)
	r.group.POST("/config/get", r.botController.GetBotConfig)
	r.group.POST("/config/search", r.botController.SearchBotConfig)
}
