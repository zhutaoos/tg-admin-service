package service

import (
    "app/internal/job"
    "app/internal/model"
    "app/internal/request"
    "app/internal/vo"
    "app/tools/cron"
    "app/tools/logger"
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/hibiken/asynq"
    "gorm.io/gorm"
)

// TaskService 任务服务接口
type TaskService interface {
    CreateTask(req *request.CreateTaskRequest, adminID uint) (*vo.TaskVo, error)
    UpdateTask(req *request.UpdateTaskRequest, adminID uint) (*vo.TaskVo, error)
    DeleteTask(req *request.DeleteTaskRequest, adminID uint) error
    GetTaskByID(id uint64, adminID uint) (*vo.TaskVo, error)
    ListTasks(req *request.TaskListRequest, adminID uint) (*vo.TaskListVo, error)
    GetTaskStats(adminID uint) (*vo.TaskStatsVo, error)
    SubmitTask(req *request.SubmitTaskRequest, adminID uint) (*vo.TaskVo, error)
}

type TaskServiceImpl struct {
    db         *gorm.DB
    cronUtils  *cron.CronUtils
    jobService *job.JobService
}

// NewTaskService 创建TaskService实例
func NewTaskService(db *gorm.DB, jobService *job.JobService) TaskService {
    return &TaskServiceImpl{
        db:         db,
        cronUtils:  cron.NewCronUtils(),
        jobService: jobService,
    }
}

// normalizeCronTo5 将6位（含秒）cron表达式规整为5位（去秒）。
func (t *TaskServiceImpl) normalizeCronTo5(expr string) (string, bool) {
    fields := strings.Fields(expr)
    if len(fields) == 6 {
        normalized := strings.Join(fields[1:], " ")
        logger.System("检测到6字段Cron，保存前转换为5字段", "原始", expr, "转换后", normalized)
        return normalized, true
    }
    return expr, false
}

// CreateTask 创建任务
func (t *TaskServiceImpl) CreateTask(req *request.CreateTaskRequest, adminID uint) (*vo.TaskVo, error) {
	// 参数验证
	if req.TriggerType == model.TriggerTypeSchedule && req.GetScheduleTime() == nil {
		return nil, errors.New("定时执行类型必须指定执行时间")
	}
	if req.TriggerType == model.TriggerTypeCron && req.CronExpression == "" {
		return nil, errors.New("周期执行类型必须指定Cron表达式")
	}

	// 校验任务到期时间（必填）
	if req.GetExpireTime() == nil {
		return nil, errors.New("任务到期时间必填")
	}

    // 验证 Cron 表达式
    if req.TriggerType == model.TriggerTypeCron {
        if valid, errMsg := t.cronUtils.ValidateCronExpression(req.CronExpression); !valid {
            return nil, errors.New("Cron表达式格式错误: " + errMsg)
        }
    }

    // 规整Cron为5位后再保存（仅cron类型保留），schedule类型不保存表达式
    cronExpr := req.CronExpression
    if req.TriggerType == model.TriggerTypeCron && cronExpr != "" {
        if normalized, changed := t.normalizeCronTo5(cronExpr); changed {
            cronExpr = normalized
        }
    } else if req.TriggerType == model.TriggerTypeSchedule {
        cronExpr = ""
    }

	// 构建任务模型
    task := &model.Task{
        TaskName:        req.TaskName,
        Description:     req.Description,
        Status:          -1, // 待提交
        AdminID:         adminID,
        TriggerType:     req.TriggerType,
        ScheduleTime:    req.GetScheduleTime(),
        ExpireTime:      req.GetExpireTime(),
        CronExpression:  cronExpr,
        CronPatternType: req.CronPatternType,
        ExecuteCount:    0,
        RetryCount:      0,
        MaxRetryCount:   req.MaxRetryCount,
        CreateTime:      time.Now(),
        UpdateTime:      time.Now(),
    }

	// 如果未设置最大重试次数，默认为3次
	if task.MaxRetryCount == 0 {
		task.MaxRetryCount = 3
	}

	// 处理JSON字段
	groupIDsJSON, err := json.Marshal(req.GroupIDs)
	if err != nil {
		return nil, errors.New("群组ID序列化失败")
	}
	task.GroupIDs = model.JSON(groupIDsJSON)

	messageIDsJSON, err := json.Marshal(req.MessageIDs)
	if err != nil {
		return nil, errors.New("消息ID序列化失败")
	}
	task.MessageIDs = model.JSON(messageIDsJSON)

	if req.CronConfig != nil {
		cronConfigJSON, err := json.Marshal(req.CronConfig)
		if err != nil {
			return nil, errors.New("Cron配置序列化失败")
		}
		task.CronConfig = model.JSON(cronConfigJSON)
	}

    // 不在创建阶段计算 next_execute_at；改为提交阶段计算

    // 保存任务（创建阶段不入队，待提交后入队）
    if err := t.db.Create(task).Error; err != nil {
        return nil, err
    }
	// 转换为VO
	return t.taskToVO(task), nil
}

