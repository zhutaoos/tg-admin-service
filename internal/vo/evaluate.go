package model

import (
	"gorm.io/datatypes"
)

// EvaluateListResponse 评价列表响应
type EvaluateListResponse struct {
	List  []JsEvaluate `json:"list"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

type JsEvaluate struct {
	Id               string         `redis:"id" json:"id" gorm:"primaryKey"`
	GroupID          string         `redis:"group_id" json:"group_id"`                           //群组id
	UserId           int64          `redis:"user_id" json:"user_id" `                            //js用户id
	UserName         string         `redis:"user_name" json:"user_name"`                         //js用户名称
	NickName         string         `redis:"nick_name" json:"nick_name" gorm:"column:nick_name"` //js用户昵称
	EvaluateUserName string         `redis:"evaluate_user_name" json:"evaluate_user_name"`       //评价人用户名称
	EvaluateUserId   int64          `redis:"evaluate_user_id" json:"evaluate_user_id"`           //评价人用户id
	EvaluateNickName string         `redis:"evaluate_nick_name" json:"evaluate_nick_name"`       //评价人昵称
	CjDate           string         `redis:"cj_date" json:"cj_date"`                             //出击日期
	Dj               int            `redis:"dj" json:"dj" gorm:"type:int,column:dj"`             //技师等级
	Rz               int            `redis:"rz" json:"rz"`                                       //人照评分
	Sc               int            `redis:"sc" json:"sc"`                                       //身材评分
	Fw               int            `redis:"fw" json:"fw"`                                       //服务评分
	Td               int            `redis:"td" json:"td"`                                       //态度评分
	Hj               int            `redis:"hj" json:"hj"`                                       //环境评分
	Zb               string         `redis:"zb" json:"zb"`                                       //罩杯大小
	Summary          string         `redis:"summary" json:"summary"`                             //总结
	CjMedia          datatypes.JSON `redis:"cj_media" json:"cj_media" gorm:"type:json"`          //出击图片/视频
	Status           int32          `redis:"status" json:"status"`                               //状态 0 待提交 1 已提交 2 审核通过 3 审核不通过
}

// EvaluateStatResponse 评价统计响应
type EvaluateStatResponse struct {
	TotalCount    int64   `json:"total_count"`
	AvgScore      float64 `json:"avg_score"`
	PendingCount  int64   `json:"pending_count"`
	ApprovedCount int64   `json:"approved_count"`
	RejectedCount int64   `json:"rejected_count"`
}
