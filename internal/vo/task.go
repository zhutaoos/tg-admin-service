package vo

import (
	"app/internal/model"
	"encoding/json"
	"time"
)

// CustomTime 自定义时间类型，用于统一JSON输出格式
type CustomTime struct {
	time.Time
}

// MarshalJSON 自定义时间JSON序列化格式
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	// 可以根据需要调整时间格式
	// "2006-01-02 15:04:05" - 标准格式
	// time.RFC3339 - ISO格式(当前使用)
	formatted := ct.Time.Format("2006-01-02 15:04:05")
	return json.Marshal(formatted)
}

// UnmarshalJSON 自定义时间JSON反序列化
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	
	if timeStr == "" || timeStr == "null" {
		ct.Time = time.Time{}
		return nil
	}
	
	// 支持多种时间格式解析
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
	}
	
	var err error
	for _, format := range formats {
		if ct.Time, err = time.Parse(format, timeStr); err == nil {
			return nil
		}
	}
	
	return err
}

// TaskVo 任务视图对象
type TaskVo struct {
	ID              uint64                  `json:"id"`
	TaskName        string                  `json:"taskName"`
	Description     string                  `json:"description"`
	Status          int                     `json:"status"`
	StatusText      string                  `json:"statusText"`
	AdminID         uint                    `json:"adminId"`
	GroupIDs        []int64                 `json:"groupIds"`
	MessageIDs      []uint64                `json:"messageIds"`
	TriggerType     model.TriggerType       `json:"triggerType"`
	TriggerTypeText string                  `json:"triggerTypeText"`
    ScheduleTime    *CustomTime             `json:"scheduleTime"`
    ExpireTime      *CustomTime             `json:"expireTime"`
    CronExpression  string                  `json:"cronExpression"`
	CronPatternType *model.CronPatternType  `json:"cronPatternType"`
	CronConfig      map[string]interface{}  `json:"cronConfig"`
	LastExecutedAt  *CustomTime             `json:"lastExecutedAt"`
	NextExecuteAt   *CustomTime             `json:"nextExecuteAt"`
	ExecuteCount    int                     `json:"executeCount"`
	RetryCount      int                     `json:"retryCount"`
	MaxRetryCount   int                     `json:"maxRetryCount"`
	ErrorMessage    string                  `json:"errorMessage"`
	CreateTime      CustomTime              `json:"createTime"`
	UpdateTime      CustomTime              `json:"updateTime"`
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
    case -1:
        return "待提交"
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
