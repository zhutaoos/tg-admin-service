package repository

import (
	"app/internal/model"

	"gorm.io/gorm"
)

// EvaluateRepo 评价数据访问接口
type EvaluateRepo interface {
	GetList(groupId string) ([]*model.JsEvaluateDB, error)
}

// EvaluateRepository 管理员数据访问实现
type EvaluateRepository struct {
	db *gorm.DB
}

// NewEvaluateRepository 创建管理员仓储实例
func NewEvaluateRepository(db *gorm.DB) EvaluateRepo {
	return &EvaluateRepository{db: db}
}

func (e *EvaluateRepository) GetList(groupId string) ([]*model.JsEvaluateDB, error) {
	var list []*model.JsEvaluateDB
	err := e.db.Where("group_id = ?", groupId).Find(&list).Error
	return list, err
}
