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
	GroupNamePrefix  string `json:"groupNamePrefix"`      // 群组名称前缀
}

// 机器人配置数据结构
type BotConfigListVo struct {
	Id              uint   `json:"id"`
	Region          string `json:"region"`
	Name            string `json:"name"`            // 机器人名称
	GroupID         int64  `json:"groupId"`         // 群组ID
	GroupNamePrefix string `json:"groupNamePrefix"` // 群组名称前缀
	CreateTime      string `json:"createTime"`
}
