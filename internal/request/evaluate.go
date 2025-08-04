package request

type EvaluateSearchRequest struct {
	PageRequest
	GroupId string `json:"group_id"`
}
