package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"leslie-blog-server/internal/modules/category/model"
	postModel "leslie-blog-server/internal/modules/post/model"
)

// gormCategoryRepository
//
// CategoryRepository 的 GORM 实现。
type gormCategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository
//
// 创建 Category Repository。
func NewCategoryRepository(
	db *gorm.DB,
) CategoryRepository {
	return &gormCategoryRepository{
		db: db,
	}
}

// FindByID
//
// 根据 ID 查询分类。
func (r *gormCategoryRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Category, error) {

	var category model.Category

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByName
//
// 根据名称查询分类。
func (r *gormCategoryRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.Category, error) {

	var category model.Category

	err := r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindBySlug
//
// 根据 slug 查询分类。
func (r *gormCategoryRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.Category, error) {

	var category model.Category

	err := r.db.
		WithContext(ctx).
		Where("slug = ?", slug).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindAll
//
// 查询所有未删除分类。
//
// 排序规则：
// 1. sort ASC
// 2. created_at ASC
func (r *gormCategoryRepository) FindAll(
	ctx context.Context,
) ([]*model.Category, error) {

	var categories []*model.Category

	err := r.db.
		WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("sort ASC").
		Order("created_at ASC").
		Find(&categories).
		Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create
//
// 创建分类。
func (r *gormCategoryRepository) Create(
	ctx context.Context,
	category *model.Category,
) error {

	return r.db.
		WithContext(ctx).
		Create(category).
		Error
}

// Update
//
// 更新分类。
func (r *gormCategoryRepository) Update(
	ctx context.Context,
	category *model.Category,
) error {

	return r.db.
		WithContext(ctx).
		Save(category).
		Error
}

// Delete
//
// 删除分类。
//
// 这里采用软删除。
// 不直接 DELETE 数据。
func (r *gormCategoryRepository) Delete(
	ctx context.Context,
	id string,
) error {

	return r.db.
		WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Update(
			"deleted_at",
			gorm.Expr("NOW()"),
		).
		Error
}

// CountPosts
//
// 查询某个分类下面有多少篇文章。
func (r *gormCategoryRepository) CountPosts(
	ctx context.Context,
	categoryID string,
) (int64, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Model(&postModel.Post{}).
		Where("category_id = ?", categoryID).
		Where("deleted_at IS NULL").
		Count(&count).
		Error

	return count, err
}

func (r *gormCategoryRepository) FindPage(
	ctx context.Context,
	query CategoryListQuery,
) ([]*model.Category, int64, error) {

	var (
		categories []*model.Category
		total      int64
	)

	db := r.db.
		WithContext(ctx).
		Model(&model.Category{}).
		Where("deleted_at IS NULL")

	// 关键字搜索。
	if query.Keyword != "" {

		keyword := "%" +
			strings.TrimSpace(query.Keyword) +
			"%"

		db = db.Where(
			"name LIKE ? OR slug LIKE ?",
			keyword,
			keyword,
		)
	}

	// 状态筛选。
	if query.Status != nil {

		db = db.Where(
			"status = ?",
			*query.Status,
		)
	}

	// Count。
	if err := db.
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	// List。
	if err := db.
		Order("sort ASC").
		Order("created_at DESC").
		Limit(query.PageSize).
		Offset(query.Offset()).
		Find(&categories).
		Error; err != nil {

		return nil, 0, err
	}

	return categories, total, nil
}
