package service

import (
	"app/internal/model"
	"app/internal/request"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理员服务接口
type AdminService interface {
	Login(req request.AdminLoginRequest) (*LoginResult, error)
	GetProfile(adminId int64) (*model.Admin, error)
	InitPwd(req request.InitPwdRequest) error
	GetAdminById(id int64, admin *model.Admin) error
	GetAdminWithGroups(ctx context.Context, adminId int64) (*model.Admin, error)
}

// LoginResult 登录结果
type LoginResult struct {
	Token     string       `json:"token"`
	TokenInfo interface{}  `json:"token_info"`
	User      *model.Admin `json:"user"`
}

type AdminServiceImpl struct {
	db          *gorm.DB
	tokenLogic  TokenService
	groupService GroupService
}

// NewAdminService 创建AdminService实例
func NewAdminService(
	db *gorm.DB,
	tokenLogic TokenService,
	groupService GroupService,
) AdminService {
	return &AdminServiceImpl{
		db:          db,
		tokenLogic:  tokenLogic,
		groupService: groupService,
	}
}

func (as *AdminServiceImpl) GetAdminById(id int64, admin *model.Admin) error {
	return as.db.Where("id = ?", id).First(admin).Error
}

// Login 管理员登录
func (as *AdminServiceImpl) Login(req request.AdminLoginRequest) (*LoginResult, error) {
	// 根据账号查询管理员
	admin := &model.Admin{}
	err := as.db.Where("account = ?", req.Username).First(admin).Error
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
	return as.GetAdminWithGroups(context.Background(), adminId)
}

// GetAdminWithGroups 获取管理员信息及其关联的群组
func (as *AdminServiceImpl) GetAdminWithGroups(ctx context.Context, adminId int64) (*model.Admin, error) {
	admin := &model.Admin{}
	err := as.db.Where("id = ?", adminId).First(admin).Error
	if err != nil {
		return nil, err
	}

	// 获取管理员的群组信息
	groups, err := as.groupService.GetGroupsByAdminID(ctx, int(adminId))
	if err != nil {
		return admin, nil // 即使没有群组信息也返回管理员基本信息
	}

	// 转换为 GroupInfo 格式
	var groupInfos []model.GroupInfo
	for _, group := range groups {
		groupInfos = append(groupInfos, model.GroupInfo{
			ID:   group.ID,
			Name: group.Name,
		})
	}
	admin.GroupInfo = groupInfos

	return admin, nil
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
