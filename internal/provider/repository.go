package provider

import (
	"app/internal/repository"

	"gorm.io/gorm"
)

// NewRepository 创建统一的仓储Provider
func NewRepository(db *gorm.DB) repository.Repository {
	return &repository.RepositoryImpl{
		UserRepo:  repository.NewUserRepository(db),
		AdminRepo: repository.NewAdminRepository(db),
		TokenRepo: repository.NewTokenRepository(db),
	}
}
