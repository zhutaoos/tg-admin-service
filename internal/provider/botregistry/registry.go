package botregistry

import (
    "app/internal/model"
    "app/tools/logger"
    "context"
    "strings"

    "gorm.io/gorm"
)

// Registry 基于DB的简单实现：从 bot_config 中读取该群的机器人配置，返回其 Token 作为候选。
// 后续如需“群-多Bot”支持，可改为按一对多表查询。
type Registry struct{
    db *gorm.DB
}

func NewRegistry(db *gorm.DB) *Registry { return &Registry{db: db} }

func (r *Registry) Candidates(ctx context.Context, chatID int64) ([]string, error) {
    if r.db == nil { return nil, nil }
    var bc model.BotConfig
    if err := r.db.WithContext(ctx).Where("group_id = ?", chatID).First(&bc).Error; err != nil {
        return nil, nil // 无bot配置时返回空
    }
    // 从Config(JSON)中提取 token。由于当前没有专门的结构体解析，这里做个简单字符串查找。
    // 若存在正式DTO，可改为 json 解析。
    s := string(bc.Config)
    token := extractTokenQuick(s)
    if token == "" { return nil, nil }
    return []string{token}, nil
}

// 轻量级提取 token（避免额外结构定义）。
func extractTokenQuick(s string) string {
    // 查找 "token":"..."
    i := strings.Index(s, "\"token\"")
    if i < 0 { return "" }
    s = s[i+7:]
    i = strings.Index(s, ":")
    if i < 0 { return "" }
    s = s[i+1:]
    // 跳过空白与引号
    s = strings.TrimLeft(s, " \t\r\n")
    if len(s) == 0 || s[0] != '"' { return "" }
    s = s[1:]
    j := strings.IndexByte(s, '"')
    if j <= 0 { return "" }
    token := s[:j]
    if token == "" { logger.System("BotRegistry: token为空") }
    return token
}
