package response

import "time"

// EvaluateResponse 评价响应对象，类似Java中的VO
type EvaluateResponse struct {
	ID               string          `json:"id"`
	GroupID          string          `json:"group_id"`
	UserId           int64           `json:"user_id"`
	UserName         string          `json:"user_name"`
	NickName         string          `json:"nick_name"`
	EvaluateUserName string          `json:"evaluate_user_name"`
	EvaluateUserId   int64           `json:"evaluate_user_id"`
	EvaluateNickName string          `json:"evaluate_nick_name"`
	CjDate           string          `json:"cj_date"`
	Dj               int             `json:"dj"`
	Scores           ScoreDetail     `json:"scores"` // 评分详情
	Zb               string          `json:"zb"`
	Summary          string          `json:"summary"`
	MediaList        []MediaResponse `json:"media_list"` // 媒体文件列表
	Status           int32           `json:"status"`
	StatusText       string          `json:"status_text"` // 状态文本
	CreatedAt        time.Time       `json:"created_at"`
}

// ScoreDetail 评分详情
type ScoreDetail struct {
	Rz    int     `json:"rz"`    // 人照评分
	Sc    int     `json:"sc"`    // 身材评分
	Fw    int     `json:"fw"`    // 服务评分
	Td    int     `json:"td"`    // 态度评分
	Hj    int     `json:"hj"`    // 环境评分
	Total float64 `json:"total"` // 总分
}

// MediaResponse 媒体文件响应
type MediaResponse struct {
	Type string `json:"type"` // image/video
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size,omitempty"`
}

// EvaluateListResponse 评价列表响应
type EvaluateListResponse struct {
	List  []EvaluateResponse `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

// EvaluateStatResponse 评价统计响应
type EvaluateStatResponse struct {
	TotalCount    int64   `json:"total_count"`
	AvgScore      float64 `json:"avg_score"`
	PendingCount  int64   `json:"pending_count"`
	ApprovedCount int64   `json:"approved_count"`
	RejectedCount int64   `json:"rejected_count"`
}
