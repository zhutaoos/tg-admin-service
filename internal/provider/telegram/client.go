package telegram

import (
    "context"
    "app/internal/queue"
)

// Client 是Telegram发送的抽象。当前为桩实现：直接返回成功。
type Client struct{}

func NewClient() *Client { return &Client{} }

// Send 模拟发送：不做网络请求，直接认为成功。
func (c *Client) Send(ctx context.Context, bot string, chatID int64, payload string) (providerMsgID string, status queue.SendStatus, retryAfterSec int) {
    return "mock-msg-id", queue.SendOK, 0
}