// UpdateTask 更新任务
func (t *TaskServiceImpl) UpdateTask(req *request.UpdateTaskRequest, adminID uint) (*vo.TaskVo, error) {
    // 查找任务
    task := &model.Task{}
    if err := t.db.Where("id = ? AND admin_id = ? AND is_delete = 0", req.ID, adminID).First(task).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("任务不存在或无权限操作")
        }
        return nil, err
    }

    // 更新前旧值若需用于校验，可在此读取（当前策略：不做调度同步）

    // 前端默认不允许编辑；若后续允许，仅支持待提交状态编辑
    if task.Status != -1 {
        return nil, errors.New("当前状态不允许编辑（仅待提交可编辑）")
    }

    // 参数验证
    if req.TriggerType == model.TriggerTypeSchedule && req.GetScheduleTime() == nil {
        return nil, errors.New("定时执行类型必须指定执行时间")
    }
    if req.TriggerType == model.TriggerTypeCron && req.CronExpression == "" {
        return nil, errors.New("周期执行类型必须指定Cron表达式")
    }

	// 验证 Cron 表达式
    if req.TriggerType == model.TriggerTypeCron {
        if valid, errMsg := t.cronUtils.ValidateCronExpression(req.CronExpression); !valid {
            return nil, errors.New("Cron表达式格式错误: " + errMsg)
        }
    }

    // 规整Cron为5位（仅cron类型保留），schedule类型不保存表达式
    cronExpr := req.CronExpression
    if req.TriggerType == model.TriggerTypeCron && cronExpr != "" {
        if normalized, changed := t.normalizeCronTo5(cronExpr); changed {
            cronExpr = normalized
        }
    } else if req.TriggerType == model.TriggerTypeSchedule {
        cronExpr = ""
    }

	// 更新字段
    updates := map[string]interface{}{
        "task_name":         req.TaskName,
        "description":       req.Description,
        "trigger_type":      req.TriggerType,
        "schedule_time":     req.GetScheduleTime(),
        "cron_expression":   cronExpr,
        "cron_pattern_type": req.CronPatternType,
        "max_retry_count":   req.MaxRetryCount,
        "update_time":       time.Now(),
    }

    // 处理JSON字段
    groupIDsJSON, err := json.Marshal(req.GroupIDs)
	if err != nil {
		return nil, errors.New("群组ID序列化失败")
	}
	updates["group_ids"] = model.JSON(groupIDsJSON)

	messageIDsJSON, err := json.Marshal(req.MessageIDs)
	if err != nil {
		return nil, errors.New("消息ID序列化失败")
	}
	updates["message_ids"] = model.JSON(messageIDsJSON)

	if req.CronConfig != nil {
		cronConfigJSON, err := json.Marshal(req.CronConfig)
		if err != nil {
			return nil, errors.New("Cron配置序列化失败")
		}
		updates["cron_config"] = model.JSON(cronConfigJSON)
	}

    // 不在更新阶段计算 next_execute_at；改为提交阶段计算

    // 执行更新
    if err := t.db.Model(task).Updates(updates).Error; err != nil {
        return nil, err
    }

	// 重新查询更新后的数据
    if err := t.db.Where("id = ?", req.ID).First(task).Error; err != nil {
        return nil, err
    }

    // 当前策略：创建/更新阶段不入队，提交时统一入队

    return t.taskToVO(task), nil
}

