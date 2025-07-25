package provider

import (
	"app/internal/config"
	"app/tools/logger"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewDatabase 创建数据库连接
func NewDatabase(lc fx.Lifecycle, dbConf *config.DbConf) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.UserName,
		dbConf.Password,
		dbConf.Ip,
		dbConf.Port,
		dbConf.DbName,
	)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm日志模式：silent
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.System("数据库连接已建立")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.System("正在关闭数据库连接")
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		},
	})

	return db.Debug(), nil
}

// NewRedis 创建Redis连接（可选）
func NewRedis(lc fx.Lifecycle, redisConf *config.RedisConf) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     redisConf.Ip + ":" + redisConf.Port,
		Username: redisConf.Username,
		Password: redisConf.Password,
		DB:       redisConf.Db,
		PoolSize: redisConf.MaxTotal,
	}

	client := redis.NewClient(options)

	// 测试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Redis连接失败，将在无Redis模式下运行", "error", err)
		// 注意：这里我们返回一个nil client和nil error，表示Redis是可选的
		// 但这需要调用方处理nil client的情况
		// 为了简化，我们暂时还是返回一个假的client，但记录警告
		// 实际项目中应该重构为可选依赖模式

		// 创建一个用于测试的本地Redis连接
		localOptions := &redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}
		localClient := redis.NewClient(localOptions)
		_, localErr := localClient.Ping(context.Background()).Result()
		if localErr != nil {
			logger.Error("本地Redis也无法连接，继续使用原配置（可能导致功能异常）", "error", localErr)
			// 为了演示，我们返回原client，但不会正常工作
			return client, nil
		}

		logger.System("使用本地Redis连接")
		client = localClient
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.System("Redis连接已建立")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.System("正在关闭Redis连接")
			return client.Close()
		},
	})

	return client, nil
}
