package queue

// Job 表示发送作业的载荷（存入Stream/ZSET）
type Job struct {
	ID       string `json:"id"` // 队列内部唯一ID（每次入队生成）
	ChatID   int64  `json:"chat_id"`
	Payload  string `json:"payload"`            // 应用层消息内容（JSON或字符串）
	Idem     string `json:"idem,omitempty"`     // 幂等键（用于跨任务去重/记录发送结果）
	Attempts int    `json:"attempts,omitempty"` // 已重试次数（Worker 根据计数计算退避）
}

// Config 队列与限流配置
type Config struct {
	GlobalRatePerSec     int     // 每个Bot的默认全局速率（固定窗口简化版）
	MoverBatch           int     // 搬运器每批处理数量
	MoverIntervalMs      int     // 搬运器轮询间隔
	HorizonSec           int     // 背压窗口
	StreamMaxLen         int64   // 流最大长度（0表示不限制）
	MaxChatRunners       int     // 允许同时活跃的 chat runner 数量（0表示不设上限）
	IdleRunnerTTLMs      int     // chat runner 空闲多久自动回收（毫秒，<=0 表示不回收）
	CandidateCacheTTLMs  int64   // 候选 bot 缓存周期（毫秒）
	FailureBackoffMs     []int64 // 单 bot 失败退避阶梯（毫秒）
	FailureMaxDurationMs int64   // 允许连续失败的最大时长（毫秒），超过后默认阻断
}

func DefaultConfig() *Config {
	return &Config{
		GlobalRatePerSec:    15,
		MoverBatch:          200,
		MoverIntervalMs:     100,
		HorizonSec:          120,
		StreamMaxLen:        0,
		MaxChatRunners:      0,
		IdleRunnerTTLMs:     600000,  // 默认10分钟回收
		CandidateCacheTTLMs: 3600000, // 默认缓存60分钟
		FailureBackoffMs: []int64{
			60_000,
			300_000,
			900_000,
			3_600_000,
			10_800_000,
			21_600_000,
			43_200_000,
			86_400_000,
		},
		FailureMaxDurationMs: 86_400_000, // 最多尝试1天
	}
}
