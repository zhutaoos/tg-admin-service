package request

type UserSearchRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Nickname string `json:"nickname" form:"nickname"`
	Status   int    `json:"status" form:"status"`
}
