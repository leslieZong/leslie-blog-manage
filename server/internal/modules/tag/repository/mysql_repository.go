package repository

import (
	"context"

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
