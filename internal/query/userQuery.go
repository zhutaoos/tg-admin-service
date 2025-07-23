package query

import (
	"app/internal/config"
	"app/internal/model"
	"app/internal/request"
)

type UserQuery struct {
}

func (u *UserQuery) GetList(req request.UserSearchRequest) ([]model.User, int64) {
	query := config.Db().Model(&model.User{})
	// 构建查询条件
	if req.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 分页查询
	var total int64
	query.Count(&total)

	var users []model.User
	query.Offset(req.GetOffset()).Limit(req.Limit).Find(&users)

	return users, total
}
