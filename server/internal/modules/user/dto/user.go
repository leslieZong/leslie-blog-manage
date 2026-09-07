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

// CreateUserRequest 创建用户请求。
type CreateUserRequest struct {

	// 用户名：
	//
	// 不能为空
	// 最少 3 个字符
	// 最多 50 个字符
	Username string `json:"username" validate:"required,min=3,max=50"`

	// 密码：
	//
	// 不能为空
	// 最少 6 个字符
	// bcrypt 实际最大支持 72 bytes。
	//
	// 注意这里暂时使用字符长度约束，
	// 后面可以进一步处理 UTF-8 byte 长度问题。
	Password string `json:"password" validate:"required,min=6,max=72"`

	// Email 可以为空。
	//
	// 如果传了，就必须符合 email 格式。
	Email *string `json:"email" validate:"omitempty,email"`

	// 显示名称。
	DisplayName string `json:"displayName" validate:"required,max=100"`

	// 头像地址可以为空。
	AvatarURL string `json:"avatarUrl" validate:"omitempty,max=500"`
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
	// Email 可以为空。
	//
	// 如果传入，则必须符合 Email 格式。
	Email *string `json:"email" validate:"omitempty,email"`

	// 显示名称不能为空，并且最多 100 个字符。
	DisplayName string `json:"displayName" validate:"required,max=100"`

	// 头像 URL 可以为空，但长度不能超过 500。
	AvatarURL string `json:"avatarUrl" validate:"omitempty,max=500"`

	// 用户状态。
	//
	// 目前约定：
	//
	// 1 = 正常
	// 0 = 禁用
	Status int8 `json:"status" validate:"oneof=0 1"`
}
