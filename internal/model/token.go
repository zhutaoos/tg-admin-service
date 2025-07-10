package model

import (
	"app/internal/config"
)

type Token struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	UserId          uint           `json:"user_id" gorm:"type:INT(8) UNSIGNED NOT NULL;default:0"`
	ExpireTime      int64          `json:"expire_time" gorm:"type:BIGINT UNSIGNED NOT NULL;default:0"`
	Token           string         `json:"token" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	CreateTime      int64          `json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"`
	CreateTimeStr   string         `json:"create_time_str" gorm:"-:all"`
}

func (t *Token) CheckToken(token string) *Token {
	config.Db().First(t, "token = ? ", token)
	return t
}

func (t *Token) CreateToken() *Token {
	config.Db().Create(&t)
	return t
}

func (t *Token) DelToken() {
	where := make(map[string]interface{})
	if t.UserId > 0 {
		where["user_id"] = t.UserId
	}
	if t.Id > 0 {
		where["id"] = t.Id
	}
	if t.Token != "" {
		where["token"] = t.Token
	}
	config.Db().Where(where).Delete(&t)
}
