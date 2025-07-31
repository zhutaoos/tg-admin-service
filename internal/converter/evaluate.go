package converter

import (
	"app/internal/model"
	"app/internal/response"
	"encoding/json"
	"strconv"
	"time"
)

// ToEvaluateResponse 将 Model 转换为 Response（类似Java中的转换工具）
func ToEvaluateResponse(evaluate *model.JsEvaluateDB) *response.EvaluateResponse {
	if evaluate == nil {
		return nil
	}

	resp := &response.EvaluateResponse{
		ID:               evaluate.Id,
		GroupID:          evaluate.GroupID,
		UserId:           evaluate.UserId,
		UserName:         evaluate.UserName,
		NickName:         evaluate.NickName,
		EvaluateUserName: evaluate.EvaluateUserName,
		EvaluateUserId:   evaluate.EvaluateUserId,
		EvaluateNickName: evaluate.EvaluateNickName,
		CjDate:           evaluate.CjDate,
		Dj:               evaluate.Dj,
		Zb:               evaluate.Zb,
		Summary:          evaluate.Summary,
		Status:           evaluate.Status,
		StatusText:       getStatusText(evaluate.Status),
		CreatedAt:        time.Now(), // 如果 model 有创建时间字段，使用实际值
	}

	// 计算评分详情
	resp.Scores = response.ScoreDetail{
		Rz:    evaluate.Rz,
		Sc:    evaluate.Sc,
		Fw:    evaluate.Fw,
		Td:    evaluate.Td,
		Hj:    evaluate.Hj,
		Total: calculateTotalScore(evaluate.Rz, evaluate.Sc, evaluate.Fw, evaluate.Td, evaluate.Hj),
	}

	// 解析媒体数据
	resp.MediaList = parseMediaData(evaluate.CjMedia)

	return resp
}

// ToEvaluateListResponse 转换列表响应
func ToEvaluateListResponse(evaluates []*model.JsEvaluateDB, total int64, page, limit int) *response.EvaluateListResponse {
	list := make([]response.EvaluateResponse, 0, len(evaluates))
	for _, evaluate := range evaluates {
		if resp := ToEvaluateResponse(evaluate); resp != nil {
			list = append(list, *resp)
		}
	}

	return &response.EvaluateListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Limit: limit,
	}
}

// parseMediaData 解析JSON媒体数据
func parseMediaData(jsonData []byte) []response.MediaResponse {
	if len(jsonData) == 0 {
		return []response.MediaResponse{}
	}

	var mediaItems []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(jsonData, &mediaItems); err != nil {
		return []response.MediaResponse{}
	}

	result := make([]response.MediaResponse, 0, len(mediaItems))
	for _, item := range mediaItems {
		result = append(result, response.MediaResponse{
			Type: item.Type,
			URL:  item.URL,
			Name: item.Name,
		})
	}

	return result
}

// calculateTotalScore 计算总分
func calculateTotalScore(scores ...int) float64 {
	if len(scores) == 0 {
		return 0
	}

	total := 0
	for _, score := range scores {
		total += score
	}
	return float64(total) / float64(len(scores))
}

// getStatusText 获取状态文本
func getStatusText(status int32) string {
	switch status {
	case 0:
		return "待提交"
	case 1:
		return "已提交"
	case 2:
		return "审核通过"
	case 3:
		return "审核不通过"
	default:
		return "未知状态:" + strconv.Itoa(int(status))
	}
}
