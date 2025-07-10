package logic

import (
	"app/internal/model"
)

type UserLogic struct {
}

var UserLogicInstance UserLogic

func init() {
	UserLogicInstance = UserLogic{}
}

// LoadUser 根据 uid 搜索用户
func (u UserLogic) LoadUser(uid uint) *model.User {
	userModel := &model.User{Id: uid}
	if userModel.Id == 0 {
		userModel.GetUserInfo()
	}
	return userModel
}

func (u UserLogic) SearchUser(search map[string]interface{}) *model.User {
	user := &model.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	user.GetUserInfo()
	return user
}
