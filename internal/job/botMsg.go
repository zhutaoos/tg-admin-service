package job

import (
	"app/tools/logger"
	"context"
	"encoding/json"
	"time"
)

var (
	BotMsgType = "bot_msg"
)

type BotMsgPayload struct {
	MsgType string `json:"msg_type"`
	Content string `json:"content"`
}

type BotMsgHandler struct {
}

func NewBotMsgHandler(jobService *JobService) {
	handler := &BotMsgHandler{}
	jobService.RegisterHandler(handler)
}

func (b *BotMsgHandler) TaskType() string {
	return BotMsgType
}

func (b *BotMsgHandler) Process(ctx context.Context, payload []byte) error {
	logger.System("BotMsgHandler开始处理任务", "执行时间", time.Now().Format("2006-01-02 15:04:05"), "payload", string(payload))

	var botMsg BotMsgPayload
	if err := json.Unmarshal(payload, &botMsg); err != nil {
		logger.Error("BotMsgHandler 反序列化失败", "error", err, "payload", string(payload))
		return err
	}

	logger.System("成功处理机器人消息", "msgType", botMsg.MsgType, "content", botMsg.Content, "处理时间", time.Now().Format("2006-01-02 15:04:05"))

	// 在这里实现具体的消息处理逻辑
	// 例如：发送消息到Telegram、处理回调等

	return nil
}
