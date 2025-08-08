package request

import "app/internal/model"

type CreateBotConfigRequest struct {
	GroupID int64                  `json:"group_id" binding:"required"`
	Config  model.BotConfigData    `json:"config" binding:"required"`
}

type UpdateBotConfigRequest struct {
	GroupID int64                  `json:"group_id" binding:"required"`
	Config  model.BotConfigData    `json:"config" binding:"required"`
}

type GetBotConfigRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
}