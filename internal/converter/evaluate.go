package converter

import (
	"app/internal/model"
	"app/internal/vo"
)

// ToEvaluateResponse 将 Model 转换为 Response（类似Java中的转换工具）
func ToEvaluateResponse(evaluate *model.JsEvaluateDB) *vo.JsEvaluateVo {
	if evaluate == nil {
		return nil
	}

	resp := &vo.JsEvaluateVo{
		Id:               evaluate.Id,
		GroupID:          evaluate.GroupID,
		UserId:           evaluate.UserId,
		UserName:         evaluate.UserName,
		NickName:         evaluate.NickName,
		EvaluateUserName: evaluate.EvaluateUserName,
		EvaluateUserId:   evaluate.EvaluateUserId,
		EvaluateNickName: evaluate.EvaluateNickName,
		CjDate:           evaluate.CjDate.Format("2006-01-02"),
		Dj:               int(evaluate.Dj),
		Zb:               evaluate.Zb,
		Rz:               evaluate.Rz,
		Sc:               evaluate.Sc,
		Fw:               evaluate.Fw,
		Td:               evaluate.Td,
		Hj:               evaluate.Hj,
		Summary:          evaluate.Summary,
		Status:           evaluate.Status,
	}

	return resp
}

// ToEvaluateListResponse 转换列表响应
func ToEvaluateListResponse(evaluates []*model.JsEvaluateDB, total int64, page, limit int) *vo.JsEvaluateListVo {
	list := make([]*vo.JsEvaluateVo, 0, len(evaluates))
	for _, evaluate := range evaluates {
		if resp := ToEvaluateResponse(evaluate); resp != nil {
			list = append(list, resp)
		}
	}

	return &vo.JsEvaluateListVo{
		List:  list,
		Total: total,
	}
}