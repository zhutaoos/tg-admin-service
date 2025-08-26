package request

import (
	"app/internal/model"
	"time"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskName        string                  `json:"taskName" binding:"required" validate:"required"`
	Description     string                  `json:"description"`
	GroupIDs        []uint64               `json:"groupIds" binding:"required" validate:"required,min=1"`
	MessageIDs      []uint64               `json:"messageIds" binding:"required" validate:"required,min=1"`
	TriggerType     model.TriggerType      `json:"triggerType" binding:"required" validate:"required,oneof=schedule cron"`
	ScheduleTime    *time.Time             `json:"scheduleTime"`
	CronExpression  string                 `json:"cronExpression"`
	CronPatternType *model.CronPatternType `json:"cronPatternType"`
	CronConfig      map[string]interface{} `json:"cronConfig"`
	MaxRetryCount   int                    `json:"maxRetryCount" validate:"min=0,max=10"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	ID              uint64                 `json:"id" binding:"required" validate:"required"`
	TaskName        string                 `json:"taskName" binding:"required" validate:"required"`
	Description     string                 `json:"description"`
	GroupIDs        []uint64               `json:"groupIds" binding:"required" validate:"required,min=1"`
	MessageIDs      []uint64               `json:"messageIds" binding:"required" validate:"required,min=1"`
	TriggerType     model.TriggerType      `json:"triggerType" binding:"required" validate:"required,oneof=schedule cron"`
	ScheduleTime    *time.Time             `json:"scheduleTime"`
	CronExpression  string                 `json:"cronExpression"`
	CronPatternType *model.CronPatternType `json:"cronPatternType"`
	CronConfig      map[string]interface{} `json:"cronConfig"`
	MaxRetryCount   int                    `json:"maxRetryCount" validate:"min=0,max=10"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	PageRequest
	Status    *int     `json:"status" form:"status"`
	GroupID   *uint64  `json:"groupId" form:"groupId"`
	MessageID *uint64  `json:"messageId" form:"messageId"`
	AdminID   *uint64  `json:"adminId" form:"adminId"`
	TaskName  string   `json:"taskName" form:"taskName"`
	StartTime string   `json:"startTime" form:"startTime"`
	EndTime   string   `json:"endTime" form:"endTime"`
}

// DeleteTaskRequest 删除任务请求
type DeleteTaskRequest struct {
	ID uint64 `json:"id" binding:"required" validate:"required"`
}