// DeleteTask 删除任务
func (t *TaskServiceImpl) DeleteTask(req *request.DeleteTaskRequest, adminID uint) error {
    // 查找任务
    task := &model.Task{}
    if err := t.db.Where("id = ? AND admin_id = ? AND is_delete = 0", req.ID, adminID).First(task).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("任务不存在或无权限操作")
        }
        return err
    }

    // 允许删除：待提交(-1)、待执行(0) 和 失败(3)
    if task.Status != -1 && task.Status != 0 && task.Status != 3 {
        return errors.New("只有待提交、待执行和失败状态的任务可以删除")
    }

    // 软删除标记
    updates := map[string]interface{}{
        "is_delete":  1,
        "update_time": time.Now(),
    }
    if err := t.db.Model(task).Updates(updates).Error; err != nil {
        return err
    }

    // 同步移除 asynq 中的对应任务
    go func(taskCopy model.Task) {
        defer func() { recover() }()
        if taskCopy.TriggerType == model.TriggerTypeSchedule {
            if err := t.jobService.DeleteScheduledByDBTaskID(taskCopy.ID); err != nil {
                logger.Error("移除一次性定时任务失败", "error", err, "taskID", taskCopy.ID)
            }
        } else if taskCopy.TriggerType == model.TriggerTypeCron && taskCopy.CronExpression != "" {
            if _, err := t.jobService.UnregisterCronByTask(taskCopy.CronExpression, taskCopy.ID); err != nil {
                logger.Error("卸载cron任务失败", "error", err, "taskID", taskCopy.ID)
            }
        }
    }(*task)

    return nil
}

// GetTaskByID 根据ID获取任务详情
func (t *TaskServiceImpl) GetTaskByID(id uint64, adminID uint) (*vo.TaskVo, error) {
    task := &model.Task{}
    if err := t.db.Where("id = ? AND admin_id = ? AND is_delete = 0", id, adminID).First(task).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("任务不存在或无权限查看")
        }
        return nil, err
    }

	return t.taskToVO(task), nil
}

