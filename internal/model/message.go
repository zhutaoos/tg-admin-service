package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Message struct {
	ID         uint           `gorm:"primaryKey; type:INT NOT NULL AUTO_INCREMENT"`
	AdminID    int            `gorm:"column:admin_id; type:INT; default:NULL"`
	GroupID    int            `gorm:"column:group_id; type:INT; default:NULL; comment:'群组id,chatid'"`
	Content    string         `gorm:"column:content; type:VARCHAR(2048); default:NULL; comment:'消息内容'"`
	Images     JSONStringSlice `gorm:"column:images; type:JSON; default:NULL; comment:'图片'"`
	Medias     JSONStringSlice `gorm:"column:medias; type:JSON; default:NULL; comment:'视频'"`
	CreateTime *time.Time     `gorm:"column:create_time; type:DATETIME; default:NULL"`
	UpdateTime *time.Time     `gorm:"column:update_time; type:DATETIME; default:NULL"`
	Status     int            `gorm:"column:status; type:INT; default:0; comment:'是否删除 0:正常 1:删除'"`
}

func (Message) TableName() string {
	return "message"
}

// JSONStringSlice 自定义 JSON 类型用于处理 []string
type JSONStringSlice []string

// Value 实现 driver.Valuer 接口，将 Go 类型转换为数据库值
func (j JSONStringSlice) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 Go 类型
func (j *JSONStringSlice) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, j)
}