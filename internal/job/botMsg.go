package job

import (
	"app/tools/logger"
	"context"
	"encoding/json"
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
	var botMsg BotMsgPayload
	if err := json.Unmarshal(payload, &botMsg); err != nil {
		logger.Error("BotMsgHandler 反序列化失败", "error", err)
		return err
	}

	logger.System("处理机器人消息", "msgType", botMsg.MsgType, "content", botMsg.Content)

	// 在这里实现具体的消息处理逻辑
	// 例如：发送消息到Telegram、处理回调等

	return nil
}
