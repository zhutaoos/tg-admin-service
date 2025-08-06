package request

type EvaluateSearchRequest struct {
	PageRequest
	GroupIds []string `json:"group_ids"`
}
