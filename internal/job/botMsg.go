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

    now := time.Now()
    // 按群×消息拆分作业（新设计，默认）
    mids := botMsg.MessageIds
    if len(mids) == 0 {
        logger.System("BotMsgHandler 未获取到消息ID列表，跳过入队", "taskId", botMsg.TaskID)
        return nil
    }
    jobs := make([]queue.Job, 0, len(groupIDs)*len(mids))
    for _, gid := range groupIDs {
        for mIdx, mid := range mids {
            idem := buildIdemMessage(botMsg.TaskID, gid, mid)
            mp := map[string]any{
                "taskId":    botMsg.TaskID,
                "messageId": mid,
            }
            if botMsg.MsgType != "" { mp["msgType"] = botMsg.MsgType }
            if botMsg.ExpireTime != "" { mp["expireTime"] = botMsg.ExpireTime }
            bpayload, _ := json.Marshal(mp)
            j := queue.Job{
                JID:         fmt.Sprintf("%d-%d-%d-%d", botMsg.TaskID, gid, mid, now.UnixNano()),
                TaskID:      botMsg.TaskID,
                MsgIdx:      mIdx,
                ChatID:      gid,
                Payload:     string(bpayload),
                Idem:        idem,
                Attempts:    0,
                CreatedAtMs: now.UnixMilli(),
            }
            jobs = append(jobs, j)
        }
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

func buildIdemMessage(taskID uint64, chatID int64, messageID uint64) string {
    h := sha1.New()
    h.Write([]byte(fmt.Sprintf("%d|%d|%d", taskID, chatID, messageID)))
    return hex.EncodeToString(h.Sum(nil))
}
