package service

import (
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type EvaluateService interface {
	GetList(request request.EvaluateSearchRequest) (vo.PageResultVo[vo.JsEvaluateVo], error)
}

type EvaluateServiceImpl struct {
	db *gorm.DB
}

func NewEvaluateService(db *gorm.DB) EvaluateService {
	return &EvaluateServiceImpl{
		db: db,
	}
}

func (e *EvaluateServiceImpl) GetList(request request.EvaluateSearchRequest) (vo.PageResultVo[vo.JsEvaluateVo], error) {
	groupIds := request.GroupIds
	if len(groupIds) == 0 {
		return vo.PageResultVo[vo.JsEvaluateVo]{}, errors.New("group_id is required")
	}

	var list []*model.JsEvaluateDB
	var total int64
	err := e.db.Where("group_id in ?", groupIds).Offset(request.GetOffset()).Limit(request.Limit).Find(&list).Error
	if err != nil {
		return vo.PageResultVo[vo.JsEvaluateVo]{}, err
	}
	err = e.db.Model(&model.JsEvaluateDB{}).Where("group_id in ?", groupIds).Count(&total).Error
	if err != nil {
		return vo.PageResultVo[vo.JsEvaluateVo]{}, err
	}

	voList := make([]vo.JsEvaluateVo, len(list))
	for i, db := range list {
		voList[i] = vo.JsEvaluateVo{}
		copier.Copy(&voList[i], db)

		// 手动处理日期字段，格式化为字符串
		voList[i].CjDate = db.CjDate.Format("2006-01-02")
	}

	pageResultVo := vo.PageResultVo[vo.JsEvaluateVo]{
		Total: total,
		List:  voList,
	}

	return pageResultVo, nil
}
