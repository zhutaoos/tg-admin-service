package message

import (
	"app/internal/controller"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	controller.BaseController
	messageService *service.MessageService
}

func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// CreateMessage 创建消息
func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var req request.CreateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	if currentUserId == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户未登录"}).Response()
		return
	}
	
	err := c.messageService.CreateMessage(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建消息失败: " + err.Error()}).Response()
		return
	}
	
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建消息成功"}).Response()
}

// UpdateMessage 更新消息
func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	var req request.UpdateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	if currentUserId == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户未登录"}).Response()
		return
	}
	
	err := c.messageService.UpdateMessage(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "更新消息失败: " + err.Error()}).Response()
		return
	}
	
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新消息成功"}).Response()
}

// GetMessage 获取消息详情
func (c *MessageController) GetMessage(ctx *gin.Context) {
	var req request.GetMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	if currentUserId == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户未登录"}).Response()
		return
	}
	
	messageData, err := c.messageService.GetMessage(ctx, req.ID, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取消息失败: " + err.Error()}).Response()
		return
	}
	
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取消息成功", Data: messageData}).Response()
}

// SearchMessages 搜索消息
func (c *MessageController) SearchMessages(ctx *gin.Context) {
	var req request.SearchMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	if currentUserId == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户未登录"}).Response()
		return
	}
	
	messageData, err := c.messageService.SearchMessages(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "搜索消息失败: " + err.Error()}).Response()
		return
	}
	
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "搜索消息成功", Data: messageData}).Response()
}

// DeleteMessage 删除消息
func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	var req request.DeleteMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	if currentUserId == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户未登录"}).Response()
		return
	}
	
	err := c.messageService.DeleteMessage(ctx, req.ID, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "删除消息失败: " + err.Error()}).Response()
		return
	}
	
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "删除消息成功"}).Response()
}

// RegisterRoutes 注册路由
func (c *MessageController) RegisterRoutes(router *gin.RouterGroup) {
	message := router.Group("/message")
	{
		message.POST("/create", c.CreateMessage)
		message.POST("/update", c.UpdateMessage)
		message.POST("/get", c.GetMessage)
		message.POST("/search", c.SearchMessages)
		message.POST("/delete", c.DeleteMessage)
	}
}