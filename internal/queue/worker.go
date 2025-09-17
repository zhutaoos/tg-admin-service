package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TelegramProvider interface {
	Send(ctx context.Context, bot string, chatID int64, payload string) (providerMsgID string, status SendStatus, retryAfterSec int)
}

type BotRegistry interface {
	// 返回该 chat 可用的候选 bot 列表（bot 标识或 token 别名）
	Candidates(ctx context.Context, chatID int64) ([]string, error)
}

type SendStatus int

const (
	SendOK SendStatus = iota
	SendRetryable
	SendTooManyRequests
	SendFatal
)

// Worker 负责消费单个 chat 的消息流（保持顺序）
type Worker struct {
	rdb      *redis.Client
	cfg      *Config
	limiter  *Limiter
	failure  *FailureTracker
	tg       TelegramProvider
	registry BotRegistry

	chatID   int64
	consumer string
	stream   string
	group    string

	onActive func()

	candidateCache []string
	cacheExpiry    time.Time
}

func NewWorker(rdb *redis.Client, cfg *Config, limiter *Limiter, failure *FailureTracker, tg TelegramProvider, registry BotRegistry, chatID int64, onActive func()) *Worker {
	return &Worker{
		rdb:      rdb,
		cfg:      cfg,
		limiter:  limiter,
		failure:  failure,
		tg:       tg,
		registry: registry,
		chatID:   chatID,
		consumer: genConsumerName(),
		stream:   streamReady(chatID),
		group:    consumerGroup(chatID),
		onActive: onActive,
	}
}

func (w *Worker) ensureGroup(ctx context.Context) {
	_ = w.rdb.XGroupCreateMkStream(ctx, w.stream, w.group, "$").Err()
}

func (w *Worker) Run(ctx context.Context) error {
	w.ensureGroup(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			res, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    w.group,
				Consumer: w.consumer,
				Streams:  []string{w.stream, ">"},
				Count:    32,
				Block:    time.Second,
			}).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				time.Sleep(200 * time.Millisecond)
				continue
			}
			for _, stream := range res {
				for _, msg := range stream.Messages {
					if w.onActive != nil {
						w.onActive()
					}
					opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					w.handleOne(opCtx, msg)
					cancel()
				}
			}
		}
	}
}

func (w *Worker) handleOne(ctx context.Context, msg redis.XMessage) {
	var j Job
	j.ID = getString(msg.Values["id"])
	j.Payload = getString(msg.Values["payload"])
	j.Idem = getString(msg.Values["idem"])
	j.Attempts = getInt(msg.Values["attempts"])
	j.ChatID = getInt64(msg.Values["chat_id"])
	now := time.Now()
	nowMs := now.UnixMilli()

	candidates := w.getCandidates(ctx)
	choose := ""
	var minWaitMs int64 = 1 << 62

	for _, bot := range candidates {
		allow, nextAllowed, blocked, err := w.failure.Allow(ctx, bot, j.ChatID, now)
		if err != nil {
			allow = true
		}
		if blocked {
			w.pruneCandidate(bot, now)
			continue
		}
		if !allow {
			wait := nextAllowed - nowMs
			if wait > 0 && wait < minWaitMs {
				minWaitMs = wait
			}
			continue
		}
		if ok, wait, _ := w.limiter.CheckPerChatGap(ctx, bot, j.ChatID, now); !ok {
			if wait > 0 && wait < minWaitMs {
				minWaitMs = wait
			}
			continue
		}
		if ok, wait, _ := w.limiter.TryAcquireGlobal(ctx, bot, now); !ok {
			if wait > 0 && wait < minWaitMs {
				minWaitMs = wait
			}
			continue
		}
		choose = bot
		break
	}

	if choose == "" {
		if minWaitMs == 1<<62 || minWaitMs <= 0 {
			minWaitMs = 60_000
		}
		if minWaitMs < 100 {
			minWaitMs = 100
		}
		j.Attempts++
		score := now.Add(time.Duration(minWaitMs) * time.Millisecond).UnixMilli()
		body, _ := json.Marshal(j)
		_ = w.rdb.ZAdd(ctx, zsetDelayed(w.chatID), redis.Z{Score: float64(score), Member: string(body)}).Err()
		_ = w.rdb.XAck(ctx, w.stream, w.group, msg.ID).Err()
		return
	}

	providerMsgID, status, retryAfter := w.tg.Send(ctx, choose, j.ChatID, j.Payload)
	switch status {
	case SendOK:
		if j.Idem != "" {
			_ = w.rdb.SetNX(ctx, keyIdem(j.Idem), providerMsgID, 24*time.Hour).Err()
		}
		w.failure.ReportSuccess(ctx, choose, j.ChatID)
		nextGap := now.Add(time.Duration(w.cfg.PerChatMinGapMs) * time.Millisecond).UnixMilli()
		w.limiter.SetPerChatGap(ctx, choose, j.ChatID, nextGap)
		_ = w.rdb.XAck(ctx, w.stream, w.group, msg.ID).Err()
	case SendTooManyRequests:
		if retryAfter <= 0 {
			retryAfter = 1
		}
		j.Attempts++
		delayMs := int64(retryAfter) * 1000
		nextRetry, blocked, _ := w.failure.ReportFailure(ctx, choose, j.ChatID, now)
		if nextRetry > nowMs {
			wait := nextRetry - nowMs
			if wait > delayMs {
				delayMs = wait
			}
		}
		if delayMs < 1000 {
			delayMs = 1000
		}
		score := nowMs + delayMs
		body, _ := json.Marshal(j)
		_ = w.rdb.ZAdd(ctx, zsetDelayed(w.chatID), redis.Z{Score: float64(score), Member: string(body)}).Err()
		w.limiter.SetPerChatGap(ctx, choose, j.ChatID, score)
		if blocked {
			w.pruneCandidate(choose, now)
		}
		_ = w.rdb.XAck(ctx, w.stream, w.group, msg.ID).Err()
	case SendRetryable:
		j.Attempts++
		delay := ComputeBackoff(j.Attempts)
		delayMs := int64(delay / time.Millisecond)
		if delayMs < 500 {
			delayMs = 500
		}
		nextRetry, blocked, _ := w.failure.ReportFailure(ctx, choose, j.ChatID, now)
		if nextRetry > nowMs {
			wait := nextRetry - nowMs
			if wait > delayMs {
				delayMs = wait
			}
		}
		score := nowMs + delayMs
		body, _ := json.Marshal(j)
		_ = w.rdb.ZAdd(ctx, zsetDelayed(w.chatID), redis.Z{Score: float64(score), Member: string(body)}).Err()
		if blocked {
			w.pruneCandidate(choose, now)
		}
		_ = w.rdb.XAck(ctx, w.stream, w.group, msg.ID).Err()
	case SendFatal:
		_, blocked, _ := w.failure.ReportFailure(ctx, choose, j.ChatID, now)
		if blocked {
			w.pruneCandidate(choose, now)
		}
		_ = w.rdb.XAck(ctx, w.stream, w.group, msg.ID).Err()
	}
}

