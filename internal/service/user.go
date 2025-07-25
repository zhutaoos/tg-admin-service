package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/internal/request"
)

// UserService 用户服务接口
type UserService interface {
	UserList(req request.UserSearchRequest) ([]model.User, int64, error)
	LoadUser(uid string) (*model.User, error)
	SearchUser(search map[string]interface{}) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

type UserServiceImpl struct {
	userRepo repository.UserRepo
}

// NewUserService 创建UserService实例
func NewUserService(userRepo repository.UserRepo) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (u *UserServiceImpl) UserList(req request.UserSearchRequest) ([]model.User, int64, error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	return u.userRepo.GetList(req)
}

// LoadUser 根据 uid 搜索用户
func (u *UserServiceImpl) LoadUser(uid string) (*model.User, error) {
	userModel := &model.User{UserId: uid}
	err := u.userRepo.GetUserInfo(userModel)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (u *UserServiceImpl) SearchUser(search map[string]interface{}) (*model.User, error) {
	user := &model.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	err := u.userRepo.GetUserInfo(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser 创建用户
func (u *UserServiceImpl) CreateUser(user *model.User) error {
	return u.userRepo.Create(user)
}

// UpdateUser 更新用户
func (u *UserServiceImpl) UpdateUser(user *model.User) error {
	return u.userRepo.Update(user)
}

// DeleteUser 删除用户
func (u *UserServiceImpl) DeleteUser(id int64) error {
	return u.userRepo.Delete(id)
}
