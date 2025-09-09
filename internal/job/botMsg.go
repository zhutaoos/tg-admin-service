package job

import (
    "app/internal/queue"
    "app/tools/logger"
    "context"
    "crypto/sha1"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"

    "gorm.io/gorm"
)

var (
	BotMsgType = "bot_msg"
)

type BotMsgPayload struct {
	MessageIds []uint64 `json:"messageIds"`
	GroupIds   []int64  `json:"groupIds"`
	MsgType    string   `json:"msg_type"`
	TaskID     uint64   `json:"taskId,omitempty"`
	ExpireTime string   `json:"expireTime,omitempty"`
}

type BotMsgHandler struct {
    producer *queue.Producer
    cfg      *queue.Config
}

func NewBotMsgHandler(jobService *JobService, db *gorm.DB, producer *queue.Producer, cfg *queue.Config) {
    // 当前实现不再依赖 DB 回查
    handler := &BotMsgHandler{producer: producer, cfg: cfg}
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
    // 直接使用 payload 中的 GroupIds（上游已保证传入）
    groupIDs := botMsg.GroupIds
    if len(groupIDs) == 0 {
        logger.System("BotMsgHandler 未获取到群组列表，跳过入队", "taskId", botMsg.TaskID)
        return nil
    }

	// 2) 组装作业
	now := time.Now()
	jobs := make([]queue.Job, 0, len(groupIDs))
    for idx, gid := range groupIDs {
        idem := buildIdem(botMsg.TaskID, idx, gid)
        // 将消息ID等轻量信息放入payload，供下游发送方使用
        mp := map[string]any{"taskId": botMsg.TaskID}
        if len(botMsg.MessageIds) > 0 { mp["messageIds"] = botMsg.MessageIds }
        bpayload, _ := json.Marshal(mp)
        j := queue.Job{
            JID:         fmt.Sprintf("%d-%d-%d", botMsg.TaskID, gid, now.UnixNano()),
            TaskID:      botMsg.TaskID,
            MsgIdx:      idx,
            ChatID:      gid,
            Payload:     string(bpayload),
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

func buildIdem(taskID uint64, idx int, chatID int64) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%d|%d|%d", taskID, idx, chatID)))
	return hex.EncodeToString(h.Sum(nil))
}
