package model

import (
	"app/tools/logger"
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type TriggerType string
type CronPatternType string

const (
	TriggerTypeSchedule TriggerType = "schedule"
	TriggerTypeCron     TriggerType = "cron"
)

const (
	CronPatternMinute  CronPatternType = "minute"
	CronPatternHour    CronPatternType = "hour"
	CronPatternDaily   CronPatternType = "daily"
	CronPatternWeekly  CronPatternType = "weekly"
	CronPatternMonthly CronPatternType = "monthly"
	CronPatternCustom  CronPatternType = "custom"
)

type Task struct {
	*MysqlBaseModel `gorm:"-:all"`
	ID              uint64           `json:"id" gorm:"primaryKey;type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT;comment:主键ID"`
	TaskName        string           `json:"taskName" gorm:"type:VARCHAR(50) NOT NULL;comment:任务名称"`
	Description     string           `json:"description" gorm:"type:TEXT;comment:任务描述"`
	Status          int              `json:"status" gorm:"type:INT NOT NULL;default:0;comment:任务状态：-1-待提交，0-待执行，1-执行中，2-已完成，3-执行失败"`
	AdminID         uint             `json:"adminId" gorm:"type:BIGINT NOT NULL;comment:创建者ID"`
	GroupIDs        JSON             `json:"groupIds" gorm:"type:JSON NOT NULL;comment:群组ID列表，JSON格式存储"`
	MessageIDs      JSON             `json:"messageIds" gorm:"type:JSON NOT NULL;comment:消息ID列表，JSON格式存储"`
	TriggerType     TriggerType      `json:"triggerType" gorm:"type:ENUM('schedule','cron') NOT NULL;comment:触发类型：schedule-定时执行，cron-周期执行"`
	ScheduleTime    *time.Time       `json:"scheduleTime" gorm:"type:DATETIME;comment:定时执行时间，当trigger_type=schedule时使用"`
	ExpireTime      *time.Time       `json:"expireTime" gorm:"type:DATETIME;comment:任务到期日期"`
	CronExpression  string           `json:"cronExpression" gorm:"type:VARCHAR(100) NOT NULL;comment:Cron表达式，统一存储所有类型的执行规则"`
	CronPatternType *CronPatternType `json:"cronPatternType" gorm:"type:ENUM('minute','hour','daily','weekly','monthly','custom');comment:Cron模式类型，用于编辑时回显"`
	CronConfig      JSON             `json:"cronConfig" gorm:"type:JSON;comment:Cron配置快照，用于编辑时精确回显表单数据"`
	LastExecutedAt  *time.Time       `json:"lastExecutedAt" gorm:"type:DATETIME;comment:上次执行时间"`
	NextExecuteAt   *time.Time       `json:"nextExecuteAt" gorm:"type:DATETIME;comment:下次执行时间，由调度系统计算"`
	ExecuteCount    int              `json:"executeCount" gorm:"type:INT NOT NULL;default:0;comment:已执行次数"`
	RetryCount      int              `json:"retryCount" gorm:"type:INT NOT NULL;default:0;comment:当前重试次数"`
	MaxRetryCount   int              `json:"maxRetryCount" gorm:"type:INT NOT NULL;default:3;comment:最大重试次数"`
	ErrorMessage    string           `json:"errorMessage" gorm:"type:TEXT;comment:错误信息，执行失败时记录"`
	IsDelete        int              `json:"isDelete" gorm:"type:INT NOT NULL DEFAULT 0;comment:是否删除 0:正常 1:删除"`
	CreateTime      time.Time        `json:"createTime" gorm:"type:DATETIME NOT NULL;comment:创建时间"`
	UpdateTime      time.Time        `json:"updateTime" gorm:"type:DATETIME NOT NULL ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// JSON 自定义类型用于处理 JSON 字段
type JSON []byte

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return json.Unmarshal(data, j)
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// UnmarshalTo 快速将 JSON 内容反序列化为指定类型
// 为空或为 null 时不报错，目标保持零值
func (j JSON) UnmarshalTo(dest any) error {
	s := bytes.TrimSpace(j)
	if len(s) == 0 || bytes.EqualFold(s, []byte("null")) {
		return nil
	}
	return json.Unmarshal(s, dest)
}

// Int64s 将 JSON 解码为 []int64
func (j JSON) Int64s() []int64 {
	var v []int64
	err := j.UnmarshalTo(&v)
	if err != nil {
		logger.Error("Int64s 反序列化失败", "error", err, "payload", string(j))
		return nil
	}
	return v
}

// Uint64s 将 JSON 解码为 []uint64
func (j JSON) Uint64s() []uint64 {
	var v []uint64
	err := j.UnmarshalTo(&v)
	if err != nil {
		logger.Error("Uint64s 反序列化失败", "error", err, "payload", string(j))
		return nil
	}
	return v
}

// Strings 将 JSON 解码为 []string
func (j JSON) Strings() ([]string, error) { var v []string; return v, j.UnmarshalTo(&v) }

// Map 将 JSON 解码为 map[string]any
func (j JSON) Map() (map[string]any, error) { var v map[string]any; return v, j.UnmarshalTo(&v) }

// TableName 指定表名
func (Task) TableName() string {
	return "task"
}
