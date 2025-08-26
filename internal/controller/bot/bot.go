package bot

import (
	"app/internal/controller"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

type BotController struct {
	controller.BaseController
	botService *service.BotService
}

func NewBotController(botService *service.BotService) *BotController {
	return &BotController{
		botService: botService,
	}
}

// Bot Config Related Methods

func (c *BotController) CreateBotConfig(ctx *gin.Context) {
	var req request.CreateBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	currentUserId := c.CurrentUserId(ctx)
	err := c.botService.CreateBotConfig(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建机器人配置失败: " + err.Error()}).Response()
		return
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建成功"}).Response()
}

func (c *BotController) UpdateBotConfig(ctx *gin.Context) {
	var req request.UpdateBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	
	// 临时只允许更新 BotFeature 配置
	if req.BotFeature == nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "当前只支持更新机器人功能配置"}).Response()
		return
	}
	
	currentUserId := c.CurrentUserId(ctx)
	err := c.botService.UpdateBotConfig(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "更新机器人配置失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新机器人配置成功"}).Response()
}

func (c *BotController) GetBotConfig(ctx *gin.Context) {
	var req request.GetBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	currentUserId := c.CurrentUserId(ctx)

	configData, err := c.botService.GetBotConfigData(ctx, req.Id, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取机器人配置失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人配置成功", Data: configData}).Response()
}

func (c *BotController) SearchBotConfig(ctx *gin.Context) {
	var req request.SearchBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	currentUserId := c.CurrentUserId(ctx)
	configData, err := c.botService.SearchBotConfig(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取机器人配置失败: " + err.Error()}).Response()
		return
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人配置成功", Data: configData}).Response()
}

func (c *BotController) DelBotConfig(ctx *gin.Context) {
	var req request.GetBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	userId := c.CurrentUserId(ctx)
	err := c.botService.DeleteBotConfig(ctx, req.Id, userId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail}).Response()
		return
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人配置成功"}).Response()
}

