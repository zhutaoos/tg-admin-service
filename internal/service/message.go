package service

import (
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"app/tools/logger"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessage 创建消息
func (s *MessageService) CreateMessage(ctx context.Context, req request.CreateMessageRequest, adminID uint) error {
	if adminID == 0 {
		return errors.New("用户ID不能为空")
	}

	now := time.Now()
	message := model.Message{
		AdminID:       int(adminID),
		Content:       req.Content,
		Images:        req.Images,
		Medias:        req.Medias,
		AdNickname:    req.AdNickname,
		AdUserID:      req.AdUserID,
		AdGroupLink:   req.AdGroupLink,
		AdChannelLink: req.AdChannelLink,
		CreateTime:    &now,
		UpdateTime:    &now,
		Status:        0,
	}

	return s.db.WithContext(ctx).Create(&message).Error
}

// UpdateMessage 更新消息
func (s *MessageService) UpdateMessage(ctx context.Context, req request.UpdateMessageRequest, adminID uint) error {
	if adminID == 0 {
		return errors.New("用户ID不能为空")
	}

	// 验证记录是否存在且是本人的
	var count int64
	err := s.db.WithContext(ctx).Model(&model.Message{}).
		Where("id = ? AND admin_id = ? AND status = 0", req.ID, adminID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("消息不存在或无权限操作")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"update_time": now,
	}

	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}
	if req.Medias != nil {
		updates["medias"] = req.Medias
	}
	if req.AdNickname != nil {
		updates["ad_nickname"] = req.AdNickname
	}
	if req.AdUserID != nil {
		updates["ad_user_id"] = req.AdUserID
	}
	if req.AdGroupLink != nil {
		updates["ad_group_link"] = req.AdGroupLink
	}
	if req.AdChannelLink != nil {
		updates["ad_channel_link"] = req.AdChannelLink
	}

	return s.db.WithContext(ctx).Model(&model.Message{}).
		Where("id = ? AND status = 0", req.ID).
		Updates(updates).Error
}

// GetMessage 获取消息详情 - 只能查询本人的消息
func (s *MessageService) GetMessage(ctx context.Context, messageID uint, adminID uint) (*vo.MessageVO, error) {
	if adminID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	var message model.Message
	err := s.db.WithContext(ctx).Where("id = ? AND admin_id = ? AND status = 0", messageID, adminID).First(&message).Error
	if err != nil {
		return nil, err
	}

	return &vo.MessageVO{
		ID:            message.ID,
		Content:       message.Content,
		Images:        message.Images,
		Medias:        message.Medias,
		AdNickname:    message.AdNickname,
		AdUserID:      message.AdUserID,
		AdGroupLink:   message.AdGroupLink,
		AdChannelLink: message.AdChannelLink,
		CreateTime:    message.CreateTime,
	}, nil
}

// SearchMessages 搜索消息 - 只能查询本人的消息
func (s *MessageService) SearchMessages(ctx context.Context, req request.SearchMessageRequest, adminID uint) (*vo.PageVO[vo.MessageVO], error) {
	if adminID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	var messages []model.Message
	var total int64

	// 基础查询条件：只查询本人的消息
	query := s.db.WithContext(ctx).Model(&model.Message{}).Where("admin_id = ?", adminID)

	// 添加其他筛选条件
	if req.Content != "" {
		query = query.Where("content LIKE ?", "%"+req.Content+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		query = query.Where("status = 0") // 默认只查询未删除的消息
	}
	// 添加广告字段的搜索条件
	if req.AdNickname != nil && *req.AdNickname != "" {
		query = query.Where("ad_nickname LIKE ?", "%"+*req.AdNickname+"%")
	}
	if req.AdUserID != nil {
		query = query.Where("ad_user_id = ?", *req.AdUserID)
	}
	if req.AdGroupLink != nil && *req.AdGroupLink != "" {
		query = query.Where("ad_group_link LIKE ?", "%"+*req.AdGroupLink+"%")
	}
	if req.AdChannelLink != nil && *req.AdChannelLink != "" {
		query = query.Where("ad_channel_link LIKE ?", "%"+*req.AdChannelLink+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		logger.Error("统计消息数量失败", "err", err.Error())
		return nil, err
	}

	// 分页查询
	offset := req.GetOffset()
	err = query.Offset(offset).Limit(req.Limit).Order("create_time DESC").Find(&messages).Error
	if err != nil {
		logger.Error("查询消息列表失败", "err", err.Error())
		return nil, err
	}

	// 转换为VO
	var messageVOs []vo.MessageVO
	for _, message := range messages {
		messageVOs = append(messageVOs, vo.MessageVO{
			ID:            message.ID,
			Content:       message.Content,
			Images:        message.Images,
			Medias:        message.Medias,
			AdNickname:    message.AdNickname,
			AdUserID:      message.AdUserID,
			AdGroupLink:   message.AdGroupLink,
			AdChannelLink: message.AdChannelLink,
			CreateTime:    message.CreateTime,
		})
	}

	return &vo.PageVO[vo.MessageVO]{
		List:  messageVOs,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// DeleteMessage 删除消息（软删除）- 确保只能删除本人的消息
func (s *MessageService) DeleteMessage(ctx context.Context, messageID uint, adminID uint) error {
	if adminID == 0 {
		return errors.New("用户ID不能为空")
	}

	return s.db.WithContext(ctx).Model(&model.Message{}).
		Where("id = ? AND admin_id = ?", messageID, adminID).
		Update("status", 1).Error
}
