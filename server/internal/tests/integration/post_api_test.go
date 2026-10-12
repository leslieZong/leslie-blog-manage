package integration

import (
	"leslie-blog-server/internal/modules/post/dto"
	"leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/pkg/ulid"
	"net/http"
	"testing"
)

func TestPostAPI_Create_Unauthorized(
	t *testing.T,
) {
	app := NewTestApplication(t)

	body := map[string]any{
		"title":   "Go Testing",
		"slug":    "go-testing",
		"content": "test",
	}

	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts",
		body,
		"",
	)

	AssertStatus(
		t,
		resp,
		http.StatusUnauthorized,
	)
}

func TestPostAPI_Create_AdminSuccess(
	t *testing.T,
) {
	// Arrange
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

	tag2 := CreateTagFixture(
		t,
		app.Server.DB(),
		ulid.New(),
		"GORM-"+ulid.New(),
		"gorm-"+ulid.New(),
	)

	token := GenerateTestToken(
		t,
		app,
		admin,
	)

	body := dto.CreatePostRequest{
		Title:      "Go Testing",
		Slug:       "go-testing",
		CategoryID: category.ID,
		Content:    "# Go Testing",
		TagIDs:     []string{tag1.ID, tag2.ID},
		Featured:   false,
	}

	// Act
	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts",
		body,
		token,
	)

	// Assert
	AssertStatus(
		t,
		resp,
		http.StatusOK,
	)

	response := DecodeAPIResponse(
		t,
		resp,
	)

	if response.Code != 0 {
		t.Fatalf(
			"expected success code 0, got %d",
			response.Code,
		)
	}

	var post model.Post

	err := app.Server.DB().
		Where("slug = ?", "go-testing").
		First(&post).
		Error

	if err != nil {
		t.Fatalf(
			"post should exist in database: %v",
			err,
		)
	}

	// Author 必须来自 JWT 对应的 Actor。
	if post.AuthorID != admin.ID {
		t.Fatalf(
			"expected author %s, got %s",
			admin.ID,
			post.AuthorID,
		)
	}
	// Category 必须正确。
	if post.CategoryID != category.ID {
		t.Fatalf(
			"expected category %s, got %s",
			category.ID,
			post.CategoryID,
		)
	}
	// 新文章默认必须是 draft。
	if post.Status != string(model.PostStatusDraft) {
		t.Fatalf(
			"expected draft status, got %s",
			post.Status,
		)
	}
	// Featured 应该与请求一致。
	if post.Featured != body.Featured {
		t.Fatalf(
			"expected featured %v, got %v",
			body.Featured,
			post.Featured,
		)
	}

	var relations []model.PostTag

	err = app.Server.DB().
		Where(
			"post_id = ?",
			post.ID,
		).
		Find(&relations).
		Error

	if err != nil {
		t.Fatalf(
			"find post tags failed: %v",
			err,
		)
	}
	if len(relations) != 2 {
		t.Fatalf(
			"expected 2 post tags, got %d",
			len(relations),
		)
	}
	for _, relation := range relations {
		if relation.TagID != tag1.ID && relation.TagID != tag2.ID {
			t.Fatalf(
				"expected tag %s or %s, got %s",
				tag1.ID,
				tag2.ID,
				relation.TagID,
			)
		}
	}

}
