package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/internal/request"

	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务接口
type AdminService interface {
	Login(req request.AdminLoginRequest) (*LoginResult, error)
	GetProfile(adminId int64) (*model.Admin, error)
	InitPwd(req request.InitPwdRequest) error
}

// LoginResult 登录结果
type LoginResult struct {
	Token     string       `json:"token"`
	TokenInfo interface{}  `json:"token_info"`
	User      *model.Admin `json:"user"`
}

type AdminServiceImpl struct {
	adminRepo  repository.AdminRepo
	tokenLogic TokenService
}

// NewAdminService 创建AdminService实例
func NewAdminService(
	adminRepo repository.AdminRepo,
	tokenLogic TokenService,
) AdminService {
	return &AdminServiceImpl{
		adminRepo:  adminRepo,
		tokenLogic: tokenLogic,
	}
}

// Login 管理员登录
func (as *AdminServiceImpl) Login(req request.AdminLoginRequest) (*LoginResult, error) {
	// 根据账号查询管理员
	admin, err := as.adminRepo.GetByAccount(req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	// 生成JWT
	token, userJwt := as.tokenLogic.GenerateJwt(admin.Id, 0)

	// 组装登录结果
	result := &LoginResult{
		Token:     token,
		TokenInfo: userJwt,
		User:      admin,
	}

	return result, nil
}

// GetProfile 获取管理员信息
func (as *AdminServiceImpl) GetProfile(adminId int64) (*model.Admin, error) {
	return as.adminRepo.GetByID(uint(adminId))
}

// InitPwd 初始化密码
func (as *AdminServiceImpl) InitPwd(req request.InitPwdRequest) error {
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 这里可以根据具体业务逻辑决定如何处理初始化密码
	// 例如创建默认管理员账号或更新特定管理员密码
	// 暂时返回哈希后的密码供调用方使用
	_ = string(hashedPassword)

	return nil
}

// CreateAdmin 创建管理员
func (as *AdminServiceImpl) CreateAdmin(admin *model.Admin) error {
	// 密码加密
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin.Password = string(hashedPassword)
	}

	return as.adminRepo.Create(admin)
}

// UpdateAdmin 更新管理员信息
func (as *AdminServiceImpl) UpdateAdmin(admin *model.Admin) error {
	return as.adminRepo.Update(admin)
}

// UpdatePassword 更新管理员密码
func (as *AdminServiceImpl) UpdatePassword(adminId uint, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return as.adminRepo.UpdatePassword(adminId, string(hashedPassword))
}

// DeleteAdmin 删除管理员
func (as *AdminServiceImpl) DeleteAdmin(id int64) error {
	return as.adminRepo.Delete(id)
}
