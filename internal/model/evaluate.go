package model

import (
	"app/internal/config"
	"time"

	"gorm.io/datatypes"
)

type JsEvaluateDB struct {
	Id               string         `redis:"id" json:"id" gorm:"primaryKey"`
	GroupID          int64          `redis:"group_id" json:"group_id" gorm:"column:group_id"`                               //群组id
	UserId           int64          `redis:"user_id" json:"user_id" gorm:"column:user_id"`                                  //js用户id
	UserName         string         `redis:"user_name" json:"user_name" gorm:"column:user_name"`                            //js用户名称
	NickName         string         `redis:"nick_name" json:"nick_name" gorm:"column:nick_name"`                            //js用户昵称
	EvaluateUserName string         `redis:"evaluate_user_name" json:"evaluate_user_name" gorm:"column:evaluate_user_name"` //评价人用户名称
	EvaluateUserId   int64          `redis:"evaluate_user_id" json:"evaluate_user_id" gorm:"column:evaluate_user_id"`       //评价人用户id
	EvaluateNickName string         `redis:"evaluate_nick_name" json:"evaluate_nick_name" gorm:"column:evaluate_nick_name"` //评价人昵称
	CjDate           time.Time      `redis:"cj_date" json:"cj_date" gorm:"column:cj_date;type:date"`                        //出击日期
	Dj               int            `redis:"dj" json:"dj" gorm:"type:int,column:dj"`                                        //技师等级
	Rz               int            `redis:"rz" json:"rz" gorm:"column:rz"`                                                 //人照评分
	Sc               int            `redis:"sc" json:"sc" gorm:"column:sc"`                                                 //身材评分
	Fw               int            `redis:"fw" json:"fw" gorm:"column:fw"`                                                 //服务评分
	Td               int            `redis:"td" json:"td" gorm:"column:td"`                                                 //态度评分
	Hj               int            `redis:"hj" json:"hj" gorm:"column:hj"`                                                 //环境评分
	Zb               string         `redis:"zb" json:"zb" gorm:"column:zb"`                                                 //罩杯大小
	Summary          string         `redis:"summary" json:"summary" gorm:"column:summary"`                                  //总结
	CjMedia          datatypes.JSON `redis:"cj_media" json:"cj_media" gorm:"type:json,column:cj_media"`                     //出击图片/视频
	Status           int32          `redis:"status" json:"status" gorm:"column:status"`                                     //状态 0 待提交 1 待审核 2 审核通过 3 审核不通过
}

func (e *JsEvaluateDB) TableName() string {
	return "evaluate"
}

func (e *JsEvaluateDB) List() []JsEvaluateDB {
	list := make([]JsEvaluateDB, 0)
	config.Db().Where(e).Find(&list)
	return list
}
