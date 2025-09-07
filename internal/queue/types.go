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
    Shard            string // 分片，默认 "default"
    GlobalRatePerSec int    // 每个Bot的默认全局速率（固定窗口简化版）
    PerChatMinGapMs  int64  // 每群最小间隔ms
    MoverBatch       int    // 搬运器每批处理数量
    MoverIntervalMs  int    // 搬运器轮询间隔
    HorizonSec       int    // 背压窗口
    StreamMaxLen     int64  // 流最大长度（0表示不限制）
}

func DefaultConfig() *Config {
    return &Config{
        Shard:            "default",
        GlobalRatePerSec: 25,
        PerChatMinGapMs:  1000,
        MoverBatch:       200,
        MoverIntervalMs:  100,
        HorizonSec:       120,
        StreamMaxLen:     0,
    }
}

