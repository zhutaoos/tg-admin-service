package queue

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type Producer struct {
    rdb *redis.Client
    cfg *Config
}

func NewProducer(rdb *redis.Client, cfg *Config) *Producer { return &Producer{rdb: rdb, cfg: cfg} }

// EnsureGroup 确保消费组存在
func (p *Producer) EnsureGroup(ctx context.Context) error {
    stream := streamReady(p.cfg.Shard)
    group := consumerGroup(p.cfg.Shard)
    // XGROUP CREATE mkstream
    if err := p.rdb.XGroupCreateMkStream(ctx, stream, group, "$" /* latest */).Err(); err != nil {
        // 忽略已存在的错误
        if err.Error() == "BUSYGROUP Consumer Group name already exists" {
            return nil
        }
        return err
    }
    return nil
}

// Backlog 读取就绪/延迟/待处理（pending）数量
func (p *Producer) Backlog(ctx context.Context) (ready, delayed, pending int64) {
    stream := streamReady(p.cfg.Shard)
    zdel := zsetDelayed(p.cfg.Shard)
    group := consumerGroup(p.cfg.Shard)

    ready = p.rdb.XLen(ctx, stream).Val()
    delayed = p.rdb.ZCard(ctx, zdel).Val()
    // XPENDING，若不存在组则返回0
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
    // 确保组存在（容忍已存在错误）
    _ = p.EnsureGroup(ctx)

    ready, delayed, pending := p.Backlog(ctx)
    backlog := ready + delayed + pending
    // 估算阈值
    cap := int64(p.cfg.GlobalRatePerSec * p.cfg.HorizonSec)
    // group近似：不同chat数量
    uniqChats := map[int64]struct{}{}
    for _, j := range jobs { uniqChats[j.ChatID] = struct{}{} }
    limit := cap
    if gc := int64(len(uniqChats) * 2); gc > limit { limit = gc }

    now := time.Now()
    stream := streamReady(p.cfg.Shard)
    zdel := zsetDelayed(p.cfg.Shard)

    if backlog > limit {
        // 计算延迟秒
        over := backlog - limit
        if over < 0 { over = 0 }
        delaySec := int64(over)/int64(p.cfg.GlobalRatePerSec)
        if delaySec < 1 { delaySec = 1 }
        score := now.Add(time.Duration(delaySec) * time.Second).UnixMilli()
        // 批量写入 ZSET（存储为JSON串）
        zs := make([]*redis.Z, 0, len(jobs))
        for _, j := range jobs {
            b, _ := json.Marshal(j)
            zs = append(zs, &redis.Z{Score: float64(score), Member: string(b)})
        }
        return p.rdb.ZAdd(ctx, zdel, zs...).Err()
    }

    // 直接写入Stream
    for _, j := range jobs {
        fields := map[string]interface{}{
            "jid":        j.JID,
            "task_id":    j.TaskID,
            "msg_idx":    j.MsgIdx,
            "chat_id":    j.ChatID,
            "payload":    j.Payload,
            "idem":       j.Idem,
            "attempts":   j.Attempts,
            "created_at": j.CreatedAtMs,
        }
        if len(j.BotCandidates) > 0 {
            b, _ := json.Marshal(j.BotCandidates)
            fields["bot_candidates"] = string(b)
        }
        // XADD
        if err := p.rdb.XAdd(ctx, &redis.XAddArgs{
            Stream: stream,
            Values: fields,
            // MaxLen 设为近似修剪（0表示不修剪）
            Approx: true,
            MaxLen: p.cfg.StreamMaxLen,
        }).Err(); err != nil {
            return err
        }
    }
    return nil
}
