package request

type CreateBotFeatureRequest struct {
	GroupID     int64       `json:"group_id" binding:"required"`
	FeatureName string      `json:"feature_name" binding:"required"`
	Enabled     bool        `json:"enabled"`
	Config      interface{} `json:"config"`
}

type UpdateBotFeatureRequest struct {
	GroupID     int64       `json:"group_id" binding:"required"`
	FeatureName string      `json:"feature_name" binding:"required"`
	Enabled     bool        `json:"enabled"`
	Config      interface{} `json:"config"`
}

type GetBotFeatureRequest struct {
	GroupID     int64  `json:"group_id" binding:"required"`
	FeatureName string `json:"feature_name" binding:"required"`
}

type CreateSubscribeCheckRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type ListBotFeaturesRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
}