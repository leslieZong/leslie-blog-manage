package repository

import (
	"context"

	"leslie-blog-server/internal/modules/post/model"

	"gorm.io/gorm"
)

// postTagRepository 是 PostTagRepository 的 GORM 实现。
type postTagRepository struct {
	db *gorm.DB
}

// NewPostTagRepository 创建 Repository。
func NewPostTagRepository(
	db *gorm.DB,
) PostTagRepository {

	return &postTagRepository{
		db: db,
	}
}

// ReplaceTags 替换文章的全部 Tag。
func (r *postTagRepository) ReplaceTags(
	ctx context.Context,
	postID string,
	tagIDs []string,
) error {

	// -----------------------------------------------------
	// 第一步：删除旧关系
	// -----------------------------------------------------

	if err := r.db.
		WithContext(ctx).
		Where("post_id = ?", postID).
		Delete(&model.PostTag{}).
		Error; err != nil {

		return err
	}

	// -----------------------------------------------------
	// 第二步：如果没有 Tag，就到这里结束。
	// -----------------------------------------------------

	if len(tagIDs) == 0 {
		return nil
	}

	// -----------------------------------------------------
	// 第三步：构造新的关联记录。
	// -----------------------------------------------------

	relations := make(
		[]model.PostTag,
		0,
		len(tagIDs),
	)

	for _, tagID := range tagIDs {

		relations = append(
			relations,
			model.PostTag{
				PostID: postID,
				TagID:  tagID,
			},
		)
	}

	// -----------------------------------------------------
	// 第四步：批量插入。
	// -----------------------------------------------------

	return r.db.
		WithContext(ctx).
		Create(&relations).
		Error
}

// DeleteByPostID 删除文章全部 Tag。
func (r *postTagRepository) DeleteByPostID(
	ctx context.Context,
	postID string,
) error {

	return r.db.
		WithContext(ctx).
		Where("post_id = ?", postID).
		Delete(&model.PostTag{}).
		Error
}
