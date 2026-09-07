package dto

import (
	"time"

	"leslie-blog-server/internal/modules/user/model"
)

// ============================================================
// UserResponse
// ============================================================
//
// UserResponse 是返回给前端的用户数据。
//
// 注意：
// PasswordHash 永远不能返回给前端。
//
// 数据库 Model：
//
//	User
//	├── ID
//	├── Username
//	├── PasswordHash  ← 敏感数据
//	└── ...
//
// API Response：
//
//	UserResponse
//	├── ID
//	├── Username
//	└── ...
//
// 这就是 DTO 的意义。
// ============================================================

type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       *string    `json:"email"`
	DisplayName string     `json:"displayName"`
	AvatarURL   string     `json:"avatarUrl"`
	Status      int8       `json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// FromUser 将数据库 Model 转换成 API Response。
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

// ============================================================
// CreateUserRequest
// ============================================================
//
// 创建用户时，前端只需要告诉后端：
//
//	username
//	password
//	email
//	displayName
//	avatarUrl
//
// 不需要：
//
//	id
//	passwordHash
//	createdAt
//	updatedAt
//
// 这些字段由后端负责。
// ============================================================

type CreateUserRequest struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Email       *string `json:"email"`
	DisplayName string  `json:"displayName"`
	AvatarURL   string  `json:"avatarUrl"`
}

// ============================================================
// UpdateUserRequest
// ============================================================
//
// 修改用户时，不允许前端修改：
//
//	ID
//	PasswordHash
//	CreatedAt
//
// Password 后面如果需要修改，建议单独设计：
//
//	POST /users/:id/password
//
// 不要把所有修改操作塞到一个接口里。
// ============================================================

type UpdateUserRequest struct {
	Email       *string `json:"email"`
	DisplayName string  `json:"displayName"`
	AvatarURL   string  `json:"avatarUrl"`
	Status      int8    `json:"status"`
}
