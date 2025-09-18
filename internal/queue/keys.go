package queue

import "fmt"

// 统一Redis键名规范

func streamReady(chatID int64) string   { return fmt.Sprintf("tg:send:ready:%d", chatID) }
func zsetDelayed(chatID int64) string   { return fmt.Sprintf("tg:send:delayed:%d", chatID) }
func consumerGroup(chatID int64) string { return fmt.Sprintf("tg:send:cg:%d", chatID) }

func keyBotFixedWindow(bot string, sec int64) string {
	return fmt.Sprintf("tg:lim:botcnt:%s:%d", bot, sec)
}
func keyIdem(idem string) string { return fmt.Sprintf("tg:idem:%s", idem) }
func keyFailureState(bot string, chatID int64) string {
	return fmt.Sprintf("tg:fail:chat:%s:%d", bot, chatID)
}
