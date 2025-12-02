package error

import "errors"

var (
	ErrInvalidRequest            error = errors.New("无效的请求")
	ErrRecordNotFound            error = errors.New("记录不存在")
	GroupSendBotNotSupportFeuter error = errors.New("群发机器人不支持配置功能")
)
