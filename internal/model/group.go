package model

import "time"

type Group struct {
	ID         uint      `gorm:"primaryKey; type:INT NOT NULL AUTO_INCREMENT"`
	AdminID    int       `gorm:"column:admin_id; type:INT; default:NULL"`
	GroupID    int64     `gorm:"column:group_id; type:BIGINT; default:NULL"`
	GroupName  string    `gorm:"column:group_name; type:VARCHAR(128); default:NULL"`
	Status     int       `gorm:"column:status; type:INT; default:0; comment:'状态 0:正常 1:删除'"`
	CreateTime time.Time `gorm:"column:create_time; type:DATETIME; default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"column:update_time; type:DATETIME; default:NULL"`
}

func (Group) TableName() string {
	return "admin_group"
}

// GroupInfo 用于 Admin 模型的群组信息结构
type GroupInfo struct {
	ID   string `json:"id"`   // 群组ID，适配字符串格式如 "-1002714549168"
	Name string `json:"name"` // 群组名称
}