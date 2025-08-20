package vo

// 机器人配置数据结构
type BotConfigVo struct {
	Id               uint   `json:"id"`
	Region           string `json:"region"`
	Name             string `json:"name"`                 // 机器人名称
	Token            string `json:"token"`                // 机器人token
	GroupID          int64  `json:"groupId"`              // 群组ID
	InviteLink       string `json:"inviteLink"`           // 群组邀请链接
	SubscribeChannel string `json:"subscribeChannelLink"` // 订阅频道链接
	GroupNamePrefix  string        `json:"groupNamePrefix"`      // 群组名称前缀
	BotFeature       *BotFeatureVo `json:"bot_feature,omitempty"` // 机器人功能配置
}

// 机器人配置数据结构
type BotConfigListVo struct {
	Id              uint   `json:"id"`
	Region          string `json:"region"`
	Name            string `json:"name"`            // 机器人名称
	GroupID         int64  `json:"groupId"`         // 群组ID
	GroupNamePrefix string `json:"groupNamePrefix"` // 群组名称前缀
	CreateTime      string        `json:"createTime"`
	BotFeature      *BotFeatureVo `json:"bot_feature,omitempty"` // 机器人功能配置
}

// BotFeatureVo 机器人功能配置响应
type BotFeatureVo struct {
	Features FeaturesVo `json:"features"`
	Configs  ConfigsVo  `json:"configs"`
}

// FeaturesVo 功能开关响应
type FeaturesVo struct {
	User UserFeatureVo `json:"user"`
}

// UserFeatureVo 用户功能开关响应
type UserFeatureVo struct {
	Mute      bool `json:"mute"`
	Verify    bool `json:"verify"`
	Subscribe bool `json:"subscribe"`
}

// ConfigsVo 功能配置响应
type ConfigsVo struct {
	User UserConfigsVo `json:"user"`
}

// UserConfigsVo 用户功能配置响应
type UserConfigsVo struct {
	Mute      UserMuteConfigVo      `json:"mute"`
	Verify    UserVerifyConfigVo    `json:"verify"`
	Subscribe UserSubscribeConfigVo `json:"subscribe"`
}

// UserMuteConfigVo 禁言功能配置响应
type UserMuteConfigVo struct {
	Enabled bool `json:"enabled"`
}

// UserVerifyConfigVo 验证功能配置响应
type UserVerifyConfigVo struct {
	Enabled bool `json:"enabled"`
}

// UserSubscribeConfigVo 订阅功能配置响应
type UserSubscribeConfigVo struct {
	Enabled    bool              `json:"enabled"`
	ReplyItems []SubscribeItemVo `json:"replyItems"`
}

// SubscribeItemVo 订阅项响应
type SubscribeItemVo struct {
	SubscribeUrl   string `json:"subscribeUrl"`
	ForceSubscribe bool   `json:"forceSubscribe"`
}
