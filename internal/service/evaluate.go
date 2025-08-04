package service

import (
	"app/internal/model"
	"app/internal/request"

	"gorm.io/gorm"
)

type EvaluateService interface {
	GetList(request request.EvaluateSearchRequest) ([]*model.JsEvaluateDB, error)
}

type EvaluateServiceImpl struct {
	db *gorm.DB
}

func NewEvaluateService(db *gorm.DB) EvaluateService {
	return &EvaluateServiceImpl{
		db: db,
	}
}

func (e *EvaluateServiceImpl) GetList(request request.EvaluateSearchRequest) ([]*model.JsEvaluateDB, error) {
	var list []*model.JsEvaluateDB
	err := e.db.Where("group_id = ?", request.GroupId).Offset(request.GetOffset()).Limit(request.Limit).Find(&list).Error
	return list, err
}
