package service

import (
	"context"
	"encoding/json"

	"app/internal/dto"
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"app/tools/logger"

	"github.com/jinzhu/copier"
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
func (s *BotService) CreateBotConfig(ctx context.Context, request request.CreateBotConfigRequest) error {
	configJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	botConfig := &model.BotConfig{
		GroupID: request.GroupID,
		Region:  request.Region,
		Config:  configJSON,
	}

	return s.db.WithContext(ctx).Create(botConfig).Error
}

// UpdateBotConfig updates bot configuration by group_id
func (s *BotService) UpdateBotConfig(ctx context.Context, request request.UpdateBotConfigRequest) error {
	configJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.BotConfig{}).
		Where("id = ?", request.Id).
		Update("config", configJSON).Error
}

// GetBotConfig retrieves bot configuration by group_id
func (s *BotService) GetBotConfig(ctx context.Context, id int64) (*model.BotConfig, error) {
	var botConfig model.BotConfig
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&botConfig).Error
	if err != nil {
		return nil, err
	}
	return &botConfig, nil
}

// GetBotConfigData retrieves and parses bot configuration data by group_id
func (s *BotService) GetBotConfigData(ctx context.Context, id int64) (*vo.BotConfigVo, error) {
	botConfig, err := s.GetBotConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	var configData vo.BotConfigVo
	err = json.Unmarshal(botConfig.Config, &configData)
	if err != nil {
		return nil, err
	}
	configData.Id = botConfig.ID
	configData.Region = botConfig.Region
	configData.GroupID = botConfig.GroupID

	return &configData, nil
}

func (s *BotService) SearchBotConfig(ctx context.Context, request request.SearchBotConfigRequest) (vo.PageResultVo[vo.BotConfigListVo], error) {
	var botConfigs []model.BotConfig
	var total int64
	var result []vo.BotConfigListVo

	query := s.db.WithContext(ctx).Model(&model.BotConfig{})
	if len(request.GroupIds) != 0 {
		query = query.Where("group_id in ?", request.GroupIds)
	}
	if request.Region != "" {
		query = query.Where("region = ?", request.Region)
	}
	err := query.Count(&total).Error

	if err != nil {
		logger.Error("查询机器人配置列表失败 错误信息: %s", err.Error())
		return vo.PageResultVo[vo.BotConfigListVo]{}, err
	}

	err = query.Offset(request.GetOffset()).Limit(request.Limit).Find(&botConfigs).Error
	if err != nil {
		logger.Error("查询机器人配置列表失败 错误信息: %s", err.Error())
		return vo.PageResultVo[vo.BotConfigListVo]{}, err
	}

	for _, botConfig := range botConfigs {
		var configData vo.BotConfigListVo
		copier.Copy(&configData, botConfig)

		var cfgData dto.BotConfigData
		err := json.Unmarshal(botConfig.Config, &cfgData)
		if err != nil {
			logger.Error("解析机器人配置数据失败错误信息: %s", err.Error())
		}
		configData.Name = cfgData.Name
		configData.GroupNamePrefix = cfgData.GroupNamePrefix
		configData.CreateTime = botConfig.CreateTime.Format("2006-01-02 15:04:05")

		result = append(result, configData)
	}

	return vo.PageResultVo[vo.BotConfigListVo]{
		Total: total,
		List:  result,
	}, nil
}

// Bot Features Related Methods

// CreateBotFeature creates a new bot feature configuration
// func (s *BotService) CreateBotFeature(ctx context.Context, groupID int64, featureName string, enabled bool, config interface{}) error {
// 	configJSON, err := json.Marshal(config)
// 	if err != nil {
// 		return err
// 	}

// 	botFeature := &model.BotFeatures{
// 		GroupID:     groupID,
// 		FeatureName: featureName,
// 		Enabled:     enabled,
// 		Config:      configJSON,
// 	}

// 	return s.db.WithContext(ctx).Create(botFeature).Error
// }

// // UpdateBotFeature updates bot feature configuration
// func (s *BotService) UpdateBotFeature(ctx context.Context, groupID int64, featureName string, enabled bool, config interface{}) error {
// 	configJSON, err := json.Marshal(config)
// 	if err != nil {
// 		return err
// 	}

// 	return s.db.WithContext(ctx).Model(&model.BotFeatures{}).
// 		Where("group_id = ? AND feature_name = ?", groupID, featureName).
// 		Updates(map[string]interface{}{
// 			"enabled": enabled,
// 			"config":  configJSON,
// 		}).Error
// }

// // GetBotFeature retrieves bot feature configuration by group_id and feature_name
// func (s *BotService) GetBotFeature(ctx context.Context, groupID int64, featureName string) (*model.BotFeatures, error) {
// 	var botFeature model.BotFeatures
// 	err := s.db.WithContext(ctx).
// 		Where("group_id = ? AND feature_name = ?", groupID, featureName).
// 		First(&botFeature).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &botFeature, nil
// }

// // GetSubscribeCheckConfig retrieves and parses subscribe check configuration
// func (s *BotService) GetSubscribeCheckConfig(ctx context.Context, groupID int64) (*model.SubscribeCheckConfig, error) {
// 	botFeature, err := s.GetBotFeature(ctx, groupID, "subscribe_check")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var config model.SubscribeCheckConfig
// 	err = json.Unmarshal(botFeature.Config, &config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &config, nil
// }

// // CreateSubscribeCheckFeature creates subscribe check feature with default configuration
// func (s *BotService) CreateSubscribeCheckFeature(ctx context.Context, groupID int64) error {
// 	config := model.SubscribeCheckConfig{
// 		Enabled:        true,
// 		Channels:       []int64{-1002483637578},
// 		WelcomeMessage: "@%s 欢迎入裙,请遵守群规,订阅上新频道，未订阅不能发言！",
// 	}

// 	return s.CreateBotFeature(ctx, groupID, "subscribe_check", true, config)
// }

// // ListBotFeaturesByGroup lists all features for a group
// func (s *BotService) ListBotFeaturesByGroup(ctx context.Context, groupID int64) ([]model.BotFeatures, error) {
// 	var features []model.BotFeatures
// 	err := s.db.WithContext(ctx).
// 		Where("group_id = ?", groupID).
// 		Find(&features).Error
// 	return features, err
// }
