package service

import (
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// TaskService 任务服务接口
type TaskService interface {
	CreateTask(req *request.CreateTaskRequest, adminID uint64) (*vo.TaskVo, error)
	UpdateTask(req *request.UpdateTaskRequest, adminID uint64) (*vo.TaskVo, error)
	DeleteTask(req *request.DeleteTaskRequest, adminID uint64) error
	GetTaskByID(id uint64, adminID uint64) (*vo.TaskVo, error)
	ListTasks(req *request.TaskListRequest, adminID uint64) (*vo.TaskListVo, error)
	GetTaskStats(adminID uint64) (*vo.TaskStatsVo, error)
}

type TaskServiceImpl struct {
	db *gorm.DB
}

// NewTaskService 创建TaskService实例
func NewTaskService(db *gorm.DB) TaskService {
	return &TaskServiceImpl{
		db: db,
	}
}

// CreateTask 创建任务
func (t *TaskServiceImpl) CreateTask(req *request.CreateTaskRequest, adminID uint64) (*vo.TaskVo, error) {
	// 参数验证
	if req.TriggerType == model.TriggerTypeSchedule && req.ScheduleTime == nil {
		return nil, errors.New("定时执行类型必须指定执行时间")
	}
	if req.TriggerType == model.TriggerTypeCron && req.CronExpression == "" {
		return nil, errors.New("周期执行类型必须指定Cron表达式")
	}

	// 构建任务模型
	task := &model.Task{
		TaskName:        req.TaskName,
		Description:     req.Description,
		Status:          0, // 待执行
		AdminID:         adminID,
		TriggerType:     req.TriggerType,
		ScheduleTime:    req.ScheduleTime,
		CronExpression:  req.CronExpression,
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

	// 计算下次执行时间
	if req.TriggerType == model.TriggerTypeSchedule && req.ScheduleTime != nil {
		task.NextExecuteAt = req.ScheduleTime
	}
	// TODO: 对于Cron类型，需要根据表达式计算下次执行时间

	// 保存任务
	if err := t.db.Create(task).Error; err != nil {
		return nil, err
	}

	// 转换为VO
	return t.taskToVO(task), nil
}

// UpdateTask 更新任务
func (t *TaskServiceImpl) UpdateTask(req *request.UpdateTaskRequest, adminID uint64) (*vo.TaskVo, error) {
	// 查找任务
	task := &model.Task{}
	if err := t.db.Where("id = ? AND admin_id = ?", req.ID, adminID).First(task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在或无权限操作")
		}
		return nil, err
	}

	// 只有待执行状态的任务可以编辑
	if task.Status != 0 {
		return nil, errors.New("只有待执行状态的任务可以编辑")
	}

	// 参数验证
	if req.TriggerType == model.TriggerTypeSchedule && req.ScheduleTime == nil {
		return nil, errors.New("定时执行类型必须指定执行时间")
	}
	if req.TriggerType == model.TriggerTypeCron && req.CronExpression == "" {
		return nil, errors.New("周期执行类型必须指定Cron表达式")
	}

	// 更新字段
	updates := map[string]interface{}{
		"task_name":         req.TaskName,
		"description":       req.Description,
		"trigger_type":      req.TriggerType,
		"schedule_time":     req.ScheduleTime,
		"cron_expression":   req.CronExpression,
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

	// 更新下次执行时间
	if req.TriggerType == model.TriggerTypeSchedule && req.ScheduleTime != nil {
		updates["next_execute_at"] = req.ScheduleTime
	}

	// 执行更新
	if err := t.db.Model(task).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询更新后的数据
	if err := t.db.Where("id = ?", req.ID).First(task).Error; err != nil {
		return nil, err
	}

	return t.taskToVO(task), nil
}

// DeleteTask 删除任务
func (t *TaskServiceImpl) DeleteTask(req *request.DeleteTaskRequest, adminID uint64) error {
	// 查找任务
	task := &model.Task{}
	if err := t.db.Where("id = ? AND admin_id = ?", req.ID, adminID).First(task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在或无权限操作")
		}
		return err
	}

	// 只有待执行和失败状态的任务可以删除
	if task.Status != 0 && task.Status != 3 {
		return errors.New("只有待执行和失败状态的任务可以删除")
	}

	// 执行删除
	return t.db.Delete(task).Error
}

// GetTaskByID 根据ID获取任务详情
func (t *TaskServiceImpl) GetTaskByID(id uint64, adminID uint64) (*vo.TaskVo, error) {
	task := &model.Task{}
	if err := t.db.Where("id = ? AND admin_id = ?", id, adminID).First(task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在或无权限查看")
		}
		return nil, err
	}

	return t.taskToVO(task), nil
}

// ListTasks 获取任务列表
func (t *TaskServiceImpl) ListTasks(req *request.TaskListRequest, adminID uint64) (*vo.TaskListVo, error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	var tasks []model.Task
	var total int64

	query := t.db.Model(&model.Task{}).Where("admin_id = ?", adminID)

	// 构建查询条件
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.TaskName != "" {
		query = query.Where("task_name LIKE ?", "%"+req.TaskName+"%")
	}
	if req.GroupID != nil {
		query = query.Where("JSON_CONTAINS(group_ids, ?)", *req.GroupID)
	}
	if req.MessageID != nil {
		query = query.Where("JSON_CONTAINS(message_ids, ?)", *req.MessageID)
	}
	if req.StartTime != "" {
		query = query.Where("create_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("create_time <= ?", req.EndTime)
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
func (t *TaskServiceImpl) GetTaskStats(adminID uint64) (*vo.TaskStatsVo, error) {
	stats := &vo.TaskStatsVo{}

	// 总数
	if err := t.db.Model(&model.Task{}).Where("admin_id = ?", adminID).Count(&stats.TotalCount).Error; err != nil {
		return nil, err
	}

	// 各状态统计
	if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ?", adminID, 0).Count(&stats.PendingCount).Error; err != nil {
		return nil, err
	}
	if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ?", adminID, 1).Count(&stats.RunningCount).Error; err != nil {
		return nil, err
	}
	if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ?", adminID, 2).Count(&stats.CompletedCount).Error; err != nil {
		return nil, err
	}
	if err := t.db.Model(&model.Task{}).Where("admin_id = ? AND status = ?", adminID, 3).Count(&stats.FailedCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
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
		ScheduleTime:    task.ScheduleTime,
		CronExpression:  task.CronExpression,
		CronPatternType: task.CronPatternType,
		LastExecutedAt:  task.LastExecutedAt,
		NextExecuteAt:   task.NextExecuteAt,
		ExecuteCount:    task.ExecuteCount,
		RetryCount:      task.RetryCount,
		MaxRetryCount:   task.MaxRetryCount,
		ErrorMessage:    task.ErrorMessage,
		CreateTime:      task.CreateTime,
		UpdateTime:      task.UpdateTime,
	}

	// 设置状态和类型文本
	taskVO.StatusText = taskVO.GetStatusText()
	taskVO.TriggerTypeText = taskVO.GetTriggerTypeText()

	// 解析JSON字段
	if len(task.GroupIDs) > 0 {
		var groupIDs []uint64
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