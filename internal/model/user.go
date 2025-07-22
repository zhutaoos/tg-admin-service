package model

import (
	"app/internal/config"
	"app/internal/request"
)

type User struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `redis:"Id" json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	UserId          string         `redis:"UserId" json:"user_id" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Nickname        string         `redis:"Nickname" json:"nickname" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Avatar          string         `redis:"Avatar" json:"avatar" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Status          uint8          `redis:"Status" json:"status" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	CreateTime      int64          `redis:"CreateTime" json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"` // 自动写入时间戳
	CreateTimeStr   string         `redis:"-" json:"create_time_str" gorm:"-:all"`                                              // -:all 无读写迁移权限，该字段不在数据库中
}

func (user *User) GetUserInfo() {
	config.Db().Where(user).First(user)
}

func (user *User) GetList(req request.UserSearchRequest) []User {
	list := make([]User, 0)
	config.Db().Where(user).Find(&list)

	return list
}
