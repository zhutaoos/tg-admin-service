package request

import "app/internal/model"

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	Content       string              `json:"content" validate:"required"`     // 消息内容
	Images        model.JSONFileSlice `json:"images,omitempty"`                // 图片列表
	Medias        model.JSONFileSlice `json:"medias,omitempty"`                // 视频列表
	AdNickname    *string             `json:"adNickname,omitempty"`            // 被推广人花名
	AdUserID      *int                `json:"adUserId,omitempty"`              // 被推广人用户id
	AdGroupLink   *string             `json:"adGroupLink,omitempty"`           // 被推广人群组链接
	AdChannelLink *string             `json:"adChannelLink,omitempty"`         // 被推广人频道链接
}

// UpdateMessageRequest 更新消息请求
type UpdateMessageRequest struct {
	ID      uint                `json:"id" binding:"required" validate:"required"` // 消息ID
	Content string              `json:"content,omitempty"`                         // 消息内容
	Images  model.JSONFileSlice `json:"images,omitempty"`                          // 图片列表
	Medias  model.JSONFileSlice `json:"medias,omitempty"`                          // 视频列表
}

// GetMessageRequest 获取消息请求
type GetMessageRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"` // 消息ID
}

// SearchMessageRequest 搜索消息请求
type SearchMessageRequest struct {
	Content       string  `json:"content,omitempty"`       // 消息内容（模糊搜索）
	Status        *int    `json:"status,omitempty"`        // 状态筛选
	AdNickname    *string `json:"adNickname,omitempty"`    // 被推广人花名（模糊搜索）
	AdUserID      *int    `json:"adUserId,omitempty"`      // 被推广人用户id
	AdGroupLink   *string `json:"adGroupLink,omitempty"`   // 被推广人群组链接（模糊搜索）
	AdChannelLink *string `json:"adChannelLink,omitempty"` // 被推广人频道链接（模糊搜索）
	PageRequest
}

// DeleteMessageRequest 删除消息请求
type DeleteMessageRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"` // 消息ID
}
