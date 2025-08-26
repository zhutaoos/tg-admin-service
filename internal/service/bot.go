package service

import (
	bizErrors "app/internal/error"
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
func (s *BotService) CreateBotConfig(ctx context.Context, request request.CreateBotConfigRequest, userid uint) error {
	configJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	botConfig := &model.BotConfig{
		AdminId: userid,
		Type:    *request.Type,
		GroupID: request.GroupID,
		Region:  request.Region,
		Config:  configJSON,
	}

	return s.db.WithContext(ctx).Create(botConfig).Error
}

// UpdateBotConfig updates bot configuration by group_id
func (s *BotService) UpdateBotConfig(ctx context.Context, request request.UpdateBotConfigRequest, userid uint) error {
	var botConfig model.BotConfig
	err := s.db.WithContext(ctx).Model(&model.BotConfig{}).Where("id = ?", request.Id).First(&botConfig).Error
	if err != nil {
		return err
	}
	if botConfig.AdminId != userid {
		return bizErrors.ErrInvalidRequest
	}

	// 准备更新的字段
	updates := map[string]interface{}{}

	// 如果请求中包含 BotFeature，则只更新 Features 字段
	if request.BotFeature != nil {
		featuresJSON, err := json.Marshal(request.BotFeature)
		if err != nil {
			return err
		}
		updates["features"] = featuresJSON
		return s.db.WithContext(ctx).Model(&model.BotConfig{}).
			Where("id = ?", request.Id).
			Updates(updates).Error
	}
	return nil
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
func (s *BotService) GetBotConfigData(ctx context.Context, id int64, userId uint) (*vo.BotConfigVo, error) {
	botConfig, err := s.GetBotConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	if botConfig.AdminId != userId {
		return nil, bizErrors.ErrInvalidRequest
	}

	// 解析基础配置数据
	var requestData request.UpdateBotConfigRequest
	err = json.Unmarshal(botConfig.Config, &requestData)
	if err != nil {
		return nil, err
	}

	// 构建响应数据
	configData := &vo.BotConfigVo{
		Id:               botConfig.ID,
		Type:             botConfig.Type,
		Region:           botConfig.Region,
		Name:             requestData.Name,
		Token:            requestData.Token,
		GroupID:          requestData.GroupID,
		InviteLink:       requestData.InviteLink,
		SubscribeChannel: requestData.SubscribeChannel,
		GroupNamePrefix:  requestData.GroupNamePrefix,
	}

	// 解析 BotFeature 数据
	if len(botConfig.Features) > 0 {
		var botFeature vo.BotFeatureVo
		err = json.Unmarshal(botConfig.Features, &botFeature)
		if err != nil {
			logger.Error("解析机器人功能配置数据失败，数据: %s, 错误: %s", string(botConfig.Features), err.Error())
		} else {
			configData.BotFeature = &botFeature
		}
	}

	return configData, nil
}

func (s *BotService) DeleteBotConfig(ctx context.Context, id int64, userId uint) error {
	return s.db.WithContext(ctx).Where("admin_id = ? AND id = ?", userId, id).Delete(&model.BotConfig{}).Error
}

func (s *BotService) SearchBotConfig(ctx context.Context, request request.SearchBotConfigRequest, userId uint) (vo.PageResultVo[vo.BotConfigListVo], error) {
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
	if request.Type != nil {
		query = query.Where("type = ?", request.Type)
	}
	query.Where("admin_id = ?", userId)
	err := query.Count(&total).Error

	if err != nil {
		logger.Error("查询机器人配置列表失败 错误信息: %s", err.Error())
		return vo.PageResultVo[vo.BotConfigListVo]{}, err
	}

	err = query.Order("create_time desc").Offset(request.GetOffset()).Limit(request.Limit).Find(&botConfigs).Error
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

		// 解析 BotFeature 数据
		if len(botConfig.Features) > 0 {
			var botFeature vo.BotFeatureVo
			err := json.Unmarshal(botConfig.Features, &botFeature)
			if err != nil {
				logger.Error("解析机器人功能配置数据失败错误信息: %s", err.Error())
			} else {
				configData.BotFeature = &botFeature
			}
		}

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
