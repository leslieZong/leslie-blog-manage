package dto

import (
	"leslie-blog-server/internal/modules/user/model"
	"time"
)

// UserResponse 是返回给前端的用户信息。
//
// 注意：
//
// 这个结构不是数据库 Model。
//
// 它专门用于：
//
// Go Backend
//
//	↓
//
// HTTP Response
//
//	↓
//
// # Frontend
//
// 最重要的作用之一就是：
//
// “只返回允许暴露的数据”。
type UserResponse struct {

	// 用户 ID。
	ID string `json:"id"`

	// 登录用户名。
	Username string `json:"username"`

	// 用户邮箱。
	//
	// 使用 *string 是为了保留 NULL。
	Email *string `json:"email"`

	// 用户显示名称。
	DisplayName string `json:"displayName"`

	// 用户头像地址。
	AvatarURL string `json:"avatarUrl"`

	// 用户状态。
	//
	// 1 = 正常
	// 0 = 禁用
	Status int8 `json:"status"`

	// 最后登录时间。
	LastLoginAt *time.Time `json:"lastLoginAt"`

	// 创建时间。
	CreatedAt time.Time `json:"createdAt"`

	// 更新时间。
	UpdatedAt time.Time `json:"updatedAt"`
}

// FromUser 将 Model 转换成 Response DTO。
//
// 这个函数非常重要。
//
// 它明确告诉我们：
//
// User Model
//
//	↓
//
// # UserResponse
//
// 并且这里故意没有：
//
// PasswordHash
// DeletedAt
//
// 因此这两个字段不会返回给前端。
func FromUser(user *model.User) *UserResponse {

	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
