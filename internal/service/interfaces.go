package service

import (
	"app/internal/model"
)

// Service 业务逻辑层接口
type Service interface {
	User() UserService
	Admin() AdminService
	Token() TokenService
}

// ServiceImpl 实现统一的Service接口
type ServiceImpl struct {
	UserService  UserService
	AdminService AdminService
	TokenService TokenService
}

// User 返回用户服务
func (s *ServiceImpl) User() UserService {
	return s.UserService
}

// Admin 返回管理员服务
func (s *ServiceImpl) Admin() AdminService {
	return s.AdminService
}

// Token 返回Token服务
func (s *ServiceImpl) Token() TokenService {
	return s.TokenService
}

// LoginResult 登录结果
type LoginResult struct {
	Token     string       `json:"token"`
	TokenInfo interface{}  `json:"token_info"`
	User      *model.Admin `json:"user"`
}
