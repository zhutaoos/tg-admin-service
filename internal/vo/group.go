package vo

import (
	"time"
)

// GroupVo 群组信息VO
type GroupVo struct {
	ID         uint      `json:"id"`
	AdminID    int       `json:"adminId"`
	GroupID    int64     `json:"groupId"`
	GroupName  string    `json:"groupName"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// GroupListVo 群组列表VO
type GroupListVo struct {
	ID         uint   `json:"id"`
	AdminID    int    `json:"adminId"`
	GroupID    int64  `json:"groupId"`
	GroupName  string `json:"groupName"`
	Status     int    `json:"status"`
	CreateTime string `json:"createTime"`
}