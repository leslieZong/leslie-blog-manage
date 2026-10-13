package integration

import (
	"context"
	"errors"
	"testing"

	postmodel "leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

func TestTransaction_RollbackPost(t *testing.T) {
	app := NewTestApplication(t)
	admin := CreateAdminFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())
	ctx := context.Background()

	post := &postmodel.Post{
		ID:         ulid.New(),
		Title:      "Transaction Test",
		Slug:       "transaction-test-" + ulid.New(),
		Content:    "# Transaction Test",
		Status:     string(postmodel.PostStatusDraft),
		AuthorID:   admin.User.ID,
		CategoryID: category.ID,
	}
	err := app.TransactionManager().WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			// 创建 Post
			if err := tx.Create(post).Error; err != nil {
				return err
			}

			// 故意制造错误
			return errors.New("force rollback")
		},
	)
	if err == nil {
		t.Fatal("expected transaction to fail")
	}
	var count int64

	err = app.DB().
		Model(&postmodel.Post{}).
		Where("id = ?", post.ID).
		Count(&count).
		Error

	if err != nil {
		t.Fatalf(
			"count posts failed: %v",
			err,
		)
	}
	if count != 0 {
		t.Fatalf(
			"expected post to be rolled back, got %d",
			count,
		)
	}
}

func TestPostAPI_Create_RollbackPost(t *testing.T) {
	app := NewTestApplication(t)

	admin := CreateAdminFixture(
		t,
		app,
	)

	category := CreateCategoryFixture(
		t,
		app.Server.DB(),
	)

	tag1 := CreateTagFixture(
		t,
		app.Server.DB(),
		ulid.New(),
		"Go-"+ulid.New(),
		"go-"+ulid.New(),
	)

	// tag2 := CreateTagFixture(
	// 	t,
	// 	app.Server.DB(),
	// 	ulid.New(),
	// 	"GORM-"+ulid.New(),
	// 	"gorm-"+ulid.New(),
	// )

	slug := "go-testing-" + ulid.New()
	post := &postmodel.Post{
		ID:         ulid.New(),
		Title:      "Go Testing",
		Slug:       slug,
		Status:     string(postmodel.PostStatusDraft),
		CategoryID: category.ID,
		Content:    "# Go Testing",
		AuthorID:   admin.User.ID,
		Featured:   false,
	}
	ctx := context.Background()
	err := app.TransactionManager().WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			if err := tx.Create(post).Error; err != nil {
				return err
			}

			// -----------------------------------------
			// 2. 创建第一个 PostTag
			// -----------------------------------------

			postTag1 := &postmodel.PostTag{
				PostID: post.ID,
				TagID:  tag1.ID,
			}

			if err := tx.Create(postTag1).Error; err != nil {
				return err
			}
			// -----------------------------------------
			// 3. 创建非法 PostTag
			// -----------------------------------------

			postTag2 := &postmodel.PostTag{
				PostID: post.ID,
				TagID:  ulid.New(), // 不存在的TagID
			}

			if err := tx.Create(postTag2).Error; err != nil {
				return err
			}

			return nil
		},
	)
	if err == nil {
		t.Fatal("expected transaction to fail")
	}

	var count int64

	err = app.DB().
		Model(&postmodel.Post{}).
		Where("id = ?", post.ID).
		Count(&count).
		Error

	if err != nil {
		t.Fatalf("count posts failed: %v", err)
	}

	if count != 0 {
		t.Fatalf(
			"expected post to be rolled back, got %d",
			count,
		)
	}

	var postTagCount int64

	err = app.DB().
		Model(&postmodel.PostTag{}).
		Where("post_id = ?", post.ID).
		Count(&postTagCount).
		Error

	if err != nil {
		t.Fatalf(
			"count post tags failed: %v",
			err,
		)
	}

	if postTagCount != 0 {
		t.Fatalf(
			"expected post tags rollback, got %d",
			postTagCount,
		)
	}
}
