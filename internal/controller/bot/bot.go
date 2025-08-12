package controller

import (
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

type BotController struct {
	botService *service.BotService
}

func NewBotController(botService *service.BotService) *BotController {
	return &BotController{
		botService: botService,
	}
}

// Bot Config Related Methods

// CreateBotConfig creates a new bot configuration
func (c *BotController) CreateBotConfig(ctx *gin.Context) {
	var req request.CreateBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	err := c.botService.CreateBotConfig(ctx, req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建机器人配置失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建成功"}).Response()
}

// UpdateBotConfig updates bot configuration by group_id
func (c *BotController) UpdateBotConfig(ctx *gin.Context) {
	var req request.UpdateBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	err := c.botService.UpdateBotConfig(ctx, req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "更新机器人配置失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新机器人配置成功"}).Response()
}

// GetBotConfig retrieves bot configuration by group_id
func (c *BotController) GetBotConfig(ctx *gin.Context) {
	var req request.GetBotConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	configData, err := c.botService.GetBotConfigData(ctx, req.Id)
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
	configData, err := c.botService.SearchBotConfig(ctx, req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取机器人配置失败: " + err.Error()}).Response()
		return
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人配置成功", Data: configData}).Response()
}

// Bot Features Related Methods

// CreateBotFeature creates a new bot feature configuration
// func (c *BotController) CreateBotFeature(ctx *gin.Context) {
// 	var req request.CreateBotFeatureRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
// 		return
// 	}

// 	err := c.botService.CreateBotFeature(ctx, req.GroupID, req.FeatureName, req.Enabled, req.Config)
// 	if err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建机器人功能失败: " + err.Error()}).Response()
// 		return
// 	}

// 	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建机器人功能成功"}).Response()
// }

// // UpdateBotFeature updates bot feature configuration
// func (c *BotController) UpdateBotFeature(ctx *gin.Context) {
// 	var req request.UpdateBotFeatureRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
// 		return
// 	}

// 	err := c.botService.UpdateBotFeature(ctx, req.GroupID, req.FeatureName, req.Enabled, req.Config)
// 	if err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "更新机器人功能失败: " + err.Error()}).Response()
// 		return
// 	}

// 	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新机器人功能成功"}).Response()
// }

// // GetBotFeature retrieves bot feature configuration
// func (c *BotController) GetBotFeature(ctx *gin.Context) {
// 	var req request.GetBotFeatureRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
// 		return
// 	}

// 	botFeature, err := c.botService.GetBotFeature(ctx, req.GroupID, req.FeatureName)
// 	if err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取机器人功能失败: " + err.Error()}).Response()
// 		return
// 	}

// 	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人功能成功", Data: botFeature}).Response()
// }

// // CreateSubscribeCheck creates subscribe check feature with default configuration
// func (c *BotController) CreateSubscribeCheck(ctx *gin.Context) {
// 	var req request.CreateSubscribeCheckRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
// 		return
// 	}

// 	err := c.botService.CreateSubscribeCheckFeature(ctx, req.GroupID)
// 	if err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建订阅检查功能失败: " + err.Error()}).Response()
// 		return
// 	}

// 	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建订阅检查功能成功"}).Response()
// }

// // ListBotFeatures lists all features for a group
// func (c *BotController) ListBotFeatures(ctx *gin.Context) {
// 	var req request.ListBotFeaturesRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
// 		return
// 	}

// 	features, err := c.botService.ListBotFeaturesByGroup(ctx, req.GroupID)
// 	if err != nil {
// 		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取机器人功能列表失败: " + err.Error()}).Response()
// 		return
// 	}

// 	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取机器人功能列表成功", Data: features}).Response()
// }

// RegisterRoutes registers all bot related routes
func (c *BotController) RegisterRoutes(router *gin.RouterGroup) {
	bot := router.Group("/bot")
	{
		// Bot Config routes
		bot.POST("/config/create", c.CreateBotConfig)
		bot.POST("/config/update", c.UpdateBotConfig)
		bot.POST("/config/get", c.GetBotConfig)

		// Bot Features routes
		// bot.POST("/features/create", c.CreateBotFeature)
		// bot.POST("/features/update", c.UpdateBotFeature)
		// bot.POST("/features/get", c.GetBotFeature)
		// bot.POST("/features/subscribe-check/create", c.CreateSubscribeCheck)
		// bot.POST("/features/list", c.ListBotFeatures)
	}
}
