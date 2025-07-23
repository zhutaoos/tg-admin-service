package model

type User struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              int64          `redis:"Id" json:"id" gorm:"primaryKey;type:BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT;comment:ID"`
	UserId          string         `redis:"UserId" json:"user_id" gorm:"type:VARCHAR(255) NOT NULL;default:'';comment:用户ID"`
	Nickname        string         `redis:"Nickname" json:"nickname" gorm:"type:VARCHAR(255) NOT NULL;default:'';comment:昵称"`
	Avatar          string         `redis:"Avatar" json:"avatar" gorm:"type:VARCHAR(255) NOT NULL;default:'';comment:头像"`
	Status          int            `redis:"Status" json:"status" gorm:"type:INT(11) UNSIGNED NOT NULL;default:0;comment:状态"`
	CreateTime      int64          `redis:"CreateTime" json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;comment:创建时间"` // 自动写入时间戳
	CreateTimeStr   string         `redis:"-" json:"create_time_str" gorm:"-:all"`                                                          // -:all 无读写迁移权限，该字段不在数据库中
}
