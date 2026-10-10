package repository

import (
	"context"
	"errors"
	"leslie-blog-server/internal/modules/post/model"
	pkgDatabase "leslie-blog-server/internal/pkg/database"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/testutil"
	"testing"

	"gorm.io/gorm"
)

func TestPostRepository_FindByID(t *testing.T) {

	db := testutil.NewTestDB(t)

	repo := NewPostRepository(db)

	// Arrange
	post := testutil.CreateFullPostFixture(t, db)

	// Act
	actual, err := repo.FindByID(
		context.Background(),
		post.Post.ID,
	)

	// Assert
	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if actual == nil {
		t.Fatal("expected post, got nil")
	}

	if actual.ID != post.Post.ID {
		t.Fatalf(
			"expected id %s, got %s",
			post.Post.ID,
			actual.ID,
		)
	}
	// Check category
	if actual.Category == nil {
		t.Fatal("expected category to be preloaded")
	}

	if actual.Category.ID != post.Category.ID {
		t.Fatalf(
			"expected category %s, got %s",
			post.Category.ID,
			actual.Category.ID,
		)
	}
	// Check tags
	if len(actual.Tags) != 2 {
		t.Fatalf(
			"expected 2 tags, got %d",
			len(actual.Tags),
		)
	}
	foundGo := false
	foundGORM := false

	for _, tag := range actual.Tags {

		switch tag.ID {
		case post.Tags[0].ID:
			foundGo = true

		case post.Tags[1].ID:
			foundGORM = true
		}
	}

	if !foundGo || !foundGORM {
		t.Fatal(
			"expected fixture tags to be preloaded",
		)
	}
	testutil.CleanupPostFixture(t, db, post)
}

func TestPostRepository_FindBySlug(t *testing.T) {

	db := testutil.NewTestDB(t)

	repo := NewPostRepository(db)

	post := testutil.CreateFullPostFixture(t, db)

	actual, err := repo.FindBySlug(
		context.Background(),
		post.Post.Slug,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if actual == nil {
		t.Fatal("expected post, got nil")
	}

	if actual.Slug != post.Post.Slug {
		t.Fatalf(
			"expected slug %s, got %s",
			post.Post.Slug,
			actual.Slug,
		)
	}
	testutil.CleanupPostFixture(t, db, post)
}

func TestPostRepository_FindByID_NotFound(
	t *testing.T,
) {

	db := testutil.NewTestDB(t)

	repo := NewPostRepository(db)

	_, err := repo.FindByID(
		context.Background(),
		"01NOTEXIST0000000000000000",
	)

	if !errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		t.Fatalf(
			"expected gorm.ErrRecordNotFound, got %v",
			err,
		)
	}
}

func TestPostRepository_Delete(t *testing.T) {

	db := testutil.NewTestDB(t)

	repo := NewPostRepository(db)

	post := testutil.CreateFullPostFixture(t, db)

	err := repo.Delete(
		context.Background(),
		post.Post.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	_, err = repo.FindByID(
		context.Background(),
		post.Post.ID,
	)

	if !errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		t.Fatalf(
			"expected record not found after delete, got %v",
			err,
		)
	}
}

func TestPostRepository_FindPage(t *testing.T) {
	db := testutil.NewTestDB(t)

	repo := NewPostRepository(db)

	var posts []*testutil.PostFixture
	var testTotal int = 10
	for i := 0; i < testTotal; i++ {
		post := testutil.CreateFullPostFixture(t, db)
		posts = append(posts, post)
	}
	var params pagination.Params = pagination.Params{
		Page:     1,
		PageSize: 2,
	}

	actual, total, err := repo.FindPage(
		context.Background(),
		PostListQuery{
			Params: params,
			// CategoryID: posts[0].Category.ID,
		},
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if actual == nil {
		t.Fatal("expected post, got nil")
	}

	if total != int64(testTotal) {
		t.Fatalf(
			"expected total %d, got %d",
			testTotal,
			int(total),
		)
	}
	if len(actual) != params.PageSize {
		t.Fatalf(
			"expected %d posts, got %d",
			params.PageSize,
			len(actual),
		)
	}
	if actual[0].ID != posts[0].Post.ID {
		t.Fatalf(
			"expected post %s, got %s",
			posts[0].Post.ID,
			actual[0].ID,
		)
	}
	for _, post := range posts {
		testutil.CleanupPostFixture(t, db, post)
	}
}

func TestTransaction_Rollback(t *testing.T) {

	db := testutil.NewTestDB(t)

	txManager := pkgDatabase.NewTransactionManager(db)

	post := testutil.CreatePendingPostFixture(t, db)

	err := txManager.WithTransaction(
		context.Background(),
		func(tx *gorm.DB) error {

			if err := tx.Create(post).Error; err != nil {
				return err
			}

			// 故意制造错误。
			return errors.New(
				"force rollback",
			)
		},
	)

	if err == nil {
		t.Fatal(
			"expected transaction error",
		)
	}

	var count int64

	err = db.
		Model(&model.Post{}).
		Where("id = ?", post.Post.ID).
		Count(&count).
		Error

	if err != nil {
		t.Fatalf(
			"count failed: %v",
			err,
		)
	}

	if count != 0 {
		t.Fatalf(
			"expected rollback, but post exists",
		)
	}
}
