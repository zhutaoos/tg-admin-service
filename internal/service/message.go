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

// CreateMessage 创建消息 - 多个群组ID拆分为多条记录
func (s *MessageService) CreateMessage(ctx context.Context, req request.CreateMessageRequest, adminID uint) error {
	if adminID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	if len(req.GroupIDs) == 0 {
		return errors.New("群组ID不能为空")
	}

	now := time.Now()
	var messages []model.Message
	
	// 为每个群组ID创建一条记录
	for _, groupID := range req.GroupIDs {
		message := model.Message{
			AdminID:    int(adminID),
			GroupID:    groupID,
			Content:    req.Content,
			Images:     req.Images,
			Medias:     req.Medias,
			CreateTime: &now,
			UpdateTime: &now,
			Status:     0,
		}
		messages = append(messages, message)
	}

	// 批量插入
	return s.db.WithContext(ctx).CreateInBatches(messages, 100).Error
}

// UpdateMessage 更新消息 - 先验证权限再执行逻辑
func (s *MessageService) UpdateMessage(ctx context.Context, req request.UpdateMessageRequest, adminID uint) error {
	if adminID == 0 {
		return errors.New("用户ID不能为空")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先验证记录是否存在且是本人的
		var count int64
		err := tx.Model(&model.Message{}).
			Where("id = ? AND admin_id = ? AND status = 0", req.ID, adminID).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("消息不存在或无权限操作")
		}

		// 2. 获取原有的群组ID列表（已经验证了权限，不需要再加adminid条件）
		var existingGroupIDs []int
		err = tx.Model(&model.Message{}).
			Where("id = ? AND status = 0", req.ID).
			Pluck("group_id", &existingGroupIDs).Error
		if err != nil {
			return err
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

		// 如果没有提供新的群组ID列表，只更新现有记录的内容
		if len(req.GroupIDs) == 0 {
			return tx.Model(&model.Message{}).
				Where("id = ? AND status = 0", req.ID).
				Updates(updates).Error
		}

		// 3. 计算需要删除和新增的群组
		toDelete := make([]int, 0)
		toKeep := make([]int, 0)
		toAdd := make([]int, 0)

		// 找出需要删除的群组ID（存在于原有但不在新列表中）
		for _, existingID := range existingGroupIDs {
			found := false
			for _, newID := range req.GroupIDs {
				if existingID == newID {
					found = true
					break
				}
			}
			if found {
				toKeep = append(toKeep, existingID)
			} else {
				toDelete = append(toDelete, existingID)
			}
		}

		// 找出需要新增的群组ID（存在于新列表但不在原有中）
		for _, newID := range req.GroupIDs {
			found := false
			for _, existingID := range existingGroupIDs {
				if newID == existingID {
					found = true
					break
				}
			}
			if !found {
				toAdd = append(toAdd, newID)
			}
		}

		// 4. 删除不需要的群组记录（软删除）
		if len(toDelete) > 0 {
			err = tx.Model(&model.Message{}).
				Where("id = ? AND group_id IN ?", req.ID, toDelete).
				Update("status", 1).Error
			if err != nil {
				return err
			}
		}

		// 5. 更新保留的群组记录
		if len(toKeep) > 0 {
			err = tx.Model(&model.Message{}).
				Where("id = ? AND group_id IN ? AND status = 0", req.ID, toKeep).
				Updates(updates).Error
			if err != nil {
				return err
			}
		}

		// 6. 为新增的群组创建记录
		if len(toAdd) > 0 {
			// 获取原始消息的创建时间和adminID
			var originalMessage model.Message
			err = tx.Where("id = ? AND status = 0", req.ID).
				First(&originalMessage).Error
			if err != nil {
				return err
			}

			var newMessages []model.Message
			for _, groupID := range toAdd {
				newMessage := model.Message{
					AdminID:    originalMessage.AdminID, // 使用原始的adminID
					GroupID:    groupID,
					Content:    req.Content,
					Images:     req.Images,
					Medias:     req.Medias,
					CreateTime: originalMessage.CreateTime, // 使用原始的创建时间
					UpdateTime: &now,
					Status:     0,
				}
				// 如果没有提供新内容，使用原始内容
				if req.Content == "" {
					newMessage.Content = originalMessage.Content
				}
				if req.Images == nil {
					newMessage.Images = originalMessage.Images
				}
				if req.Medias == nil {
					newMessage.Medias = originalMessage.Medias
				}
				newMessages = append(newMessages, newMessage)
			}

			err = tx.CreateInBatches(newMessages, 100).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
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
		ID:         message.ID,
		GroupID:    message.GroupID,
		Content:    message.Content,
		Images:     message.Images,
		Medias:     message.Medias,
		CreateTime: message.CreateTime,
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
	if req.GroupID != 0 {
		query = query.Where("group_id = ?", req.GroupID)
	}
	if req.Content != "" {
		query = query.Where("content LIKE ?", "%"+req.Content+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		query = query.Where("status = 0") // 默认只查询未删除的消息
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
			ID:         message.ID,
			GroupID:    message.GroupID,
			Content:    message.Content,
			Images:     message.Images,
			Medias:     message.Medias,
			CreateTime: message.CreateTime,
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