package repository

import (
	"context"

	"leslie-blog-server/internal/modules/tag/model"
)

// TagRepository 定义 Tag 数据访问能力。
//
// Service 不应该直接操作 GORM。
//
// 所以：
//
// Service
//
//	↓
//
// TagRepository
//
//	↓
//
// GORM
//
//	↓
//
// # MySQL
//
// Repository 只负责数据库相关操作。
type TagRepository interface {

	// FindByID 根据 ID 查询 Tag。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.Tag, error)

	// FindByName 根据名称查询 Tag。
	FindByName(
		ctx context.Context,
		name string,
	) (*model.Tag, error)

	// FindBySlug 根据 slug 查询 Tag。
	FindBySlug(
		ctx context.Context,
		slug string,
	) (*model.Tag, error)

	// FindAll 查询所有 Tag。
	FindAll(
		ctx context.Context,
	) ([]*model.Tag, error)

	// Create 创建 Tag。
	Create(
		ctx context.Context,
		tag *model.Tag,
	) error

	// Update 更新 Tag。
	Update(
		ctx context.Context,
		tag *model.Tag,
	) error

	// Delete 删除 Tag。
	//
	// 实际使用 GORM 的软删除。
	Delete(
		ctx context.Context,
		id string,
	) error

	// FindByIDs 根据多个 ID 查询 Tag。
	//
	// 这个方法非常重要。
	//
	// 创建 Post 时前端可能提交：
	//
	// tags:
	// [
	//     "01K...",
	//     "01K...",
	//     "01K..."
	// ]
	//
	// Service 需要验证：
	//
	// 这些 Tag 是否真的存在？
	FindByIDs(
		ctx context.Context,
		ids []string,
	) ([]*model.Tag, error)

	FindPage(
		ctx context.Context,
		query TagListQuery,
	) ([]*model.Tag, int64, error)

	CountPosts(
		ctx context.Context,
		tagID string,
	) (int64, error)
}
