package dto

// BotConfigData represents the structure of bot configuration data
// This struct is used for JSON marshaling/unmarshaling the Config field
// 机器人配置数据结构
type BotConfigData struct {
    Name             string `json:"name"`                 // 机器人名称
    Token            string `json:"token"`                // 机器人token
    GroupID          int64  `json:"groupId"`              // 群组ID
    InviteLink       string `json:"inviteLink"`           // 群组邀请链接
    SubscribeChannel string `json:"subscribeChannelLink"` // 订阅频道链接
    GroupNamePrefix  string `json:"groupNamePrefix"`      // 群组名称前缀
}

// SubscribeCheckConfig represents the configuration for subscribe check feature
type SubscribeCheckConfig struct {
    Enabled        bool    `json:"enabled"`
    Channels       []int64 `json:"channels"`
    WelcomeMessage string  `json:"welcome_message"`
}
