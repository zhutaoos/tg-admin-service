package service

import (
	"app/internal/model"
	"app/internal/query"
	"app/internal/request"
)

type UserService struct {
	userQuery query.UserQuery
}

func (u *UserService) UserList(req request.UserSearchRequest) ([]model.User, int64) {
	return u.userQuery.GetList(req)
}

// LoadUser 根据 uid 搜索用户
func (u UserService) LoadUser(uid string) *model.User {
	userModel := &model.User{UserId: uid}

	return userModel
}

func (u UserService) SearchUser(search map[string]interface{}) *model.User {
	user := &model.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	return user
}
