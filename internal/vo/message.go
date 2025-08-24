package vo

import (
	"app/internal/model"
	"time"
)

// MessageVO 消息响应结构
type MessageVO struct {
	ID         uint                     `json:"id"`         // 消息ID
	GroupID    int                      `json:"groupId"`    // 群组ID
	Content    string                   `json:"content"`    // 消息内容
	Images     model.JSONStringSlice    `json:"images"`     // 图片列表
	Medias     model.JSONStringSlice    `json:"medias"`     // 视频列表
	CreateTime *time.Time               `json:"createTime"` // 创建时间
}

// PageVO 分页响应结构（通用）
type PageVO[T any] struct {
	List  []T   `json:"list"`  // 数据列表
	Total int64 `json:"total"` // 总数
	Page  int   `json:"page"`  // 页码
	Limit int   `json:"limit"` // 每页数量
}