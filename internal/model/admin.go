package model

import (
	"app/internal/config"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type Admin struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey; type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT"`
	GroupInfo       datatypes.JSON `json:"groupInfo" gorm:"type:json; column:group_info; default:[]"`
	Account         string         `json:"account" gorm:"type:VARCHAR(32) NOT NULL;  default:''; comment:登录账号"`
	Password        string         `json:"-" gorm:"type:VARCHAR(128) NOT NULL;"`
	Name            string         `json:"name" gorm:"type:VARCHAR(32) NOT NULL; default:''"`
	Status          uint           `json:"status" gorm:"type:INT(4) UNSIGNED NOT NULL; default:1; comment:0 禁用 1 启用"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:DATETIME NOT NULL; default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"type:DATETIME; default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (admin *Admin) GetAdmin() *Admin {
	config.Db().Where(admin).First(admin)
	return admin
}

func (admin *Admin) GetList(pid uint) []*Admin {
	list := make([]*Admin, 0)

	return list
}

func (admin *Admin) SetAdmin() *Admin {
	if admin.Id <= 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil
		}
		admin.Password = string(hashedPassword)
		config.Db().Create(&admin)
	} else {
		config.Db().Select("title", "content", "status", "cate_id", "type").Model(&admin).Updates(&admin)
	}
	return admin
}

func (admin *Admin) DelAdmin(ids []int) {
	config.Db().Delete(&admin, ids)
}
