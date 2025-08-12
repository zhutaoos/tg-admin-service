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
	Id               int64  `json:"id" binding:"required" validate:"required"`
	Name             string `json:"name" validate:"required"`                 // 机器人名称
	Token            string `json:"token" validate:"required"`                // 机器人token
	GroupID          int64  `json:"groupId" validate:"required"`              // 群组ID
	InviteLink       string `json:"inviteLink" validate:"required"`           // 群组邀请链接
	SubscribeChannel string `json:"subscribeChannelLink" validate:"required"` // 订阅频道链接
	GroupNamePrefix  string `json:"groupNamePrefix" validate:"required"`      // 群组名称前缀
}

type GetBotConfigRequest struct {
	Id int64 `json:"id" binding:"required" validate:"required"`
}

type SearchBotConfigRequest struct {
	GroupIds []string `json:"groupIds" binding:"required"`
	Region   string   `json:"region"`
	PageRequest
}
