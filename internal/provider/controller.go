package provider

import (
	admin_controller "app/internal/controller/admin"
	user_controller "app/internal/controller/user"
	"app/internal/service"
)

// NewAdminController 创建管理员控制器Provider
func NewAdminController(service service.Service) *admin_controller.AdminController {
	return admin_controller.NewAdminController(service.Admin())
}

// NewUserController 创建用户控制器Provider
func NewUserController(service service.Service) *user_controller.UserController {
	return user_controller.NewUserController(service.User())
}
