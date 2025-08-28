package request

import (
	"app/internal/model"
	"encoding/json"
	"fmt"
	"time"
)

// FlexibleTime 支持多种时间格式的自定义时间类型
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON反序列化，支持多种时间格式
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	
	if timeStr == "" {
		return nil
	}
	
	// 支持的时间格式列表
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // ISO 8601 with timezone
		"2006-01-02T15:04:05Z",      // ISO 8601 UTC
		"2006-01-02T15:04:05",       // ISO 8601 local
		"2006-01-02 15:04:05",       // 标准格式
		"2006-01-02",                // 仅日期
	}
	
	var err error
	for _, format := range formats {
		if ft.Time, err = time.Parse(format, timeStr); err == nil {
			return nil
		}
	}
	
	return fmt.Errorf("无法解析时间格式: %s", timeStr)
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskName        string                  `json:"taskName" binding:"required" validate:"required"`
	Description     string                  `json:"description"`
	GroupIDs        []int64                `json:"groupIds" binding:"required" validate:"required,min=1"`
	MessageIDs      []uint64               `json:"messageIds" binding:"required" validate:"required,min=1"`
	TriggerType     model.TriggerType      `json:"triggerType" binding:"required" validate:"required,oneof=schedule cron"`
	ScheduleTime    *FlexibleTime          `json:"scheduleTime"`
	CronExpression  string                 `json:"cronExpression"`
	CronPatternType *model.CronPatternType `json:"cronPatternType"`
	CronConfig      map[string]interface{} `json:"cronConfig"`
	MaxRetryCount   int                    `json:"maxRetryCount" validate:"min=0,max=10"`
}

// GetScheduleTime 获取 time.Time 类型的调度时间
func (req *CreateTaskRequest) GetScheduleTime() *time.Time {
	if req.ScheduleTime == nil {
		return nil
	}
	return &req.ScheduleTime.Time
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	ID              uint64                 `json:"id" binding:"required" validate:"required"`
	TaskName        string                 `json:"taskName" binding:"required" validate:"required"`
	Description     string                 `json:"description"`
	GroupIDs        []int64                `json:"groupIds" binding:"required" validate:"required,min=1"`
	MessageIDs      []uint64               `json:"messageIds" binding:"required" validate:"required,min=1"`
	TriggerType     model.TriggerType      `json:"triggerType" binding:"required" validate:"required,oneof=schedule cron"`
	ScheduleTime    *FlexibleTime          `json:"scheduleTime"`
	CronExpression  string                 `json:"cronExpression"`
	CronPatternType *model.CronPatternType `json:"cronPatternType"`
	CronConfig      map[string]interface{} `json:"cronConfig"`
	MaxRetryCount   int                    `json:"maxRetryCount" validate:"min=0,max=10"`
}

// GetScheduleTime 获取 time.Time 类型的调度时间
func (req *UpdateTaskRequest) GetScheduleTime() *time.Time {
	if req.ScheduleTime == nil {
		return nil
	}
	return &req.ScheduleTime.Time
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	PageRequest
	Status    *int     `json:"status" form:"status"`
	GroupID   *int64   `json:"groupId" form:"groupId"`
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

// GetTaskDetailRequest 获取任务详情请求
type GetTaskDetailRequest struct {
	ID uint64 `json:"id" binding:"required" validate:"required"`
}