// ListTasks 获取任务列表
func (t *TaskServiceImpl) ListTasks(req *request.TaskListRequest, adminID uint) (*vo.TaskListVo, error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	var tasks []model.Task
	var total int64

    query := t.db.Model(&model.Task{}).Where("admin_id = ? AND is_delete = 0", adminID)

	// 构建查询条件
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.TaskName != "" {
		query = query.Where("task_name LIKE ?", "%"+req.TaskName+"%")
	}
	if req.TriggerType != nil {
		query = query.Where("trigger_type = ?", *req.TriggerType)
	}
	if len(req.GroupIDs) > 0 {
		// 构建群组ID的OR条件，查询任务的group_ids字段中包含任意一个指定群组ID的记录
		groupConditions := make([]string, 0, len(req.GroupIDs))
		groupArgs := make([]interface{}, 0, len(req.GroupIDs))
		for _, groupID := range req.GroupIDs {
			groupConditions = append(groupConditions, "JSON_CONTAINS(group_ids, ?)")
			// 将数值转换为JSON格式字符串
			groupIDJSON, _ := json.Marshal(groupID)
			groupArgs = append(groupArgs, string(groupIDJSON))
		}
		if len(groupConditions) > 0 {
			groupQuery := "(" + strings.Join(groupConditions, " OR ") + ")"
			query = query.Where(groupQuery, groupArgs...)
		}
	}
	if len(req.MessageIDs) > 0 {
		// 构建消息ID的OR条件，查询任务的message_ids字段中包含任意一个指定消息ID的记录
		messageConditions := make([]string, 0, len(req.MessageIDs))
		messageArgs := make([]interface{}, 0, len(req.MessageIDs))
		for _, messageID := range req.MessageIDs {
			messageConditions = append(messageConditions, "JSON_CONTAINS(message_ids, ?)")
			// 将数值转换为JSON格式字符串
			messageIDJSON, _ := json.Marshal(messageID)
			messageArgs = append(messageArgs, string(messageIDJSON))
		}
		if len(messageConditions) > 0 {
			messageQuery := "(" + strings.Join(messageConditions, " OR ") + ")"
			query = query.Where(messageQuery, messageArgs...)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := req.GetOffset()
	if err := query.Offset(offset).Limit(req.Limit).Order("create_time DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	// 转换为VO列表
	taskVOs := make([]vo.TaskVo, len(tasks))
	for i, task := range tasks {
		taskVOs[i] = *t.taskToVO(&task)
	}

	return &vo.TaskListVo{
		Total: total,
		List:  taskVOs,
	}, nil
}

// GetTaskStats 获取任务统计信息
func (t *TaskServiceImpl) GetTaskStats(adminID uint) (*vo.TaskStatsVo, error) {
    stats := &vo.TaskStatsVo{}

    // 总数
    if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND is_delete = 0", adminID).Count(&stats.TotalCount).Error; err != nil {
        return nil, err
    }

    // 各状态统计（仅未删除）
    if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ? AND is_delete = 0", adminID, 0).Count(&stats.PendingCount).Error; err != nil {
        return nil, err
    }
    if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ? AND is_delete = 0", adminID, 1).Count(&stats.RunningCount).Error; err != nil {
        return nil, err
    }
    if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ? AND is_delete = 0", adminID, 2).Count(&stats.CompletedCount).Error; err != nil {
        return nil, err
    }
    if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ? AND is_delete = 0", adminID, 3).Count(&stats.FailedCount).Error; err != nil {
        return nil, err
    }

	return stats, nil
}

// SubmitTask 提交任务：将待提交(-1)的任务变为待执行(0)并注册到asynq
func (t *TaskServiceImpl) SubmitTask(req *request.SubmitTaskRequest, adminID uint) (*vo.TaskVo, error) {
    // 查找任务
    task := &model.Task{}
    if err := t.db.Where("id = ? AND admin_id = ? AND is_delete = 0", req.ID, adminID).First(task).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("任务不存在或无权限操作")
        }
        return nil, err
    }

    if task.Status != -1 {
        return nil, errors.New("仅待提交状态的任务可提交")
    }

    now := time.Now()
    updates := map[string]interface{}{
        "status":      0, // 待执行
        "update_time": now,
    }

    // 预计算下一次执行时间
    if task.TriggerType == model.TriggerTypeSchedule && task.ScheduleTime != nil {
        updates["next_execute_at"] = task.ScheduleTime
    } else if task.TriggerType == model.TriggerTypeCron && task.CronExpression != "" {
        next, err := t.cronUtils.CalculateNextExecution(task.CronExpression, now)
        if err != nil {
            return nil, errors.New("计算下次执行时间失败: " + err.Error())
        }
        updates["next_execute_at"] = next
    }

    // 提交时校验到期时间
    if task.ExpireTime == nil {
        return nil, errors.New("任务到期时间必填")
    }
    if task.ExpireTime.Before(now) || task.ExpireTime.Equal(now) {
        return nil, errors.New("任务到期时间已过期")
    }
    if task.TriggerType == model.TriggerTypeSchedule && task.ScheduleTime != nil {
        if !task.ExpireTime.After(*task.ScheduleTime) {
            return nil, errors.New("到期时间必须晚于执行时间")
        }
    }
    if task.TriggerType == model.TriggerTypeCron {
        if next, ok := updates["next_execute_at"].(*time.Time); ok && next != nil {
            if !task.ExpireTime.After(*next) {
                return nil, errors.New("到期时间必须晚于下一次执行时间")
            }
        }
    }

    if err := t.db.Model(task).Updates(updates).Error; err != nil {
        return nil, err
    }

    // 注册到asynq
    var expireStr string
    if task.ExpireTime != nil {
        expireStr = task.ExpireTime.In(time.Local).Format("2006-01-02 15:04:05")
    }
    payload, _ := job.CreateJSONPayload(job.BotMsgPayload{
        MsgType:    "bot_msg",
        Content:    fmt.Sprintf("任务提交，任务ID：%d", task.ID),
        TaskID:     task.ID,
        ExpireTime: expireStr,
    })
    if task.TriggerType == model.TriggerTypeSchedule {
        if task.ScheduleTime == nil {
            return nil, errors.New("定时任务必须设置执行时间")
        }
        if task.ScheduleTime.Before(now) {
            return nil, errors.New("执行时间已过期")
        }
        taskID := fmt.Sprintf("schedule:%d", task.ID)
        if _, err := t.jobService.ScheduleTaskWithID(job.BotMsgType, payload, *task.ScheduleTime, taskID, asynq.MaxRetry(task.MaxRetryCount)); err != nil {
            return nil, fmt.Errorf("注册一次性任务失败: %v", err)
        }
    } else if task.TriggerType == model.TriggerTypeCron {
        if task.CronExpression == "" {
            return nil, errors.New("周期任务必须设置Cron表达式")
        }
        if _, err := t.jobService.AddCronTask(task.CronExpression, job.BotMsgType, payload, asynq.MaxRetry(task.MaxRetryCount)); err != nil {
            return nil, fmt.Errorf("注册周期任务失败: %v", err)
        }
    }

    // 重新查询task
    if err := t.db.Where("id = ?", task.ID).First(task).Error; err != nil {
        return nil, err
    }
    return t.taskToVO(task), nil
}

