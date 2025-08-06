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
	UpdateEvaluate(param request.EvaluateUpdateParam) error
}

type EvaluateServiceImpl struct {
	db *gorm.DB
}

func NewEvaluateService(db *gorm.DB) EvaluateService {
	return &EvaluateServiceImpl{
		db: db,
	}
}

func (e *EvaluateServiceImpl) UpdateEvaluate(param request.EvaluateUpdateParam) error {
	// 创建更新字段结构体
	updateFields := request.EvaluateUpdateFields{
		Dj:      param.Dj,
		Rz:      param.Rz,
		Sc:      param.Sc,
		Fw:      param.Fw,
		Td:      param.Td,
		Hj:      param.Hj,
		Zb:      param.Zb,
		Summary: param.Summary,
		Status:  param.Status,
	}

	return e.db.Model(&model.JsEvaluateDB{}).
		Where("id = ?", param.Id).
		Select("dj, rz, sc, fw, td, hj, zb, summary, status").
		Updates(updateFields).Error
}

func (e *EvaluateServiceImpl) GetList(request request.EvaluateSearchRequest) (vo.PageResultVo[vo.JsEvaluateVo], error) {
	groupIds := request.GroupIds
	if len(groupIds) == 0 {
		return vo.PageResultVo[vo.JsEvaluateVo]{}, errors.New("group_id is required")
	}

	var list []*model.JsEvaluateDB
	var total int64

	// 构建查询条件
	query := e.db.Where("group_id in ?", groupIds)

	// 根据请求参数动态添加where条件
	if request.Status > 0 {
		query = query.Where("status = ?", request.Status)
	}

	if request.EvaluateNickName != "" {
		query = query.Where("evaluate_nick_name LIKE ?", "%"+request.EvaluateNickName+"%")
	}

	// 执行查询
	err := query.Offset(request.GetOffset()).Limit(request.Limit).Find(&list).Error
	if err != nil {
		return vo.PageResultVo[vo.JsEvaluateVo]{}, err
	}

	// 构建计数查询条件（与查询条件保持一致）
	countQuery := e.db.Model(&model.JsEvaluateDB{}).Where("group_id in ?", groupIds)

	if request.Status > 0 {
		countQuery = countQuery.Where("status = ?", request.Status)
	}

	if request.EvaluateNickName != "" {
		countQuery = countQuery.Where("evaluate_nick_name LIKE ?", "%"+request.EvaluateNickName+"%")
	}

	err = countQuery.Count(&total).Error
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
