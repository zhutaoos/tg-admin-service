package queue

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type TelegramProvider interface {
    Send(ctx context.Context, bot string, chatID int64, payload string) (providerMsgID string, status SendStatus, retryAfterSec int)
}

type BotRegistry interface {
    // 返回该chat可用的候选bot列表（bot标识或token别名）
    Candidates(ctx context.Context, chatID int64) ([]string, error)
}

type SendStatus int

const (
    SendOK SendStatus = iota
    SendRetryable
    SendTooManyRequests
    SendFatal
)

type Worker struct {
    rdb      *redis.Client
    cfg      *Config
    limiter  *Limiter
    tg       TelegramProvider
    registry BotRegistry
    shard    string
    consumer string
}

func NewWorker(rdb *redis.Client, cfg *Config, limiter *Limiter, tg TelegramProvider, registry BotRegistry, shard string) *Worker {
    return &Worker{rdb: rdb, cfg: cfg, limiter: limiter, tg: tg, registry: registry, shard: shard, consumer: genConsumerName()}
}

func (w *Worker) ensureGroup(ctx context.Context) {
    _ = w.rdb.XGroupCreateMkStream(ctx, streamReady(w.shard), consumerGroup(w.shard), "$" ).Err()
}

func (w *Worker) Run(ctx context.Context) error {
    w.ensureGroup(ctx)
    stream := streamReady(w.shard)
    group := consumerGroup(w.shard)
    for {
        select {
        case <-ctx.Done():
            return nil
        default:
        }
        // 拉取
        res, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    group,
            Consumer: w.consumer,
            Streams:  []string{stream, ">"},
            Count:    10,
            Block:    time.Second,
        }).Result()
        if err != nil {
            if errors.Is(err, redis.Nil) { continue }
            // 其他错误短暂休眠
            time.Sleep(200 * time.Millisecond)
            continue
        }
        for _, s := range res {
            for _, msg := range s.Messages {
                w.handleOne(ctx, s.Stream, msg)
            }
        }
    }
}

func (w *Worker) handleOne(ctx context.Context, stream string, msg redis.XMessage) {
    // 解析字段
    var j Job
    j.JID = getString(msg.Values["jid"]) 
    j.Payload = getString(msg.Values["payload"]) 
    j.Idem = getString(msg.Values["idem"]) 
    j.Attempts = getInt(msg.Values["attempts"]) 
    j.ChatID = getInt64(msg.Values["chat_id"]) 
    // 候选bot
    if s := getString(msg.Values["bot_candidates"]); s != "" {
        _ = json.Unmarshal([]byte(s), &j.BotCandidates)
    }
    if len(j.BotCandidates) == 0 && w.registry != nil {
        if cands, err := w.registry.Candidates(ctx, j.ChatID); err == nil { j.BotCandidates = cands }
    }
    now := time.Now()
    // 动态选择可用bot：先按每群间隔与全局速率检查，选第一个可用的
    choose := ""
    var minWaitMs int64 = 1<<62
    for _, b := range j.BotCandidates {
        if ok, wait, _ := w.limiter.CheckPerChatGap(ctx, b, j.ChatID, now); !ok {
            if wait < minWaitMs { minWaitMs = wait }
            continue
        }
        if ok, wait, _ := w.limiter.TryAcquireGlobal(ctx, b, now); !ok {
            if wait < minWaitMs { minWaitMs = wait }
            continue
        }
        choose = b
        break
    }
    // 若无可用bot，延时重投
    if choose == "" {
        // 当无候选或均不可用时，minWaitMs 可能未被更新，此时使用保底回退
        if minWaitMs == 1<<62 || minWaitMs <= 0 { minWaitMs = 500 }
        if minWaitMs < 100 { minWaitMs = 100 }
        j.Attempts++
        score := now.Add(time.Duration(minWaitMs) * time.Millisecond).UnixMilli()
        b, _ := json.Marshal(j)
        _ = w.rdb.ZAdd(ctx, zsetDelayed(w.shard), redis.Z{Score: float64(score), Member: string(b)}).Err()
        _ = w.rdb.XAck(ctx, stream, consumerGroup(w.shard), msg.ID).Err()
        return
    }
    // 发送
    providerMsgID, status, retryAfter := w.tg.Send(ctx, choose, j.ChatID, j.Payload)
    switch status {
    case SendOK:
        // 幂等标记
        if j.Idem != "" {
            _ = w.rdb.SetNX(ctx, keyIdem(j.Idem), providerMsgID, 24*time.Hour).Err()
        }
        // 每群间隔
        w.limiter.SetPerChatGap(ctx, choose, j.ChatID, now.Add(time.Duration(w.cfg.PerChatMinGapMs)*time.Millisecond).UnixMilli())
        _ = w.rdb.XAck(ctx, stream, consumerGroup(w.shard), msg.ID).Err()
    case SendTooManyRequests:
        // 429：按 retry_after 退避
        if retryAfter <= 0 { retryAfter = 1 }
        j.Attempts++
        score := now.Add(time.Duration(retryAfter) * time.Second).UnixMilli()
        b, _ := json.Marshal(j)
        _ = w.rdb.ZAdd(ctx, zsetDelayed(w.shard), redis.Z{Score: float64(score), Member: string(b)}).Err()
        // 也更新chat限流，避免短时间再选此chat
        w.limiter.SetPerChatGap(ctx, choose, j.ChatID, now.Add(time.Duration(retryAfter)*time.Second).UnixMilli())
        _ = w.rdb.XAck(ctx, stream, consumerGroup(w.shard), msg.ID).Err()
    case SendRetryable:
        j.Attempts++
        delay := ComputeBackoff(j.Attempts)
        score := now.Add(delay).UnixMilli()
        b, _ := json.Marshal(j)
        _ = w.rdb.ZAdd(ctx, zsetDelayed(w.shard), redis.Z{Score: float64(score), Member: string(b)}).Err()
        _ = w.rdb.XAck(ctx, stream, consumerGroup(w.shard), msg.ID).Err()
    case SendFatal:
        // 记录失败后ACK（此处仅ACK）
        _ = w.rdb.XAck(ctx, stream, consumerGroup(w.shard), msg.ID).Err()
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
        var n int64
        _, _ = fmtSscanfInt64(x, &n)
        return n
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
    if i < len(s) && s[i] == '-' { sign = -1; i++ }
    for ; i < len(s); i++ {
        c := s[i]
        if c < '0' || c > '9' { break }
        x = x*10 + int(c-'0')
        n++
    }
    *p = sign * x
    return n, nil
}