func (w *Worker) getCandidates(ctx context.Context) []string {
	if w.registry == nil {
		return nil
	}
	now := time.Now()
	if len(w.candidateCache) > 0 && now.Before(w.cacheExpiry) {
		return w.candidateCache
	}
	return w.refreshCandidates(ctx, now)
}

func (w *Worker) refreshCandidates(ctx context.Context, now time.Time) []string {
	if w.registry == nil {
		w.candidateCache = nil
		w.cacheExpiry = now.Add(time.Hour)
		return nil
	}
	ttl := time.Duration(w.cfg.CandidateCacheTTLMs) * time.Millisecond
	if ttl <= 0 {
		ttl = time.Hour
	}
	refreshOnEmpty := 5 * time.Minute
	cands, err := w.registry.Candidates(ctx, w.chatID)
	if err != nil {
		if len(w.candidateCache) == 0 {
			w.cacheExpiry = now.Add(refreshOnEmpty)
			return nil
		}
		w.cacheExpiry = now.Add(refreshOnEmpty)
		return w.candidateCache
	}
	if len(cands) == 0 {
		w.candidateCache = nil
		w.cacheExpiry = now.Add(refreshOnEmpty)
		return nil
	}
	w.candidateCache = cands
	w.cacheExpiry = now.Add(ttl)
	return w.candidateCache
}

func (w *Worker) pruneCandidate(bot string, now time.Time) {
	if len(w.candidateCache) == 0 {
		return
	}
	filtered := w.candidateCache[:0]
	for _, b := range w.candidateCache {
		if b != bot {
			filtered = append(filtered, b)
		}
	}
	w.candidateCache = filtered
	shorten := 5 * time.Minute
	limit := now.Add(shorten)
	if w.cacheExpiry.After(limit) {
		w.cacheExpiry = limit
	}
}

// 辅助解析
func getString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return ""
	}
}

func getInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case string:
		var n int
		_, _ = fmtSscanfInt(x, &n)
		return n
	default:
		return 0
	}
}

func getInt64(v any) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case int:
		return int64(x)
	case string:
		if n, err := strconv.ParseInt(x, 10, 64); err == nil {
			return n
		}
		return 0
	default:
		return 0
	}
}

// 无依赖生成消费者名
func genConsumerName() string { return fmt.Sprintf("c-%d", time.Now().UnixNano()) }

// fmt-free int parse helpers
func fmtSscanfInt(s string, p *int) (n int, err error) {
	var x int
	sign := 1
	i := 0
	if i < len(s) && s[i] == '-' {
		sign = -1
		i++
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		x = x*10 + int(c-'0')
		n++
	}
	*p = sign * x
	return n, nil
}
