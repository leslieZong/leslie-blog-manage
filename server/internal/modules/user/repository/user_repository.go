package repository

import (
	"context"

	"leslie-blog-server/internal/modules/user/model"
)

// UserRepository 定义 User 数据访问能力。
//
// 注意：
//
// 这里是“接口”，
// 它只描述：
//
// “User 数据层能够做什么”。
//
// 它不关心具体使用：
//
// MySQL
// PostgreSQL
// Redis
// Mock
//
// 具体实现放到后面。
type UserRepository interface {

	// FindByID 根据用户 ID 查询用户。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.User, error)

	// FindByUsername 根据用户名查询用户。
	FindByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)

	// Create 创建一个新用户。
	Create(
		ctx context.Context,
		user *model.User,
	) error

	// Update 更新用户信息。
	Update(
		ctx context.Context,
		user *model.User,
	) error
}
