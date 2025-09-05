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

    // 优先解析包含时区的格式
    if t, err := time.Parse("2006-01-02T15:04:05Z07:00", timeStr); err == nil {
        ft.Time = t
        return nil
    }
    if t, err := time.Parse("2006-01-02T15:04:05Z", timeStr); err == nil {
        ft.Time = t
        return nil
    }

    // 其余不含时区的格式，按本地时区解析（main.go 已设置 Asia/Shanghai）
    if t, err := time.ParseInLocation("2006-01-02T15:04:05", timeStr, time.Local); err == nil {
        ft.Time = t
        return nil
    }
    if t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local); err == nil {
        ft.Time = t
        return nil
    }
    if t, err := time.ParseInLocation("2006-01-02", timeStr, time.Local); err == nil {
        ft.Time = t
        return nil
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
    // 定时执行（schedule）不再要求到期时间；仅周期任务（cron）在服务层校验必填
    ExpireTime      *FlexibleTime          `json:"expireTime"`
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

// GetExpireTime 获取 time.Time 类型的到期时间
func (req *CreateTaskRequest) GetExpireTime() *time.Time {
    if req.ExpireTime == nil {
        return nil
    }
    return &req.ExpireTime.Time
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
	Status      *int                `json:"status" form:"status"`
	TaskName    string              `json:"taskName" form:"taskName"`
	TriggerType *model.TriggerType  `json:"triggerType" form:"triggerType"`
	GroupIDs    []int64             `json:"groupIds" form:"groupIds"`
	MessageIDs  []uint64            `json:"messageIds" form:"messageIds"`
}

// DeleteTaskRequest 删除任务请求
type DeleteTaskRequest struct {
    ID uint64 `json:"id" binding:"required" validate:"required"`
}

// GetTaskDetailRequest 获取任务详情请求
type GetTaskDetailRequest struct {
    ID uint64 `json:"id" binding:"required" validate:"required"`
}

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
    ID uint64 `json:"id" binding:"required" validate:"required"`
}
