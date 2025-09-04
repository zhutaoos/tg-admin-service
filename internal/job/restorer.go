package job

import (
    "app/internal/model"
    "app/tools/logger"
    "context"
    "fmt"
    "time"

    "github.com/hibiken/asynq"
    "go.uber.org/fx"
    "gorm.io/gorm"
)

// NewTaskRestorer 在服务启动时恢复调度中的任务（cron 与未来 schedule）
func NewTaskRestorer(db *gorm.DB, js *JobService, lc fx.Lifecycle) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            // 等待 JobService 启动 Worker/Scheduler
            go func() {
                time.Sleep(800 * time.Millisecond)
                if err := restoreTasks(db, js); err != nil {
                    logger.Error("恢复任务失败", "error", err)
                }
            }()
            return nil
        },
    })
}

func restoreTasks(db *gorm.DB, js *JobService) error {
    logger.System("开始恢复未完成定时任务…", "time", time.Now().Format("2006-01-02 15:04:05"))

    // 1) 恢复 cron 周期任务（只要表达式合法就尝试注册；避免重复注册）
    var cronTasks []model.Task
    if err := db.Where("trigger_type = ? AND cron_expression <> '' AND is_delete = 0", model.TriggerTypeCron).Find(&cronTasks).Error; err != nil {
        return fmt.Errorf("查询cron任务失败: %w", err)
    }
    restoredCron := 0
    for _, t := range cronTasks {
        cronExpr := t.CronExpression
        // 跳过明显非法或空表达式
        if cronExpr == "" {
            continue
        }
        // 去重：若已存在相同 spec+taskType 的条目则跳过
        exists, err := schedulerEntryExists(js, cronExpr, BotMsgType)
        if err != nil {
            logger.Error("检查Scheduler条目失败", "error", err, "taskID", t.ID)
            continue
        }
        if exists {
            logger.System("已存在相同cron条目，跳过恢复", "taskID", t.ID, "cron", cronExpr)
            continue
        }

        // 使用结构化payload并带taskId，保持与执行链路一致
        var exp string
        if t.ExpireTime != nil {
            exp = t.ExpireTime.In(time.Local).Format("2006-01-02 15:04:05")
        }
        payload, _ := CreateJSONPayload(BotMsgPayload{MsgType: "cron_restore", Content: fmt.Sprintf("恢复注册-任务ID：%d", t.ID), TaskID: t.ID, ExpireTime: exp})
        if _, err := js.AddCronTask(cronExpr, BotMsgType, payload); err != nil {
            logger.Error("恢复注册cron任务失败", "error", err, "taskID", t.ID, "cron", cronExpr)
            continue
        }
        restoredCron++
    }

    // 2) 只观测 Redis 中的一次性定时任务，避免重复入队
    scheduledCount, earliest, err := inspectRedisScheduled(js)
    if err != nil {
        logger.Error("观测Redis定时任务失败", "error", err)
    } else {
        if earliest.IsZero() {
            logger.System("Redis中待执行的一次性定时任务数量", "count", scheduledCount)
        } else {
            logger.System("Redis中待执行的一次性定时任务数量", "count", scheduledCount, "最早执行时间", earliest.Format("2006-01-02 15:04:05"))
        }
    }

    logger.System("恢复任务完成", "restored_cron", restoredCron)
    return nil
}

func schedulerEntryExists(js *JobService, spec string, taskType string) (bool, error) {
    if js == nil || js.redisConf == nil {
        return false, fmt.Errorf("JobService或Redis配置未初始化")
    }
    redisAddr := fmt.Sprintf("%s:%s", js.redisConf.Ip, js.redisConf.Port)
    redisOpt := asynq.RedisClientOpt{
        Addr:     redisAddr,
        Username: js.redisConf.Username,
        Password: js.redisConf.Password,
        DB:       js.redisConf.Db,
        PoolSize: js.redisConf.MaxTotal,
    }
    inspector := asynq.NewInspector(redisOpt)
    defer inspector.Close()

    entries, err := inspector.SchedulerEntries()
    if err != nil {
        return false, err
    }
    for _, e := range entries {
        if e.Spec == spec && e.Task != nil && e.Task.Type() == taskType {
            return true, nil
        }
    }
    return false, nil
}

// inspectRedisScheduled 观测Redis Scheduled队列（默认队列）
func inspectRedisScheduled(js *JobService) (count int, earliest time.Time, err error) {
    if js == nil || js.redisConf == nil {
        return 0, time.Time{}, fmt.Errorf("JobService或Redis配置未初始化")
    }
    redisAddr := fmt.Sprintf("%s:%s", js.redisConf.Ip, js.redisConf.Port)
    redisOpt := asynq.RedisClientOpt{
        Addr:     redisAddr,
        Username: js.redisConf.Username,
        Password: js.redisConf.Password,
        DB:       js.redisConf.Db,
        PoolSize: js.redisConf.MaxTotal,
    }
    inspector := asynq.NewInspector(redisOpt)
    defer inspector.Close()

    tasks, err := inspector.ListScheduledTasks("default")
    if err != nil {
        return 0, time.Time{}, err
    }
    if len(tasks) == 0 {
        return 0, time.Time{}, nil
    }
    earliest = tasks[0].NextProcessAt
    return len(tasks), earliest, nil
}
