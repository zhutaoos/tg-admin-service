package repository

// Repository 数据访问层统一接口
type Repository interface {
	User() UserRepo
	Admin() AdminRepo
	Token() TokenRepo
}

// RepositoryImpl 实现统一的Repository接口
type RepositoryImpl struct {
	UserRepo  UserRepo
	AdminRepo AdminRepo
	TokenRepo TokenRepo
}

// User 返回用户仓储
func (r *RepositoryImpl) User() UserRepo {
	return r.UserRepo
}

// Admin 返回管理员仓储
func (r *RepositoryImpl) Admin() AdminRepo {
	return r.AdminRepo
}

// Token 返回Token仓储
func (r *RepositoryImpl) Token() TokenRepo {
	return r.TokenRepo
}
