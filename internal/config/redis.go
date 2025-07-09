package config

import (
	"app/tools/key_utils"
	"app/tools/logger"
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisClient 定义一个全局变量
var RedisClient = &redis.Client{}
var KeyUtils = &key_utils.KeyUtils{}

type RedisConf struct {
	Ip       string
	Port     string
	Username string
	Password string
	Db       int
	MaxTotal int
}

func InitRedis(c *RedisConf) {
	o := &redis.Options{
		Addr:     c.Ip + ":" + c.Port,
		Username: c.Username,
		Password: c.Password,
		DB:       c.Db,
		PoolSize: c.MaxTotal,
	}

	RedisClient = redis.NewClient(o)
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		println(err.Error())
		logger.Error("REDIS CONNECT FAIL", err.Error())
	} else {
		logger.System("REDIS INIT SUCCESS")
	}

	KeyUtils.BaseName = ServerName
}
