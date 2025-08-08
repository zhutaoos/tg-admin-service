package model

import (
	"gorm.io/datatypes"
)

// BotConfig represents bot configuration for a specific group
type BotConfig struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	GroupID   int64          `gorm:"uniqueIndex;not null"`
	Config    datatypes.JSON `gorm:"type:json"`
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
}

// BotConfigData represents the structure of bot configuration data
// This struct is used for JSON marshaling/unmarshaling the Config field
// 机器人配置数据结构
type BotConfigData struct {
	Name             string  `json:"name"`               // 机器人名称
	Token            string  `json:"token"`              // 机器人token
	GroupID          int64   `json:"group_id"`           // 群组ID
	InviteLink       string  `json:"invite_link"`        // 群组邀请链接
	SubscribeChannel string  `json:"subscribe_channel"`  // 订阅频道链接
	GroupNamePrefix  string  `json:"group_name_prefix"`  // 群组名称前缀
	AdminIDs         []int64 `json:"admin_ids"`          // 管理员ID列表
}

// BotFeatures represents bot features for a specific group
type BotFeatures struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	GroupID     int64          `gorm:"index;not null"`
	FeatureName string         `gorm:"index;not null"`
	Enabled     bool           `gorm:"default:true"`
	Config      datatypes.JSON `gorm:"type:json"`
	CreatedAt   int64          `gorm:"autoCreateTime"`
	UpdatedAt   int64          `gorm:"autoUpdateTime"`
}

// SubscribeCheckConfig represents the configuration for subscribe check feature
type SubscribeCheckConfig struct {
	Enabled        bool    `json:"enabled"`
	Channels       []int64 `json:"channels"`
	WelcomeMessage string  `json:"welcome_message"`
}

// TableName sets the table name for BotConfig
func (BotConfig) TableName() string {
	return "bot_config"
}

// TableName sets the table name for BotFeatures
func (BotFeatures) TableName() string {
	return "bot_features"
}