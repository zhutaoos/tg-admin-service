package vo

import (
	"gorm.io/datatypes"
)

type JsEvaluateVo struct {
	Id               string         `redis:"id" json:"id" gorm:"primaryKey"`
	GroupID          string         `redis:"group_id" json:"groupId"`                           //群组id
	UserId           int64          `redis:"user_id" json:"userId" `                            //js用户id
	UserName         string         `redis:"user_name" json:"userName"`                         //js用户名称
	NickName         string         `redis:"nick_name" json:"nickName" gorm:"column:nick_name"` //js用户昵称
	EvaluateUserName string         `redis:"evaluate_user_name" json:"evaluateUserName"`        //评价人用户名称
	EvaluateUserId   int64          `redis:"evaluate_user_id" json:"evaluateUserId"`            //评价人用户id
	EvaluateNickName string         `redis:"evaluate_nick_name" json:"evaluateNickName"`        //评价人昵称
	CjDate           string         `redis:"cj_date" json:"cjDate"`                             //出击日期
	Dj               int            `redis:"dj" json:"dj" gorm:"type:int,column:dj"`            //技师等级
	Rz               int            `redis:"rz" json:"rz"`                                      //人照评分
	Sc               int            `redis:"sc" json:"sc"`                                      //身材评分
	Fw               int            `redis:"fw" json:"fw"`                                      //服务评分
	Td               int            `redis:"td" json:"td"`                                      //态度评分
	Hj               int            `redis:"hj" json:"hj"`                                      //环境评分
	Zb               string         `redis:"zb" json:"zb"`                                      //罩杯大小
	Summary          string         `redis:"summary" json:"summary"`                            //总结
	CjMedia          datatypes.JSON `redis:"cj_media" json:"cjMedia" gorm:"type:json"`          //出击图片/视频
	Status           int32          `redis:"status" json:"status"`                              //状态 0 待提交 1 已提交 2 审核通过 3 审核不通过
}

type MediaVo struct {
	FileID string `json:"file_id"`
	Type   string `json:"type"` // photo, video, audio, document
	URL    string `json:"url"`  // 通过file_id组装的完整URL
}

type JsEvaluateListVo struct {
	Total int64           `json:"total"`
	List  []*JsEvaluateVo `json:"list"`
}
