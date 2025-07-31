package service

import (
	"app/internal/model"
	"app/internal/repository"
)

type EvaluateService interface {
	GetList(groupId string) ([]*model.JsEvaluateDB, error)
}

type EvaluateServiceImpl struct {
	evaluateRepo repository.EvaluateRepo
}

func NewEvaluateService(evaluateRepo repository.EvaluateRepo) EvaluateService {
	return &EvaluateServiceImpl{
		evaluateRepo: evaluateRepo,
	}
}

func (e *EvaluateServiceImpl) GetList(groupId string) ([]*model.JsEvaluateDB, error) {
	return e.evaluateRepo.GetList(groupId)
}
