package queue

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

// Mover 将到期的ZSET作业搬到Stream
type Mover struct {
    rdb   *redis.Client
    cfg   *Config
    shard string
}

func NewMover(rdb *redis.Client, cfg *Config, shard string) *Mover { return &Mover{rdb: rdb, cfg: cfg, shard: shard} }

func (m *Mover) Run(ctx context.Context) error {
    ticker := time.NewTicker(time.Duration(m.cfg.MoverIntervalMs) * time.Millisecond)
    defer ticker.Stop()
    zdel := zsetDelayed(m.shard)
    stream := streamReady(m.shard)
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            now := time.Now().UnixMilli()
            // 拉取到期条目
            items, err := m.rdb.ZRangeByScore(ctx, zdel, &redis.ZRangeBy{Min: "-inf", Max: fmtI64(now), Offset: 0, Count: int64(m.cfg.MoverBatch)}).Result()
            if err != nil || len(items) == 0 {
                continue
            }
            // 逐个搬运（早期实现，单实例运行即可；需并发/去重时可改Lua）
            for _, s := range items {
                var j Job
                if err := json.Unmarshal([]byte(s), &j); err != nil {
                    // 解析失败则直接删除避免卡住
                    _, _ = m.rdb.ZRem(ctx, zdel, s).Result()
                    continue
                }
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
                if err := m.rdb.XAdd(ctx, &redis.XAddArgs{Stream: stream, Values: fields, Approx: true, MaxLen: m.cfg.StreamMaxLen}).Err(); err != nil {
                    continue
                }
                // 搬运成功后删除ZSET成员
                _, _ = m.rdb.ZRem(ctx, zdel, s).Result()
            }
        }
    }
}

// 小工具：int64 -> string，减少fmt导入
func fmtI64(x int64) string {
    // 简单十进制转换
    if x == 0 { return "0" }
    neg := false
    if x < 0 { neg = true; x = -x }
    buf := [20]byte{}
    i := len(buf)
    for x > 0 {
        i--
        buf[i] = byte('0' + x%10)
        x /= 10
    }
    if neg { i--; buf[i] = '-' }
    return string(buf[i:])
}
