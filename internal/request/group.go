package request

// CreateGroupRequest 创建群组关联请求
type CreateGroupRequest struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	GroupName string `json:"group_name" binding:"required"`
}

// UpdateGroupRequest 更新群组信息请求
type UpdateGroupRequest struct {
	ID        int    `json:"id" binding:"required"`
	GroupID   int64  `json:"group_id" binding:"required"`
	GroupName string `json:"group_name" binding:"required"`
}

// SearchGroupRequest 群组列表查询请求
// 支持按 group_id 查询和分页
type SearchGroupRequest struct {
	PageRequest
	GroupID *int64 `json:"group_id,omitempty"` // 可选的群组ID过滤
	AdminID int    `json:"admin_id,omitempty"` // 可选的管理员ID过滤
	Status  *int   `json:"status,omitempty"`   // 可选的状态过滤
}

// DeleteGroupRequest 删除群组关联请求
// 支持按ID删除或按管理员ID+群组ID删除
type DeleteGroupRequest struct {
	ID      *int   `json:"id,omitempty"`       // 按记录ID删除
	AdminID *int   `json:"admin_id,omitempty"` // 按管理员ID删除
	GroupID *int64 `json:"group_id,omitempty"` // 按群组ID删除
}
