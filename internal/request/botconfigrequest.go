package request

type CreateBotConfigRequest struct {
	Region           string `json:"region" binding:"required" validate:"required"`
	Name             string `json:"name" validate:"required"`                 // 机器人名称
	Token            string `json:"token" validate:"required"`                // 机器人token
	GroupID          int64  `json:"groupId" validate:"required"`              // 群组ID
	InviteLink       string `json:"inviteLink" validate:"required"`           // 群组邀请链接
	SubscribeChannel string `json:"subscribeChannelLink" validate:"required"` // 订阅频道链接
	GroupNamePrefix  string `json:"groupNamePrefix" validate:"required"`      // 群组名称前缀
}

type UpdateBotConfigRequest struct {
	Id               int64              `json:"id" binding:"required" validate:"required"`
	Name             string             `json:"name" validate:"required"`                 // 机器人名称
	Token            string             `json:"token" validate:"required"`                // 机器人token
	GroupID          int64              `json:"groupId" validate:"required"`              // 群组ID
	InviteLink       string             `json:"inviteLink" validate:"required"`           // 群组邀请链接
	SubscribeChannel string             `json:"subscribeChannelLink" validate:"required"` // 订阅频道链接
	GroupNamePrefix  string             `json:"groupNamePrefix" validate:"required"`      // 群组名称前缀
	BotFeature       *BotFeatureRequest `json:"bot_feature,omitempty"`                    // 机器人功能配置
}

type GetBotConfigRequest struct {
	Id int64 `json:"id" binding:"required" validate:"required"`
}

type SearchBotConfigRequest struct {
	GroupIds []string `json:"groupIds" binding:"required"`
	Region   string   `json:"region"`
	PageRequest
}

// BotFeatureRequest 机器人功能配置请求
type BotFeatureRequest struct {
	Features FeaturesRequest `json:"features" binding:"required"`
	Configs  ConfigsRequest  `json:"configs" binding:"required"`
}

// FeaturesRequest 功能开关配置
type FeaturesRequest struct {
	User UserFeatureRequest `json:"user" binding:"required"`
}

// UserFeatureRequest 用户相关功能开关
type UserFeatureRequest struct {
	Mute      bool `json:"mute"`
	Verify    bool `json:"verify"`
	Subscribe bool `json:"subscribe"`
}

// ConfigsRequest 功能详细配置
type ConfigsRequest struct {
	User UserConfigsRequest `json:"user" binding:"required"`
}

// UserConfigsRequest 用户功能详细配置
type UserConfigsRequest struct {
	Mute      UserMuteConfig      `json:"mute"`
	Verify    UserVerifyConfig    `json:"verify"`
	Subscribe UserSubscribeConfig `json:"subscribe"`
}

// UserMuteConfig 禁言功能配置
type UserMuteConfig struct {
	Enabled bool `json:"enabled"`
}

// UserVerifyConfig 验证功能配置
type UserVerifyConfig struct {
	Enabled bool `json:"enabled"`
}

// UserSubscribeConfig 订阅功能配置
type UserSubscribeConfig struct {
	Enabled    bool            `json:"enabled"`
	ReplyItems []SubscribeItem `json:"replyItems"`
}

// SubscribeItem 订阅项配置
type SubscribeItem struct {
	SubscribeUrl   string `json:"subscribeUrl" binding:"required"`
	ForceSubscribe bool   `json:"forceSubscribe"`
}
