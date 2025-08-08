package provider

import (
	"app/internal/controller"
	"app/internal/service"
)

// NewBotController 创建机器人控制器Provider
func NewBotController(botService *service.BotService) *controller.BotController {
	return controller.NewBotController(botService)
}