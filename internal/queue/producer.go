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

// EnsureGroupFor 确保指定分片的消费组存在
func (p *Producer) EnsureGroupFor(ctx context.Context, shard string) error {
    stream := streamReady(shard)
    group := consumerGroup(shard)
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

// Backlog 读取指定分片的就绪/延迟/待处理（pending）数量
func (p *Producer) Backlog(ctx context.Context, shard string) (ready, delayed, pending int64) {
    stream := streamReady(shard)
    zdel := zsetDelayed(shard)
    group := consumerGroup(shard)

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
    // 按分片分组
    buckets := make(map[string][]Job)
    for _, j := range jobs {
        shard := p.cfg.ShardFor(j.ChatID)
        buckets[shard] = append(buckets[shard], j)
    }

    now := time.Now()
    for shard, items := range buckets {
        // 确保该分片消费组存在
        _ = p.EnsureGroupFor(ctx, shard)

        ready, delayed, pending := p.Backlog(ctx, shard)
        backlog := ready + delayed + pending
        // 分片内阈值估算
        cap := int64(p.cfg.GlobalRatePerSec * p.cfg.HorizonSec)
        uniq := map[int64]struct{}{}
        for _, j := range items { uniq[j.ChatID] = struct{}{} }
        limit := cap
        if gc := int64(len(uniq) * 2); gc > limit { limit = gc }

        stream := streamReady(shard)
        zdel := zsetDelayed(shard)

        if backlog > limit {
            // 计算延迟秒
            over := backlog - limit
            if over < 0 { over = 0 }
            delaySec := int64(over)/int64(p.cfg.GlobalRatePerSec)
            if delaySec < 1 { delaySec = 1 }
            score := now.Add(time.Duration(delaySec) * time.Second).UnixMilli()
            // 批量写入 ZSET（存储为JSON串）
            zs := make([]redis.Z, 0, len(items))
            for _, j := range items {
                b, _ := json.Marshal(j)
                zs = append(zs, redis.Z{Score: float64(score), Member: string(b)})
            }
            if err := p.rdb.ZAdd(ctx, zdel, zs...).Err(); err != nil { return err }
            continue
        }

        // 直接写入Stream
        for _, j := range items {
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
            if err := p.rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: fields, Approx: true, MaxLen: p.cfg.StreamMaxLen}).Err(); err != nil {
                return err
            }
        }
    }
    return nil
}
