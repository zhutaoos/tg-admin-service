package provider

import (
	"app/internal/repository"

	"gorm.io/gorm"
)

// NewUserRepository 创建用户仓储Provider
func NewUserRepository(db *gorm.DB) repository.UserRepo {
	return repository.NewUserRepository(db)
}

// NewAdminRepository 创建管理员仓储Provider
func NewAdminRepository(db *gorm.DB) repository.AdminRepo {
	return repository.NewAdminRepository(db)
}

// NewTokenRepository 创建Token仓储Provider
func NewTokenRepository(db *gorm.DB) repository.TokenRepo {
	return repository.NewTokenRepository(db)
}
