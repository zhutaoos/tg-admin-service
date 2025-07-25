package repository

import (
	"app/internal/model"

	"gorm.io/gorm"
)

// TokenRepo Token数据访问接口
type TokenRepo interface {
	CreateToken(token *model.Token) error
	DelToken(token *model.Token) error
	GetToken(token *model.Token) (*model.Token, error)
}

// TokenRepository Token数据访问实现
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository 创建Token仓储实例
func NewTokenRepository(db *gorm.DB) TokenRepo {
	return &TokenRepository{db: db}
}

// CreateToken 创建Token记录
func (tr *TokenRepository) CreateToken(token *model.Token) error {
	return tr.db.Create(token).Error
}

// DelToken 删除Token记录
func (tr *TokenRepository) DelToken(token *model.Token) error {
	where := make(map[string]interface{})

	if token.UserId > 0 {
		where["user_id"] = token.UserId
	}
	if token.Id > 0 {
		where["id"] = token.Id
	}
	if token.Token != "" {
		where["token"] = token.Token
	}

	return tr.db.Where(where).Delete(&model.Token{}).Error
}

// GetToken 根据条件查询Token
func (tr *TokenRepository) GetToken(token *model.Token) (*model.Token, error) {
	err := tr.db.Where(token).First(token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetByToken 根据token字符串查询Token记录
func (tr *TokenRepository) GetByToken(tokenStr string) (*model.Token, error) {
	var token model.Token
	err := tr.db.Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByUserId 根据用户ID查询Token记录
func (tr *TokenRepository) GetByUserId(userId uint) (*model.Token, error) {
	var token model.Token
	err := tr.db.Where("user_id = ?", userId).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteByUserId 根据用户ID删除Token记录
func (tr *TokenRepository) DeleteByUserId(userId uint) error {
	return tr.db.Where("user_id = ?", userId).Delete(&model.Token{}).Error
}

// DeleteExpiredTokens 删除过期的Token记录
func (tr *TokenRepository) DeleteExpiredTokens(currentTime int64) error {
	return tr.db.Where("expire_time < ?", currentTime).Delete(&model.Token{}).Error
}
