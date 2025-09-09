package queue

// Job 表示发送作业的载荷（存入Stream/ZSET）
type Job struct {
    JID           string   `json:"jid"`
    TaskID        uint64   `json:"task_id,omitempty"`
    MsgIdx        int      `json:"msg_idx,omitempty"`
    ChatID        int64    `json:"chat_id"`
    Payload       string   `json:"payload"`        // 应用层消息内容（JSON或字符串）
    Idem          string   `json:"idem,omitempty"` // 幂等键
    Attempts      int      `json:"attempts,omitempty"`
    CreatedAtMs   int64    `json:"created_at,omitempty"`
    BotCandidates []string `json:"bot_candidates,omitempty"`
}

// Config 队列与限流配置
type Config struct {
	ShardCount       int   // 分片数量（必填，>=1）。按 chatID 一致性映射到 [0, ShardCount)
	GlobalRatePerSec int   // 每个Bot的默认全局速率（固定窗口简化版）
	PerChatMinGapMs  int64 // 每群最小间隔ms
	MoverBatch       int   // 搬运器每批处理数量
	MoverIntervalMs  int   // 搬运器轮询间隔
	HorizonSec       int   // 背压窗口
	StreamMaxLen     int64 // 流最大长度（0表示不限制）
}

func DefaultConfig() *Config {
	return &Config{
		ShardCount:       16,
		GlobalRatePerSec: 25,
		PerChatMinGapMs:  1000,
		MoverBatch:       200,
		MoverIntervalMs:  100,
		HorizonSec:       120,
		StreamMaxLen:     0,
	}
}

// ShardIndex 基于 chatID 计算分片下标（处理负数ID）
func (c *Config) ShardIndex(chatID int64) int {
	if c == nil || c.ShardCount <= 0 {
		return 0
	}
	// 将可能为负的 chatID 映射为非负数
	u := uint64(chatID)
	// 若 chatID 为负，最高位为1；使用双补码折叠，得到非负数
	u ^= uint64(int64(u) >> 63)
	return int(u % uint64(c.ShardCount))
}

// ShardName 将分片下标格式化为零填充字符串
func (c *Config) ShardName(idx int) string {
	if c == nil || c.ShardCount <= 0 {
		return "00"
	}
	// 计算宽度，例如 16 -> 2，100 -> 3
	width := 1
	n := c.ShardCount - 1
	for n >= 10 {
		width++
		n /= 10
	}
	// 零填充
	s := itoa(idx)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// ShardFor 返回 chatID 对应的分片名
func (c *Config) ShardFor(chatID int64) string {
	return c.ShardName(c.ShardIndex(chatID))
}

// itoa 简易十进制转换（避免fmt热路径）
func itoa(x int) string {
	if x == 0 {
		return "0"
	}
	neg := false
	if x < 0 {
		neg = true
		x = -x
	}
	buf := [20]byte{}
	i := len(buf)
	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
