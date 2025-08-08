package router

import (
	"app/internal/controller"
	"app/tools/logger"

	"github.com/gin-gonic/gin"
)

// BotRoute 机器人管理路由
type BotRoute struct {
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
	logger.System("初始化机器人管理路由...")
	
	// 机器人管理路由组
	bot := engine.Group("/api/bot")
	{
		// 注册所有机器人相关路由
		r.botController.RegisterRoutes(bot)
	}
	
	logger.System("机器人管理路由初始化完成")
}