package cron

import (
    "fmt"
    "time"

    "github.com/robfig/cron/v3"
)

// CronUtils Cron 工具类（仅支持标准 5 位表达式：分 时 日 月 周）
type CronUtils struct {
    parser cron.Parser
}

// NewCronUtils 创建 CronUtils 实例（只解析 5 位表达式）
func NewCronUtils() *CronUtils {
    parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
    return &CronUtils{parser: parser}
}

// CronResult Cron 计算结果
type CronResult struct {
	NextExecuteTime time.Time `json:"nextExecuteTime"`
	IsValid         bool      `json:"isValid"`
	ErrorMessage    string    `json:"errorMessage,omitempty"`
	Description     string    `json:"description,omitempty"`
}

// ParseAndCalculateNext 解析 Cron 表达式并计算下次执行时间
func (c *CronUtils) ParseAndCalculateNext(cronExpr string, baseTime time.Time) *CronResult {
    result := &CronResult{
        IsValid: false,
    }

    if cronExpr == "" {
        result.ErrorMessage = "Cron 表达式不能为空"
        return result
    }

    // 严格使用 5 位解析
    schedule, err := c.parser.Parse(cronExpr)
    if err != nil {
        result.ErrorMessage = fmt.Sprintf("无效的 Cron 表达式: %v", err)
        return result
    }

	// 计算下次执行时间
	nextTime := schedule.Next(baseTime)
	if nextTime.IsZero() {
		result.ErrorMessage = "无法计算下次执行时间"
		return result
	}

	result.NextExecuteTime = nextTime
	result.IsValid = true
	result.Description = c.getDescription(cronExpr)

	return result
}

// CalculateNextExecution 计算下次执行时间（简化版本，仅返回时间）
func (c *CronUtils) CalculateNextExecution(cronExpr string, baseTime time.Time) (*time.Time, error) {
	result := c.ParseAndCalculateNext(cronExpr, baseTime)
	if !result.IsValid {
		return nil, fmt.Errorf(result.ErrorMessage)
	}
	return &result.NextExecuteTime, nil
}

// ValidateCronExpression 验证 Cron 表达式是否有效
func (c *CronUtils) ValidateCronExpression(cronExpr string) (bool, string) {
    if cronExpr == "" {
        return false, "Cron 表达式不能为空"
    }

    // 仅按 5 位解析
    if _, err := c.parser.Parse(cronExpr); err != nil {
        return false, fmt.Sprintf("无效的 Cron 表达式: %v", err)
    }

    return true, ""
}

// getDescription 获取 Cron 表达式的人性化描述
func (c *CronUtils) getDescription(cronExpr string) string {
	// 这里可以根据常见的 cron 表达式模式返回人性化描述
	switch cronExpr {
	case "0 * * * *":
		return "每小时执行一次"
	case "0 0 * * *":
		return "每天凌晨执行一次"
	case "0 0 * * 0":
		return "每周日凌晨执行一次"
	case "0 0 1 * *":
		return "每月1日凌晨执行一次"
	case "*/5 * * * *":
		return "每5分钟执行一次"
	case "*/10 * * * *":
		return "每10分钟执行一次"
	case "*/30 * * * *":
		return "每30分钟执行一次"
	case "0 */2 * * *":
		return "每2小时执行一次"
	case "0 */6 * * *":
		return "每6小时执行一次"
	case "0 */12 * * *":
		return "每12小时执行一次"
	default:
		return "自定义执行周期"
	}
}

// GetNextExecutions 获取多个下次执行时间（用于预览）
func (c *CronUtils) GetNextExecutions(cronExpr string, baseTime time.Time, count int) ([]*time.Time, error) {
    if count <= 0 || count > 10 {
        count = 5 // 默认返回5次
    }

    schedule, err := c.parser.Parse(cronExpr)
    if err != nil {
        return nil, fmt.Errorf("无效的 Cron 表达式: %v", err)
    }

	executions := make([]*time.Time, 0, count)
	currentTime := baseTime

	for i := 0; i < count; i++ {
		nextTime := schedule.Next(currentTime)
		if nextTime.IsZero() {
			break
		}
		executions = append(executions, &nextTime)
		currentTime = nextTime.Add(time.Second) // 加1秒避免重复
	}

	return executions, nil
}

// PresetCronExpressions 预设的常用 Cron 表达式
type PresetCronExpressions struct {
	Every5Minutes  string `json:"every5Minutes"`  // 每5分钟
	Every10Minutes string `json:"every10Minutes"` // 每10分钟
	Every30Minutes string `json:"every30Minutes"` // 每30分钟
	Hourly         string `json:"hourly"`         // 每小时
	Every2Hours    string `json:"every2Hours"`    // 每2小时
	Every6Hours    string `json:"every6Hours"`    // 每6小时
	Every12Hours   string `json:"every12Hours"`   // 每12小时
	Daily          string `json:"daily"`          // 每天
	Weekly         string `json:"weekly"`         // 每周
	Monthly        string `json:"monthly"`        // 每月
}

// GetPresetExpressions 获取预设的 Cron 表达式
func GetPresetExpressions() *PresetCronExpressions {
	return &PresetCronExpressions{
		Every5Minutes:  "*/5 * * * *",
		Every10Minutes: "*/10 * * * *",
		Every30Minutes: "*/30 * * * *",
		Hourly:         "0 * * * *",
		Every2Hours:    "0 */2 * * *",
		Every6Hours:    "0 */6 * * *",
		Every12Hours:   "0 */12 * * *",
		Daily:          "0 0 * * *",
		Weekly:         "0 0 * * 0",
		Monthly:        "0 0 1 * *",
	}
}
