package request

import "github.com/go-playground/validator/v10"

type EvaluateSearchRequest struct {
	PageRequest
	GroupIds         []string `json:"group_ids"`
	Status           int32    `json:"status"`
	EvaluateNickName string   `json:"evaluateNickName"`
}

type EvaluateUpdateParam struct {
	Id      string `json:"id" validate:"required"`
	Dj      int    `json:"dj" validate:"omitempty,min=1,max=3"`         //技师等级
	Rz      int    `json:"rz" validate:"omitempty,min=1,max=10"`        //人照评分
	Sc      int    `json:"sc" validate:"omitempty,min=1,max=10"`        //身材评分
	Fw      int    `json:"fw" validate:"omitempty,min=1,max=10"`        //服务评分
	Td      int    `json:"td" validate:"omitempty,min=1,max=10"`        //态度评分
	Hj      int    `json:"hj" validate:"omitempty,min=1,max=10"`        //环境评分
	Zb      string `json:"zb" validate:"omitempty,oneof=A B C D E F G"` //罩杯大小
	Summary string `json:"summary" validate:"omitempty"`                //总结
	Status  int32  `json:"status" validate:"omitempty,min=1,max=3"`     //状态
}

func (r *EvaluateUpdateParam) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// 专门用于更新的结构体
type EvaluateUpdateFields struct {
	Dj      int    `json:"dj"`
	Rz      int    `json:"rz"`
	Sc      int    `json:"sc"`
	Fw      int    `json:"fw"`
	Td      int    `json:"td"`
	Hj      int    `json:"hj"`
	Zb      string `json:"zb"`
	Summary string `json:"summary"`
	Status  int32  `json:"status"`
}
