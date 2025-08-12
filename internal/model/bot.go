package model

import (
	"time"

	"gorm.io/datatypes"
)

// BotConfig represents bot configuration for a specific group
type BotConfig struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	Region     string         `gorm:"not null"`
	GroupID    int64          `gorm:"uniqueIndex;not null"`
	Config     datatypes.JSON `gorm:"type:json column:config;default:{}"`
	Features   datatypes.JSON `gorm:"type:json column:features;default:[]"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime"`
}

// TableName sets the table name for BotConfig
func (BotConfig) TableName() string {
	return "bot_config"
}
