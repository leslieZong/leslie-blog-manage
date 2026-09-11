package repository

import (
	"context"
	"leslie-blog-server/internal/modules/post/model"
)

// PostRepository 定义文章的数据访问能力。
//
// 注意：这里定义的是“能力”，而不是具体实现。
//
// Service 只需要知道：
//
// 你能不能创建文章？
// 你能不能查询文章？
// 你能不能修改文章？
//
// Service 不需要知道底层到底是：
//
// GORM
// MySQL
// PostgreSQL
// Mock
//
// 这就是接口的价值。
type PostRepository interface {

	// Create 创建文章。
	Create(
		ctx context.Context,
		post *model.Post,
	) error

	// FindByID 根据文章 ID 查询文章。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.Post, error)

	// FindBySlug 根据文章 Slug 查询文章。
	FindBySlug(
		ctx context.Context,
		slug string,
	) (*model.Post, error)

	// FindAll 查询文章列表。
	FindAll(
		ctx context.Context,
	) ([]*model.Post, error)

	// Update 更新文章。
	Update(
		ctx context.Context,
		post *model.Post,
	) error

	// Delete 删除文章。
	//
	// 注意：
	// 我们的 Delete 实际上会执行软删除。
	Delete(
		ctx context.Context,
		id string,
	) error

	// IncrementViewCount 增加文章阅读量。
	IncrementViewCount(
		ctx context.Context,
		id string,
	) error
}
