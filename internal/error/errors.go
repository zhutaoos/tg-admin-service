package error

import "errors"

var (
	ErrInvalidRequest error = errors.New("无效的请求")
	ErrRecordNotFound error = errors.New("记录不存在")
)
