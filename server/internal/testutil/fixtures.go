package testutil

import (
	"fmt"
	categorymodel "leslie-blog-server/internal/modules/category/model"
	postmodel "leslie-blog-server/internal/modules/post/model"
	tagmodel "leslie-blog-server/internal/modules/tag/model"
	usermodel "leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/pkg/ulid"
	"testing"

	"gorm.io/gorm"
)

type PostFixture struct {
	User     *usermodel.User
	Category *categorymodel.Category
	Tags     []*tagmodel.Tag
	Post     *postmodel.Post
}

func CreateUserFixture(
	t *testing.T,
	db *gorm.DB,
) *usermodel.User {

	t.Helper()

	user := &usermodel.User{
		ID:           ulid.New(),
		Username:     "test-user-" + ulid.New(),
		Email:        nil,
		PasswordHash: "$2a$10$8w2.j4Q4hD8v5/zQyMfvAuw5x/yw5x/yw5x/yw5x/yw5x/yw5x/yw5x",
		Status:       1,
		DisplayName:  "Test User",
		AvatarURL:    "",
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf(
			"create test user failed: %v",
			err,
		)
	}

	return user
}

func CreateCategoryFixture(
	t *testing.T,
	db *gorm.DB,
) *categorymodel.Category {

	t.Helper()

	category := &categorymodel.Category{
		ID:          ulid.New(),
		Name:        "Test Category-" + ulid.New(),
		Slug:        "test-category-" + ulid.New(),
		Description: nil,
		Status:      1,
		Sort:        0,
	}

	if err := db.Create(category).Error; err != nil {
		t.Fatalf(
			"create test category failed: %v",
			err,
		)
	}

	return category
}

func CreateTagFixture(
	t *testing.T,
	db *gorm.DB,
	id string,
	name string,
	slug string,
) *tagmodel.Tag {

	t.Helper()

	tag := &tagmodel.Tag{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: nil,
		Status:      1,
	}

	if err := db.Create(tag).Error; err != nil {
		t.Fatalf(
			"create test tag failed: %v",
			err,
		)
	}

	return tag
}

func CreatePostFixture(
	t *testing.T,
	db *gorm.DB,
	user *usermodel.User,
	category *categorymodel.Category,
) *postmodel.Post {

	t.Helper()

	post := &postmodel.Post{
		ID:         ulid.New(),
		Title:      "Test Post-" + ulid.New(),
		Slug:       "test-post-" + ulid.New(),
		Content:    "test content-" + ulid.New(),
		Summary:    nil,
		Cover:      nil,
		Status:     string(postmodel.PostStatusDraft),
		AuthorID:   user.ID,
		CategoryID: category.ID,
	}

	if err := db.Create(post).Error; err != nil {
		t.Fatalf(
			"create test post failed: %v",
			err,
		)
	}

	return post
}

func CreateFullPostFixture(
	t *testing.T,
	db *gorm.DB,
) *PostFixture {

	t.Helper()

	user := CreateUserFixture(t, db)

	category := CreateCategoryFixture(t, db)

	tag1 := CreateTagFixture(
		t,
		db,
		ulid.New(),
		"Go-"+ulid.New(),
		"go-"+ulid.New(),
	)

	tag2 := CreateTagFixture(
		t,
		db,
		ulid.New(),
		"GORM-"+ulid.New(),
		"gorm-"+ulid.New(),
	)
	post := CreatePostFixture(
		t,
		db,
		user,
		category,
	)
	tags := []*tagmodel.Tag{
		tag1,
		tag2,
	}

	CreatePostTagFixtures(
		t,
		db,
		post,
		tags,
	)

	return &PostFixture{
		User:     user,
		Category: category,
		Tags:     []*tagmodel.Tag{tag1, tag2},
		Post:     post,
	}
}
func CreatePendingPostFixture(
	t *testing.T,
	db *gorm.DB,
) *PostFixture {

	t.Helper()

	user := CreateUserFixture(t, db)

	category := CreateCategoryFixture(t, db)

	tag1 := CreateTagFixture(
		t,
		db,
		ulid.New(),
		"Go-"+ulid.New(),
		"go-"+ulid.New(),
	)

	tag2 := CreateTagFixture(
		t,
		db,
		ulid.New(),
		"GORM-"+ulid.New(),
		"gorm-"+ulid.New(),
	)
	post := &postmodel.Post{
		ID:         ulid.New(),
		Title:      "Test Post-" + ulid.New(),
		Slug:       "test-post-" + ulid.New(),
		Content:    "test content-" + ulid.New(),
		Summary:    nil,
		Cover:      nil,
		Status:     string(postmodel.PostStatusDraft),
		AuthorID:   user.ID,
		CategoryID: category.ID,
	}

	return &PostFixture{
		User:     user,
		Category: category,
		Tags:     []*tagmodel.Tag{tag1, tag2},
		Post:     post,
	}
}

func CreatePostTagFixtures(
	t *testing.T,
	db *gorm.DB,
	post *postmodel.Post,
	tags []*tagmodel.Tag,
) {

	t.Helper()

	for _, tag := range tags {

		relation := &postmodel.PostTag{
			PostID: post.ID,
			TagID:  tag.ID,
		}

		if err := db.Create(relation).Error; err != nil {
			t.Fatalf(
				"create post tag failed: %v",
				err,
			)
		}
	}
}

func CreatePostFixtures(
	t *testing.T,
	db *gorm.DB,
	user *usermodel.User,
	category *categorymodel.Category,
	count int,
) []*postmodel.Post {

	t.Helper()

	posts := make(
		[]*postmodel.Post,
		0,
		count,
	)

	for i := 0; i < count; i++ {

		post := &postmodel.Post{
			ID:         ulid.New(),
			Title:      fmt.Sprintf("Test Post %d", i),
			Slug:       fmt.Sprintf("test-post-%d", i),
			Content:    "test content",
			Status:     string(postmodel.PostStatusDraft),
			AuthorID:   user.ID,
			CategoryID: category.ID,
		}

		if err := db.Create(post).Error; err != nil {
			t.Fatalf(
				"create test post failed: %v",
				err,
			)
		}

		posts = append(posts, post)
	}

	return posts
}
