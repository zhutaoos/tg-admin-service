package queue

import (
	"app/tools/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Mover 将到期的 ZSET 作业搬到 Stream（单 chat 版）
type Mover struct {
	rdb      *redis.Client
	cfg      *Config
	chatID   int64
	onActive func()
}

func NewMover(rdb *redis.Client, cfg *Config, chatID int64, onActive func()) *Mover {
	return &Mover{rdb: rdb, cfg: cfg, chatID: chatID, onActive: onActive}
}

func (m *Mover) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(m.cfg.MoverIntervalMs) * time.Millisecond)
	defer ticker.Stop()
	zdelayed := zsetDelayed(m.chatID)
	stream := streamReady(m.chatID)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			now := time.Now().UnixMilli()
			items, err := m.rdb.ZRangeByScore(ctx, zdelayed, &redis.ZRangeBy{Min: "-inf", Max: fmtI64(now), Offset: 0, Count: int64(m.cfg.MoverBatch)}).Result()
			if err != nil || len(items) == 0 {
				continue
			}
			if m.onActive != nil {
				m.onActive()
			}
			for _, s := range items {
				var j Job
				if err := json.Unmarshal([]byte(s), &j); err != nil {
					logger.Error("mover 解析任务失败", s)
					_, _ = m.rdb.ZRem(ctx, zdelayed, s).Result()
					continue
				}
				fields := map[string]interface{}{
					"id":       j.ID,
					"chat_id":  j.ChatID,
					"payload":  j.Payload,
					"idem":     j.Idem,
					"attempts": j.Attempts,
				}
				if err := m.rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: fields, Approx: true, MaxLen: m.cfg.StreamMaxLen}).Err(); err != nil {
					continue
				}
				_, _ = m.rdb.ZRem(ctx, zdelayed, s).Result()
			}
		}
	}
}

// 小工具：int64 -> string，减少 fmt 导入
func fmtI64(x int64) string {
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
