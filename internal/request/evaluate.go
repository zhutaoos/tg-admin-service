package request

type EvaluateSearchRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	GroupId string `json:"group_id"`
}
