package job

import (
	"app/tools/logger"
	"context"
	"encoding/json"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"app/internal/model"
	"app/internal/queue"
	"time"
	"gorm.io/gorm"
)

var (
	BotMsgType = "bot_msg"
)

type BotMsgPayload struct {
    MsgType string `json:"msg_type"`
    Content string `json:"content"`
    TaskID  uint64 `json:"taskId,omitempty"`
    ExpireTime string `json:"expireTime,omitempty"`
}

type BotMsgHandler struct{
    db *gorm.DB
    producer *queue.Producer
    cfg *queue.Config
}

func NewBotMsgHandler(jobService *JobService, db *gorm.DB, producer *queue.Producer, cfg *queue.Config) {
	handler := &BotMsgHandler{db: db, producer: producer, cfg: cfg}
	jobService.RegisterHandler(handler)
}

func (b *BotMsgHandler) TaskType() string {
	return BotMsgType
}

func (b *BotMsgHandler) Process(ctx context.Context, payload []byte) error {
	var botMsg BotMsgPayload
	if err := json.Unmarshal(payload, &botMsg); err != nil {
		logger.Error("BotMsgHandler 反序列化失败", "error", err, "payload", string(payload))
		return err
	}

	// 将任务触发转为入队（背压感知 + 就绪/延迟）
	// 1) 解析群组ID集合
	var groupIDs []int64
	if botMsg.TaskID > 0 && b.db != nil {
		var t model.Task
		if err := b.db.WithContext(ctx).Where("id = ? AND is_delete = 0", botMsg.TaskID).First(&t).Error; err == nil {
			var gids []int64
			_ = json.Unmarshal(t.GroupIDs, &gids)
			groupIDs = gids
		}
	}
	if len(groupIDs) == 0 {
		logger.System("BotMsgHandler 未获取到群组列表，跳过入队", "taskId", botMsg.TaskID)
		return nil
	}

	// 2) 组装作业
	now := time.Now()
	jobs := make([]queue.Job, 0, len(groupIDs))
	for idx, gid := range groupIDs {
		idem := buildIdem(botMsg.TaskID, idx, gid, botMsg.Content)
		j := queue.Job{
			JID:         fmt.Sprintf("%d-%d-%d", botMsg.TaskID, gid, now.UnixNano()),
			TaskID:      botMsg.TaskID,
			MsgIdx:      idx,
			ChatID:      gid,
			Payload:     botMsg.Content,
			Idem:        idem,
			Attempts:    0,
			CreatedAtMs: now.UnixMilli(),
		}
		jobs = append(jobs, j)
	}

	if b.producer == nil {
		logger.System("Producer 未初始化，跳过入队")
		return nil
	}
	if err := b.producer.EnqueueJobs(ctx, jobs); err != nil {
		logger.Error("入队失败", "error", err)
		return err
	}
	logger.System("已入队作业", "count", len(jobs), "time", now.Format("2006-01-02 15:04:05"))
	return nil
}

func buildIdem(taskID uint64, idx int, chatID int64, content string) string {
    h := sha1.New()
    h.Write([]byte(fmt.Sprintf("%d|%d|%d|%s", taskID, idx, chatID, content)))
    return hex.EncodeToString(h.Sum(nil))
}
