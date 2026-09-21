package repository

import (
	"context"
	"strings"

	"leslie-blog-server/internal/modules/tag/model"

	"gorm.io/gorm"
)

// tagRepository 是 TagRepository 的 GORM 实现。
//
// Repository 接口负责定义“能做什么”。
//
// 这个结构体负责定义“怎么做”。
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建 Tag Repository。
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

// FindByID 根据 ID 查询 Tag。
func (r *tagRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Tag, error) {

	var tag model.Tag

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&tag).
		Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// FindByName 根据名称查询 Tag。
func (r *tagRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.Tag, error) {

	var tag model.Tag

	err := r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&tag).
		Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// FindBySlug 根据 slug 查询 Tag。
func (r *tagRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.Tag, error) {

	var tag model.Tag

	err := r.db.
		WithContext(ctx).
		Where("slug = ?", slug).
		First(&tag).
		Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// FindAll 查询所有未删除 Tag。
func (r *tagRepository) FindAll(
	ctx context.Context,
) ([]*model.Tag, error) {

	var tags []*model.Tag

	err := r.db.
		WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&tags).
		Error

	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Create 创建 Tag。
func (r *tagRepository) Create(
	ctx context.Context,
	tag *model.Tag,
) error {

	return r.db.
		WithContext(ctx).
		Create(tag).
		Error
}

// Update 更新 Tag。
func (r *tagRepository) Update(
	ctx context.Context,
	tag *model.Tag,
) error {

	return r.db.
		WithContext(ctx).
		Save(tag).
		Error
}

// Delete 删除 Tag。
//
// GORM 如果 Model 包含 DeletedAt，
// Delete 默认执行软删除。
func (r *tagRepository) Delete(
	ctx context.Context,
	id string,
) error {

	return r.db.
		WithContext(ctx).
		Delete(&model.Tag{}, "id = ?", id).
		Error
}

// FindByIDs 根据多个 ID 查询 Tag。
func (r *tagRepository) FindByIDs(
	ctx context.Context,
	ids []string,
) ([]*model.Tag, error) {

	var tags []*model.Tag

	err := r.db.
		WithContext(ctx).
		Where("id IN ?", ids).
		Where("deleted_at IS NULL").
		Where("status = ?", 1).
		Find(&tags).
		Error

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) FindPage(
	ctx context.Context,
	query TagListQuery,
) ([]*model.Tag, int64, error) {

	var (
		tags  []*model.Tag
		total int64
	)

	// --------------------------------------------------
	// 1. 创建基础查询
	// --------------------------------------------------
	//
	// deleted_at IS NULL：
	// 只查询没有被软删除的 Tag。
	//
	// GORM 如果模型使用 *time.Time DeletedAt，
	// 默认也会自动处理软删除。
	//
	// 这里显式写出来，是为了让初学者清楚业务条件。
	//
	db := r.db.
		WithContext(ctx).
		Model(&model.Tag{}).
		Where("deleted_at IS NULL")

	// --------------------------------------------------
	// 2. 关键字搜索
	// --------------------------------------------------

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

	// --------------------------------------------------
	// 3. 状态筛选
	// --------------------------------------------------

	if query.Status != nil {

		db = db.Where(
			"status = ?",
			*query.Status,
		)
	}

	// --------------------------------------------------
	// 4. 查询总数量
	// --------------------------------------------------

	if err := db.
		Count(&total).
		Error; err != nil {

		return nil, 0, err
	}

	// --------------------------------------------------
	// 5. 查询当前页数据
	// --------------------------------------------------

	if err := db.
		Order("created_at DESC").
		Limit(query.PageSize).
		Offset(query.Offset()).
		Find(&tags).
		Error; err != nil {

		return nil, 0, err
	}

	return tags, total, nil
}

func (r *tagRepository) CountPosts(
	ctx context.Context,
	tagID string,
) (int64, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Table("post_tags AS pt").
		Joins(
			"JOIN posts AS p ON p.id = pt.post_id",
		).
		Where(
			"pt.tag_id = ?",
			tagID,
		).
		Where(
			"p.deleted_at IS NULL",
		).
		Count(&count).
		Error

	return count, err
}
