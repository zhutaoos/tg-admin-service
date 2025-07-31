package response

import "time"

// UserResponse 用户响应对象
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Status    int32     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	// 注意：不包含密码等敏感信息
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// UserProfileResponse 用户详情响应
type UserProfileResponse struct {
	UserResponse
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int64      `json:"login_count"`
}
