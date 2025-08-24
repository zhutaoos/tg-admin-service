package router

import (
	message_controller "app/internal/controller/message"

	"github.com/gin-gonic/gin"
)

// MessageRoute 消息路由
type MessageRoute struct {
	group             *gin.RouterGroup
	messageController *message_controller.MessageController
}

// NewMessageRoute 创建消息路由
func NewMessageRoute(
	messageController *message_controller.MessageController,
) *MessageRoute {
	return &MessageRoute{
		messageController: messageController,
	}
}

// InitRoute 初始化消息路由
func (r *MessageRoute) InitRoute(engine *gin.Engine) {
	// 消息路由组
	r.group = engine.Group("/api/message")
	r.group.POST("/create", r.messageController.CreateMessage)
	r.group.POST("/update", r.messageController.UpdateMessage)
	r.group.POST("/get", r.messageController.GetMessage)
	r.group.POST("/search", r.messageController.SearchMessages)
	r.group.POST("/delete", r.messageController.DeleteMessage)
}