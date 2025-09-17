package queue

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	rdb     *redis.Client
	cfg     *Config
	manager *RunnerManager
}

func NewProducer(rdb *redis.Client, cfg *Config, manager *RunnerManager) *Producer {
	return &Producer{rdb: rdb, cfg: cfg, manager: manager}
}

// EnsureGroupFor 确保指定 chat 的消费组存在
func (p *Producer) EnsureGroupFor(ctx context.Context, chatID int64) error {
	stream := streamReady(chatID)
	group := consumerGroup(chatID)
	if err := p.rdb.XGroupCreateMkStream(ctx, stream, group, "$").Err(); err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			return err
		}
	}
	return nil
}

// Backlog 读取指定 chat 的就绪/延迟/待处理数量
func (p *Producer) Backlog(ctx context.Context, chatID int64) (ready, delayed, pending int64) {
	stream := streamReady(chatID)
	zdelayed := zsetDelayed(chatID)
	group := consumerGroup(chatID)

	ready = p.rdb.XLen(ctx, stream).Val()
	delayed = p.rdb.ZCard(ctx, zdelayed).Val()
	if res := p.rdb.XPending(ctx, stream, group); res.Err() == nil {
		pinfo := res.Val()
		pending = pinfo.Count
	}
	return
}

// EnqueueJobs 根据背压策略入队：未超阈值→XADD；超阈值→ZADD 延迟
func (p *Producer) EnqueueJobs(ctx context.Context, jobs []Job) error {
	if len(jobs) == 0 {
		return nil
	}

	buckets := make(map[int64][]Job)
	for _, j := range jobs {
		buckets[j.ChatID] = append(buckets[j.ChatID], j)
	}

	now := time.Now()
	for chatID, items := range buckets {
		if err := p.manager.EnsureRunner(ctx, chatID); err != nil {
			return err
		}
		if err := p.EnsureGroupFor(ctx, chatID); err != nil {
			return err
		}

		ready, delayed, pending := p.Backlog(ctx, chatID)
		backlog := ready + delayed + pending
		cap := int64(p.cfg.GlobalRatePerSec * p.cfg.HorizonSec)
		if cap <= 0 {
			cap = int64(len(items) * 2)
		}
		if gc := int64(len(items) * 2); gc > cap {
			cap = gc
		}

		stream := streamReady(chatID)
		zdelayed := zsetDelayed(chatID)

		if backlog > cap {
			over := backlog - cap
			if over < 0 {
				over = 0
			}
			delaySec := int64(over) / int64(p.cfg.GlobalRatePerSec)
			if delaySec < 1 {
				delaySec = 1
			}
			score := now.Add(time.Duration(delaySec) * time.Second).UnixMilli()
			zs := make([]redis.Z, 0, len(items))
			for _, j := range items {
				b, _ := json.Marshal(j)
				zs = append(zs, redis.Z{Score: float64(score), Member: string(b)})
			}
			if err := p.rdb.ZAdd(ctx, zdelayed, zs...).Err(); err != nil {
				return err
			}
			continue
		}

		for _, j := range items {
			fields := map[string]interface{}{
				"id":       j.ID,
				"chat_id":  j.ChatID,
				"payload":  j.Payload,
				"idem":     j.Idem,
				"attempts": j.Attempts,
			}
			if err := p.rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: fields, Approx: true, MaxLen: p.cfg.StreamMaxLen}).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}
