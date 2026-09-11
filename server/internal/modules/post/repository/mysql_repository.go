package repository

import (
	"context"
	"leslie-blog-server/internal/modules/post/model"

	"gorm.io/gorm"
)

// gormPostRepository 是 PostRepository 的 GORM 实现。
//
// 注意：
// 类型名称使用小写：
//
// gormPostRepository
//
// 表示它是 package 内部实现细节。
//
// 外部只需要使用：
//
// # PostRepository
//
// 不需要关心具体实现。
type gormPostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建 PostRepository。
//
// 这里使用依赖注入：
//
// 外部传入 db
// ↓
// Repository 保存 db
//
// Repository 自己不负责创建数据库连接。
func NewPostRepository(
	db *gorm.DB,
) PostRepository {
	return &gormPostRepository{
		db: db,
	}
}

// Create 创建文章。
func (r *gormPostRepository) Create(
	ctx context.Context,
	post *model.Post,
) error {

	// WithContext 将当前请求的 context
	// 传递给 GORM。
	//
	// Create 会生成类似：
	//
	// INSERT INTO posts (...)
	// VALUES (...)
	//
	// 的 SQL。
	return r.db.
		WithContext(ctx).
		Create(post).
		Error
}

// FindByID 根据 ID 查询文章。
func (r *gormPostRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	var post model.Post

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&post).
		Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// FindBySlug 根据文章 Slug 查询文章。
func (r *gormPostRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.Post, error) {

	var post model.Post

	err := r.db.
		WithContext(ctx).
		Where("slug = ?", slug).
		First(&post).
		Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// FindAll 查询所有未删除的文章。
//
// 当前版本主要用于 Admin。
// 后面我们会增加：
//
// 分页
// 状态筛选
// 作者筛选
// 分类筛选
// 标签筛选
// 关键词搜索
//
// 到时候会进一步升级为 List / Query 对象。
func (r *gormPostRepository) FindAll(
	ctx context.Context,
) ([]*model.Post, error) {

	var posts []*model.Post

	err := r.db.
		WithContext(ctx).
		Order("created_at DESC").
		Find(&posts).
		Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Update 更新文章。
func (r *gormPostRepository) Update(
	ctx context.Context,
	post *model.Post,
) error {

	return r.db.
		WithContext(ctx).
		Save(post).
		Error
}

// Delete 对文章执行软删除。
func (r *gormPostRepository) Delete(
	ctx context.Context,
	id string,
) error {

	// 先找到文章。
	var post model.Post

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&post).
		Error

	if err != nil {
		return err
	}

	// 更新 deleted_at。
	//
	// 这里不使用真正的 DELETE。
	return r.db.
		WithContext(ctx).
		Model(&post).
		Update("deleted_at", gorm.Expr("NOW()")).
		Error
}

// IncrementViewCount 增加文章阅读量。
func (r *gormPostRepository) IncrementViewCount(
	ctx context.Context,
	id string,
) error {

	return r.db.
		WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn(
			"view_count",
			gorm.Expr("view_count + ?", 1),
		).
		Error
}
