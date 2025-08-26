package vo

import (
	"app/internal/model"
	"time"
)

// TaskVo 任务视图对象
type TaskVo struct {
	ID              uint64                  `json:"id"`
	TaskName        string                  `json:"taskName"`
	Description     string                  `json:"description"`
	Status          int                     `json:"status"`
	StatusText      string                  `json:"statusText"`
	AdminID         uint64                  `json:"adminId"`
	GroupIDs        []uint64                `json:"groupIds"`
	MessageIDs      []uint64                `json:"messageIds"`
	TriggerType     model.TriggerType       `json:"triggerType"`
	TriggerTypeText string                  `json:"triggerTypeText"`
	ScheduleTime    *time.Time              `json:"scheduleTime"`
	CronExpression  string                  `json:"cronExpression"`
	CronPatternType *model.CronPatternType  `json:"cronPatternType"`
	CronConfig      map[string]interface{}  `json:"cronConfig"`
	LastExecutedAt  *time.Time              `json:"lastExecutedAt"`
	NextExecuteAt   *time.Time              `json:"nextExecuteAt"`
	ExecuteCount    int                     `json:"executeCount"`
	RetryCount      int                     `json:"retryCount"`
	MaxRetryCount   int                     `json:"maxRetryCount"`
	ErrorMessage    string                  `json:"errorMessage"`
	CreateTime      time.Time               `json:"createTime"`
	UpdateTime      time.Time               `json:"updateTime"`
}

// TaskListVo 任务列表视图对象
type TaskListVo struct {
	Total int64    `json:"total"`
	List  []TaskVo `json:"list"`
}

// TaskStatsVo 任务统计视图对象
type TaskStatsVo struct {
	TotalCount     int64 `json:"totalCount"`
	PendingCount   int64 `json:"pendingCount"`
	RunningCount   int64 `json:"runningCount"`
	CompletedCount int64 `json:"completedCount"`
	FailedCount    int64 `json:"failedCount"`
}

// GetStatusText 获取状态文本
func (t *TaskVo) GetStatusText() string {
	switch t.Status {
	case 0:
		return "待执行"
	case 1:
		return "执行中"
	case 2:
		return "已完成"
	case 3:
		return "执行失败"
	default:
		return "未知状态"
	}
}

// GetTriggerTypeText 获取触发类型文本
func (t *TaskVo) GetTriggerTypeText() string {
	switch t.TriggerType {
	case model.TriggerTypeSchedule:
		return "定时执行"
	case model.TriggerTypeCron:
		return "周期执行"
	default:
		return "未知类型"
	}
}