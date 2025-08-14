package service

import (
	"app/internal/model"
	"app/internal/request"
	"app/internal/vo"
	"context"
	"strconv"

	"gorm.io/gorm"
)

// GroupService 群组管理服务接口
type GroupService interface {
	// CreateGroup 创建群组关联
	CreateGroup(ctx context.Context, req request.CreateGroupRequest, currentUserId uint) error
	
	// UpdateGroup 更新群组信息
	UpdateGroup(ctx context.Context, req request.UpdateGroupRequest, currentUserId uint) error
	
	// DeleteGroup 删除群组关联
	DeleteGroup(ctx context.Context, req request.DeleteGroupRequest, currentUserId uint) error
	
	// SearchGroups 分页查询群组列表
	SearchGroups(ctx context.Context, req request.SearchGroupRequest, currentUserId uint) ([]vo.GroupListVo, int64, error)
	
	// GetGroupsByAdminID 获取指定管理员的所有群组
	GetGroupsByAdminID(ctx context.Context, adminID int) ([]model.GroupInfo, error)
	
	// GetGroupsByGroupID 获取指定群组的所有管理员关联
	GetGroupsByGroupID(ctx context.Context, groupID int64) ([]vo.GroupVo, error)
	
	// GetGroupByID 根据ID获取群组信息
	GetGroupByID(ctx context.Context, id int) (*vo.GroupVo, error)
}

// GroupServiceImpl 群组服务实现
type GroupServiceImpl struct {
	db *gorm.DB
}

// NewGroupService 创建GroupService实例
func NewGroupService(db *gorm.DB) GroupService {
	return &GroupServiceImpl{
		db: db,
	}
}

// CreateGroup 创建群组关联
func (s *GroupServiceImpl) CreateGroup(ctx context.Context, req request.CreateGroupRequest, currentUserId uint) error {
	group := &model.Group{
		AdminID:   int(currentUserId),
		GroupID:   req.GroupID,
		GroupName: req.GroupName,
		Status:    0, // 默认正常状态
	}
	
	return s.db.WithContext(ctx).Create(group).Error
}

// UpdateGroup 更新群组信息
func (s *GroupServiceImpl) UpdateGroup(ctx context.Context, req request.UpdateGroupRequest, currentUserId uint) error {
	// 首先检查权限 - 查询该记录的管理员ID
	var group model.Group
	if err := s.db.WithContext(ctx).Where("id = ? AND status = 0", req.ID).First(&group).Error; err != nil {
		return err
	}
	
	// 验证权限
	if group.AdminID != int(currentUserId) {
		return gorm.ErrRecordNotFound // 返回通用错误，避免泄露信息
	}
	
	return s.db.WithContext(ctx).Model(&model.Group{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"group_id":   req.GroupID,
			"group_name": req.GroupName,
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// DeleteGroup 删除群组关联
func (s *GroupServiceImpl) DeleteGroup(ctx context.Context, req request.DeleteGroupRequest, currentUserId uint) error {
	query := s.db.WithContext(ctx).Model(&model.Group{})
	
	if req.ID != nil {
		// 按ID删除 - 先检查权限
		var group model.Group
		if err := s.db.WithContext(ctx).Where("id = ? AND status = 0", *req.ID).First(&group).Error; err != nil {
			return err
		}
		
		// 验证权限
		if group.AdminID != int(currentUserId) {
			return gorm.ErrRecordNotFound
		}
		
		return query.Where("id = ?", *req.ID).Update("status", 1).Error
	}	
	return gorm.ErrRecordNotFound
}

// SearchGroups 分页查询群组列表
func (s *GroupServiceImpl) SearchGroups(ctx context.Context, req request.SearchGroupRequest, currentUserId uint) ([]vo.GroupListVo, int64, error) {
	var groups []*model.Group
	var total int64
	
	query := s.db.WithContext(ctx).Model(&model.Group{}).Where("status = 0")
	
	// 只查询当前用户所属的群组
	query = query.Where("admin_id = ?", currentUserId)
	
	// 按群组ID过滤
	if req.GroupID != nil {
		query = query.Where("group_id = ?", *req.GroupID)
	}
	
	// 按状态过滤
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	
	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := req.GetOffset()
	if req.Limit <= 0 {
		req.Limit = 10
	}
	
	err := query.Order("id DESC").
		Limit(req.Limit).
		Offset(offset).
		Find(&groups).Error
	
	// 转换为 VO 结构
	var groupListVos []vo.GroupListVo
	for _, group := range groups {
		groupListVos = append(groupListVos, vo.GroupListVo{
			ID:         group.ID,
			AdminID:    group.AdminID,
			GroupID:    group.GroupID,
			GroupName:  group.GroupName,
			Status:     group.Status,
			CreateTime: group.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	
	return groupListVos, total, err
}

// GetGroupsByAdminID 获取指定管理员的所有群组
func (s *GroupServiceImpl) GetGroupsByAdminID(ctx context.Context, adminID int) ([]model.GroupInfo, error) {
	var groups []*model.Group
	
	err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("admin_id = ? AND status = 0", adminID).
		Order("id DESC").
		Find(&groups).Error
	
	// 转换为 GroupInfo 结构
	var groupInfos []model.GroupInfo
	for _, group := range groups {
		groupInfos = append(groupInfos, model.GroupInfo{
			ID:   strconv.FormatInt(group.GroupID, 10),
			Name: group.GroupName,
		})
	}
	
	return groupInfos, err
}

// GetGroupsByGroupID 获取指定群组的所有管理员关联
func (s *GroupServiceImpl) GetGroupsByGroupID(ctx context.Context, groupID int64) ([]vo.GroupVo, error) {
	var groups []*model.Group
	
	err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("group_id = ? AND status = 0", groupID).
		Order("id DESC").
		Find(&groups).Error
	
	// 转换为 VO 结构
	var groupVos []vo.GroupVo
	for _, group := range groups {
		groupVos = append(groupVos, vo.GroupVo{
			ID:         group.ID,
			AdminID:    group.AdminID,
			GroupID:    group.GroupID,
			GroupName:  group.GroupName,
			Status:     group.Status,
			CreateTime: group.CreateTime,
			UpdateTime: group.UpdateTime,
		})
	}
	
	return groupVos, err
}

// GetGroupByID 根据ID获取群组信息
func (s *GroupServiceImpl) GetGroupByID(ctx context.Context, id int) (*vo.GroupVo, error) {
	var group model.Group
	
	if err := s.db.WithContext(ctx).Where("id = ? AND status = 0", id).First(&group).Error; err != nil {
		return nil, err
	}
	
	// 转换为 VO 结构
	groupVo := &vo.GroupVo{
		ID:         group.ID,
		AdminID:    group.AdminID,
		GroupID:    group.GroupID,
		GroupName:  group.GroupName,
		Status:     group.Status,
		CreateTime: group.CreateTime,
		UpdateTime: group.UpdateTime,
	}
	
	return groupVo, nil
}