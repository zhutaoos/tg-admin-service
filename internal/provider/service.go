package provider

import (
	"app/internal/repository"
	"app/internal/service"

	"github.com/redis/go-redis/v9"
)

// NewService 创建统一的Service Provider
func NewService(
	repo repository.Repository,
	redis *redis.Client,
) service.Service {
	return &service.ServiceImpl{
		UserService:  service.NewUserService(repo.User()),
		AdminService: service.NewAdminService(repo.Admin(), service.NewTokenLogic(redis, repo.Token())),
		TokenService: service.NewTokenLogic(redis, repo.Token()),
	}
}
