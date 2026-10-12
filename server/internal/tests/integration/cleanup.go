package integration

import (
	categorymodel "leslie-blog-server/internal/modules/category/model"
	postmodel "leslie-blog-server/internal/modules/post/model"
	tagmodel "leslie-blog-server/internal/modules/tag/model"
	usermodel "leslie-blog-server/internal/modules/user/model"
	"testing"

	"gorm.io/gorm"
)

func CleanupPostFixture(
	t *testing.T,
	db *gorm.DB,
	fixture *PostFixture,
) {

	t.Helper()

	db.
		Unscoped().
		Where("post_id = ?", fixture.Post.ID).
		Delete(&postmodel.PostTag{})

	db.
		Unscoped().
		Where("id = ?", fixture.Post.ID).
		Delete(&postmodel.Post{})

	for _, tag := range fixture.Tags {
		db.
			Unscoped().
			Where("id = ?", tag.ID).
			Delete(&tagmodel.Tag{})
	}

	db.
		Unscoped().
		Where("id = ?", fixture.Category.ID).
		Delete(&categorymodel.Category{})

	db.
		Unscoped().
		Where("id = ?", fixture.User.ID).
		Delete(&usermodel.User{})
}
