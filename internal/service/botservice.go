package service

import (
	"context"
	"encoding/json"

	"app/internal/model"

	"gorm.io/gorm"
)

type BotService struct {
	db *gorm.DB
}

func NewBotService(db *gorm.DB) *BotService {
	return &BotService{db: db}
}

// Bot Config Related Methods

// CreateBotConfig creates a new bot configuration
func (s *BotService) CreateBotConfig(ctx context.Context, groupID int64, configData *model.BotConfigData) error {
	configJSON, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	botConfig := &model.BotConfig{
		GroupID: groupID,
		Config:  configJSON,
	}

	return s.db.WithContext(ctx).Create(botConfig).Error
}

// UpdateBotConfig updates bot configuration by group_id
func (s *BotService) UpdateBotConfig(ctx context.Context, groupID int64, configData *model.BotConfigData) error {
	configJSON, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.BotConfig{}).
		Where("group_id = ?", groupID).
		Update("config", configJSON).Error
}

// GetBotConfig retrieves bot configuration by group_id
func (s *BotService) GetBotConfig(ctx context.Context, groupID int64) (*model.BotConfig, error) {
	var botConfig model.BotConfig
	err := s.db.WithContext(ctx).Where("group_id = ?", groupID).First(&botConfig).Error
	if err != nil {
		return nil, err
	}
	return &botConfig, nil
}

// GetBotConfigData retrieves and parses bot configuration data by group_id
func (s *BotService) GetBotConfigData(ctx context.Context, groupID int64) (*model.BotConfigData, error) {
	botConfig, err := s.GetBotConfig(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var configData model.BotConfigData
	err = json.Unmarshal(botConfig.Config, &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

// Bot Features Related Methods

// CreateBotFeature creates a new bot feature configuration
func (s *BotService) CreateBotFeature(ctx context.Context, groupID int64, featureName string, enabled bool, config interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	botFeature := &model.BotFeatures{
		GroupID:     groupID,
		FeatureName: featureName,
		Enabled:     enabled,
		Config:      configJSON,
	}

	return s.db.WithContext(ctx).Create(botFeature).Error
}

// UpdateBotFeature updates bot feature configuration
func (s *BotService) UpdateBotFeature(ctx context.Context, groupID int64, featureName string, enabled bool, config interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.BotFeatures{}).
		Where("group_id = ? AND feature_name = ?", groupID, featureName).
		Updates(map[string]interface{}{
			"enabled": enabled,
			"config":  configJSON,
		}).Error
}

// GetBotFeature retrieves bot feature configuration by group_id and feature_name
func (s *BotService) GetBotFeature(ctx context.Context, groupID int64, featureName string) (*model.BotFeatures, error) {
	var botFeature model.BotFeatures
	err := s.db.WithContext(ctx).
		Where("group_id = ? AND feature_name = ?", groupID, featureName).
		First(&botFeature).Error
	if err != nil {
		return nil, err
	}
	return &botFeature, nil
}

// GetSubscribeCheckConfig retrieves and parses subscribe check configuration
func (s *BotService) GetSubscribeCheckConfig(ctx context.Context, groupID int64) (*model.SubscribeCheckConfig, error) {
	botFeature, err := s.GetBotFeature(ctx, groupID, "subscribe_check")
	if err != nil {
		return nil, err
	}

	var config model.SubscribeCheckConfig
	err = json.Unmarshal(botFeature.Config, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateSubscribeCheckFeature creates subscribe check feature with default configuration
func (s *BotService) CreateSubscribeCheckFeature(ctx context.Context, groupID int64) error {
	config := model.SubscribeCheckConfig{
		Enabled:        true,
		Channels:       []int64{-1002483637578},
		WelcomeMessage: "@%s 欢迎入裙,请遵守群规,订阅上新频道，未订阅不能发言！",
	}

	return s.CreateBotFeature(ctx, groupID, "subscribe_check", true, config)
}

// ListBotFeaturesByGroup lists all features for a group
func (s *BotService) ListBotFeaturesByGroup(ctx context.Context, groupID int64) ([]model.BotFeatures, error) {
	var features []model.BotFeatures
	err := s.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		Find(&features).Error
	return features, err
}