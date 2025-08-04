package service

import (
	"app/internal/model"
	"app/internal/request"

	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	ListUser(req request.UserSearchRequest) ([]model.User, int64, error)
	GetUserById(id int64, user *model.User) error
	LoadUser(uid string) (*model.User, error)
}

type UserServiceImpl struct {
	db *gorm.DB
}

// NewUserService 创建UserService实例
func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		db: db,
	}
}

func (u *UserServiceImpl) GetUserById(id int64, user *model.User) error {
	return u.db.Where("id = ?", id).First(user).Error
}

func (u *UserServiceImpl) LoadUser(uid string) (*model.User, error) {
	user := &model.User{}
	return user, u.db.Where("id = ?", uid).First(user).Error
}

func (u *UserServiceImpl) ListUser(req request.UserSearchRequest) ([]model.User, int64, error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	var users []model.User
	var total int64

	query := u.db.Model(&model.User{})

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
