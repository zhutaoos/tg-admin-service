package config

import (
	"gorm.io/gorm"
)

var db *gorm.DB

const UserPwdSalt = "test"

type DbConf struct {
	UserName string
	Password string
	Ip       string
	Port     string
	DbName   string
}

func Db() *gorm.DB {
	return db.Debug()
}
