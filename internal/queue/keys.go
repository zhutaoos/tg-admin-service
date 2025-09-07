package queue

import "fmt"

// 统一Redis键名规范

func streamReady(shard string) string { return fmt.Sprintf("tg:send:ready:%s", shard) }
func zsetDelayed(shard string) string { return fmt.Sprintf("tg:send:delayed:%s", shard) }
func consumerGroup(shard string) string { return fmt.Sprintf("tg:send:cg:%s", shard) }

// 限流与幂等
func keyBotTokenBucket(bot string) string { return fmt.Sprintf("tg:lim:bot:%s", bot) }
func keyBotFixedWindow(bot string, sec int64) string {
    return fmt.Sprintf("tg:lim:botcnt:%s:%d", bot, sec)
}
func keyChatNextAllowed(bot string, chatID int64) string {
    return fmt.Sprintf("tg:lim:chat:%s:%d", bot, chatID)
}
func keyIdem(idem string) string { return fmt.Sprintf("tg:idem:%s", idem) }

// 统计（可选）
func keyStatBotBacklog(bot string) string { return fmt.Sprintf("tg:stat:bot_backlog:%s", bot) }

