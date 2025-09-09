package model

import (
    "time"

    "gorm.io/datatypes"
)

// BotType 机器人类型枚举
// 0: 功能型  1: 群发机器人
type BotType uint8

const (
    BotTypeFunctional BotType = 0 // 功能型
    BotTypeBroadcast  BotType = 1 // 群发机器人
)

// Valid 判断取值是否合法
func (t BotType) Valid() bool { return t == BotTypeFunctional || t == BotTypeBroadcast }

// String 返回中文描述
func (t BotType) String() string {
    switch t {
    case BotTypeFunctional:
        return "功能型"
    case BotTypeBroadcast:
        return "群发机器人"
    default:
        return "未知"
    }
}

// BotConfig represents bot configuration for a specific group
type BotConfig struct {
    ID         uint           `gorm:"primaryKey;autoIncrement"`
    AdminId    uint           `gorm:"not null"`
    Type       BotType        `gorm:"type:TINYINT unsigned;not null;default:0;comment:类型 0:功能型 1:群发机器人"`
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
