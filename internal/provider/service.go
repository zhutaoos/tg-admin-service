package provider

import (
	"app/internal/repository"
	"app/internal/service"

	"github.com/redis/go-redis/v9"
)

// NewUserService 创建用户服务Provider
func NewUserService(userRepo repository.UserRepo) service.UserService {
	return service.NewUserService(userRepo)
}

// NewAdminService 创建管理员服务Provider
func NewAdminService(
	adminRepo repository.AdminRepo,
	tokenService service.TokenService,
) service.AdminService {
	return service.NewAdminService(adminRepo, tokenService)
}

// NewTokenService 创建Token服务Provider
func NewTokenService(redis *redis.Client, tokenRepo repository.TokenRepo) service.TokenService {
	return service.NewTokenLogic(redis, tokenRepo)
}
