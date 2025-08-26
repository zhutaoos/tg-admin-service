package vo

import (
	"app/internal/model"
	"time"
)

// MessageVO 消息响应结构
type MessageVO struct {
	ID            uint                `json:"id"`                      // 消息ID
	Content       string              `json:"content"`                 // 消息内容
	Images        model.JSONFileSlice `json:"images"`                  // 图片列表
	Medias        model.JSONFileSlice `json:"medias"`                  // 视频列表
	AdNickname    *string             `json:"adNickname,omitempty"`    // 被推广人花名
	AdUserID      *int                `json:"adUserId,omitempty"`      // 被推广人用户id
	AdGroupLink   *string             `json:"adGroupLink,omitempty"`   // 被推广人群组链接
	AdChannelLink *string             `json:"adChannelLink,omitempty"` // 被推广人频道链接
	CreateTime    *time.Time          `json:"createTime"`              // 创建时间
}

// PageVO 分页响应结构（通用）
type PageVO[T any] struct {
	List  []T   `json:"list"`  // 数据列表
	Total int64 `json:"total"` // 总数
	Page  int   `json:"page"`  // 页码
	Limit int   `json:"limit"` // 每页数量
}
