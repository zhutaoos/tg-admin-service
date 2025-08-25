package provider

import (
	"app/internal/config"
	"app/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewUserService 创建用户服务Provider
func NewUserService(db *gorm.DB) service.UserService {
	return service.NewUserService(db)
}

// NewAdminService 创建管理员服务Provider
func NewAdminService(
	db *gorm.DB,
	tokenService service.TokenService,
	groupService service.GroupService,
) service.AdminService {
	return service.NewAdminService(db, tokenService, groupService)
}

// NewTokenService 创建Token服务Provider
func NewTokenService(redis *redis.Client, db *gorm.DB) service.TokenService {
	return service.NewTokenLogic(redis, db)
}

// NewEvaluateService 创建评价服务Provider
func NewEvaluateService(db *gorm.DB) service.EvaluateService {
	return service.NewEvaluateService(db)
}

// NewBotService 创建机器人服务Provider
func NewBotService(db *gorm.DB) *service.BotService {
	return service.NewBotService(db)
}

// NewGroupService 创建群组服务Provider
func NewGroupService(db *gorm.DB) service.GroupService {
	return service.NewGroupService(db)
}

// NewMessageService 创建消息服务Provider
func NewMessageService(db *gorm.DB) *service.MessageService {
	return service.NewMessageService(db)
}

// NewFileService 创建文件服务Provider
func NewFileService(config *config.Config) service.FileService {
	return service.NewFileService(config)
}
