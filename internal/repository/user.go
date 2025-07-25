package repository

import (
	"app/internal/model"
	"app/internal/request"

	"gorm.io/gorm"
)

// UserRepo 用户数据访问接口
type UserRepo interface {
	GetUserInfo(user *model.User) error
	GetList(req request.UserSearchRequest) ([]model.User, int64, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int64) error
}

// UserRepository 用户数据访问实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepo {
	return &UserRepository{db: db}
}

// GetUserInfo 根据条件查询单个用户信息
func (ur *UserRepository) GetUserInfo(user *model.User) error {
	return ur.db.Where(user).First(user).Error
}

// GetList 获取用户列表（带分页和搜索条件）
func (ur *UserRepository) GetList(req request.UserSearchRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := ur.db.Model(&model.User{})

	// 构建查询条件
	if req.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Create 创建用户
func (ur *UserRepository) Create(user *model.User) error {
	return ur.db.Create(user).Error
}

// Update 更新用户
func (ur *UserRepository) Update(user *model.User) error {
	return ur.db.Save(user).Error
}

// Delete 删除用户
func (ur *UserRepository) Delete(id int64) error {
	return ur.db.Delete(&model.User{}, id).Error
}

// GetByUserId 根据UserId查询用户
func (ur *UserRepository) GetByUserId(userId string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByNickname 根据昵称查询用户
func (ur *UserRepository) GetByNickname(nickname string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("nickname = ?", nickname).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
