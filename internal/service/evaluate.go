package service

import (
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"encoding/json"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type EvaluateService interface {
	GetList(request request.EvaluateSearchRequest) ([]*vo.JsEvaluateVo, error)
}

type EvaluateServiceImpl struct {
	db *gorm.DB
}

func NewEvaluateService(db *gorm.DB) EvaluateService {
	return &EvaluateServiceImpl{
		db: db,
	}
}

func (e *EvaluateServiceImpl) GetList(request request.EvaluateSearchRequest) ([]*vo.JsEvaluateVo, error) {
	var list []*model.JsEvaluateDB
	err := e.db.Where("group_id = ?", request.GroupId).Offset(request.GetOffset()).Limit(request.Limit).Find(&list).Error

	voList := make([]*vo.JsEvaluateVo, len(list))
	for i, db := range list {
		voList[i] = &vo.JsEvaluateVo{}
		copier.Copy(voList[i], db)

		// 解析 CjMedia JSON 字符串为 MediaVo 数组
		var mediaList []vo.MediaVo
		if len(db.CjMedia) > 0 {
			if err := json.Unmarshal(db.CjMedia, &mediaList); err == nil {
				voList[i].CjMediaList = mediaList
			}
		}
	}

	return voList, err
}
