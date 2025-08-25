package request

import "app/internal/model"

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	GroupIDs []int                 `json:"groupIds" binding:"required" validate:"required"` // 群组ID数组
	Content  string                `json:"content" validate:"required"`                     // 消息内容
	Images   model.JSONFileSlice   `json:"images,omitempty"`                                // 图片列表
	Medias   model.JSONFileSlice   `json:"medias,omitempty"`                                // 视频列表
}

// UpdateMessageRequest 更新消息请求
type UpdateMessageRequest struct {
	ID       uint                `json:"id" binding:"required" validate:"required"` // 消息ID
	GroupIDs []int               `json:"groupIds,omitempty"`                         // 群组ID数组
	Content  string              `json:"content,omitempty"`                          // 消息内容
	Images   model.JSONFileSlice `json:"images,omitempty"`                           // 图片列表
	Medias   model.JSONFileSlice `json:"medias,omitempty"`                           // 视频列表
}

// GetMessageRequest 获取消息请求
type GetMessageRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"` // 消息ID
}

// SearchMessageRequest 搜索消息请求
type SearchMessageRequest struct {
	GroupID int    `json:"groupId,omitempty"` // 群组ID
	Content string `json:"content,omitempty"` // 消息内容（模糊搜索）
	Status  *int   `json:"status,omitempty"`  // 状态筛选
	PageRequest
}

// DeleteMessageRequest 删除消息请求
type DeleteMessageRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"` // 消息ID
}