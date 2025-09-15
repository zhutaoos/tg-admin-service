package botregistry

import (
    "app/internal/dto"
    "app/internal/model"
    "app/tools/logger"
    "context"
    "encoding/json"
    "strings"

    "gorm.io/gorm"
)

// Registry 基于DB的简单实现：从 bot_config 中读取该群的机器人配置，返回其 Token 作为候选。
// 后续如需“群-多Bot”支持，可改为按一对多表查询。
type Registry struct {
	db *gorm.DB
}

func NewRegistry(db *gorm.DB) *Registry { return &Registry{db: db} }

func (r *Registry) Candidates(ctx context.Context, chatID int64) ([]string, error) {
    if r.db == nil {
        return nil, nil
    }
    var bc model.BotConfig
    if err := r.db.WithContext(ctx).Where("group_id = ? and type = ?", chatID, model.BotTypeBroadcast).First(&bc).Error; err != nil {
        logger.Error("根据 groupId,type 查询 bot_config 失败", "error", err, "groupId", chatID, "type", model.BotTypeBroadcast)
        return nil, nil // 无bot配置时返回空
    }
    // 参考 internal/dto/bot.go 的 BotConfigData 进行 JSON 解析，获取 Token
    var cfg dto.BotConfigData
    if err := json.Unmarshal(bc.Config, &cfg); err != nil {
        logger.Error("解析 bot_config.config 失败", "error", err, "groupId", chatID)
        return nil, nil
    }
    token := strings.TrimSpace(cfg.Token)
    if token == "" {
        logger.System("BotRegistry: token为空", "groupId", chatID)
        return nil, nil
    }
    return []string{token}, nil
}