// taskToVO 将任务模型转换为VO
func (t *TaskServiceImpl) taskToVO(task *model.Task) *vo.TaskVo {
	taskVO := &vo.TaskVo{
		ID:              task.ID,
		TaskName:        task.TaskName,
		Description:     task.Description,
		Status:          task.Status,
		AdminID:         task.AdminID,
		TriggerType:     task.TriggerType,
		CronExpression:  task.CronExpression,
		CronPatternType: task.CronPatternType,
		ExecuteCount:    task.ExecuteCount,
		RetryCount:      task.RetryCount,
		MaxRetryCount:   task.MaxRetryCount,
		ErrorMessage:    task.ErrorMessage,
		CreateTime:      vo.CustomTime{Time: task.CreateTime},
		UpdateTime:      vo.CustomTime{Time: task.UpdateTime},
	}

	// 转换时间字段
    if task.ScheduleTime != nil {
        taskVO.ScheduleTime = &vo.CustomTime{Time: *task.ScheduleTime}
    }
    if task.ExpireTime != nil {
        taskVO.ExpireTime = &vo.CustomTime{Time: *task.ExpireTime}
    }
	if task.LastExecutedAt != nil {
		taskVO.LastExecutedAt = &vo.CustomTime{Time: *task.LastExecutedAt}
	}
	if task.NextExecuteAt != nil {
		taskVO.NextExecuteAt = &vo.CustomTime{Time: *task.NextExecuteAt}
	}

	// 设置状态和类型文本
	taskVO.StatusText = taskVO.GetStatusText()
	taskVO.TriggerTypeText = taskVO.GetTriggerTypeText()

	// 解析JSON字段
	if len(task.GroupIDs) > 0 {
		var groupIDs []int64
		if err := json.Unmarshal(task.GroupIDs, &groupIDs); err == nil {
			taskVO.GroupIDs = groupIDs
		}
	}

	if len(task.MessageIDs) > 0 {
		var messageIDs []uint64
		if err := json.Unmarshal(task.MessageIDs, &messageIDs); err == nil {
			taskVO.MessageIDs = messageIDs
		}
	}

	if len(task.CronConfig) > 0 {
		var cronConfig map[string]interface{}
		if err := json.Unmarshal(task.CronConfig, &cronConfig); err == nil {
			taskVO.CronConfig = cronConfig
		}
	}

	return taskVO
}
