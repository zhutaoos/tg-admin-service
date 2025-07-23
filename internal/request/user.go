package request

type UserSearchRequest struct {
	Nickname string `json:"nickname" form:"nickname"`
	Status   int    `json:"status" form:"status"`
	PageRequest
}
