package request

type PageRequest struct {
	Page  int `json:"page" form:"page" default:"1"`
	Limit int `json:"limit" form:"limit" default:"10"`
}

func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
