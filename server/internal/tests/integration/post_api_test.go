package integration

import (
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/post/dto"
	"leslie-blog-server/internal/modules/post/model"
	postmodel "leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/pkg/ulid"
	"net/http"
	"testing"
)

func TestPostAPI_Create_Unauthorized(
	t *testing.T,
) {
	app := NewTestApplication(t)

	body := map[string]any{
		"title":   "Unauthorized Post",
		"slug":    "unauthorized-post",
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

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(appErrors.ErrUnauthorized) {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrUnauthorized,
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&postmodel.Post{}).
		Where(
			"slug = ?",
			"unauthorized-post",
		).
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
			"expected no post to be created, got %d",
			count,
		)
	}
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

func TestPostAPI_Create_InvalidToken(t *testing.T) {
	app := NewTestApplication(t)

	body := map[string]any{
		"title":   "Unauthorized Post",
		"slug":    "unauthorized-post",
		"content": "test",
	}

	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts",
		body,
		"invalid-token",
	)

	AssertStatus(
		t,
		resp,
		http.StatusUnauthorized,
	)

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(appErrors.ErrUnauthorized) {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrUnauthorized,
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&postmodel.Post{}).
		Where(
			"slug = ?",
			"unauthorized-post",
		).
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
			"expected no post to be created, got %d",
			count,
		)
	}
}
func TestPostAPI_Create_ErrForbidden(t *testing.T) {
	// Arrange
	app := NewTestApplication(t)

	viewer := CreateViewerFixture(
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
		viewer,
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
		http.StatusForbidden,
	)

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(appErrors.ErrForbidden) {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrForbidden,
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&postmodel.Post{}).
		Where(
			"slug = ?",
			"unauthorized-post",
		).
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
			"expected no post to be created, got %d",
			count,
		)
	}
}

func TestPostAPI_Create_InvalidParams(
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

	slug := "invalid-post-" + ulid.New()
	body := dto.CreatePostRequest{
		Title:      "",
		Slug:       slug,
		CategoryID: category.ID,
		Content:    "",
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
		http.StatusBadRequest,
	)

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(appErrors.ErrInvalidParams) {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrInvalidParams,
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&postmodel.Post{}).
		Where(
			"slug = ?",
			slug,
		).
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
			"expected no post to be created, got %d",
			count,
		)
	}

}
func TestPostAPI_Create_SameSlug(
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

	slug := "go-testing"
	body := dto.CreatePostRequest{
		Title:      "Same Post",
		Slug:       slug,
		CategoryID: category.ID,
		Content:    "# Same Post",
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
		http.StatusConflict,
	)

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(appErrors.ErrPostSlugExists) {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrPostSlugExists,
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&postmodel.Post{}).
		Where(
			"slug = ?",
			slug,
		).
		Count(&count).
		Error

	if err != nil {
		t.Fatalf(
			"count posts failed: %v",
			err,
		)
	}

	if count != 1 {
		t.Fatalf(
			"expected 1 post to be created, got %d",
			count,
		)
	}

}

func TestPostAPI_Update_NotFound(
	t *testing.T,
) {
	// Arrange
	app := NewTestApplication(t)

	admin := CreateAdminFixture(
		t,
		app,
	)

	token := GenerateTestToken(
		t,
		app,
		admin,
	)

	// Act
	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodGet,
		"/api/admin/v1/posts/not-exists",
		nil,
		token,
	)

	// Assert
	AssertStatus(
		t,
		resp,
		http.StatusNotFound,
	)

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != int(
		appErrors.ErrPostNotFound,
	) {
		t.Fatalf(
			"expected post not found, got %d",
			apiResponse.Code,
		)
	}

}

func TestPostAPI_Update_NonOwnerForbidden(
	t *testing.T,
) {
	// Arrange
	app := NewTestApplication(t)

	user := CreateEditorFixture(
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
		user,
	)

	slug := "go-testing-" + ulid.New()
	body := dto.CreatePostRequest{
		Title:      "Go Testing-" + ulid.New(),
		Slug:       slug,
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

	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)

	if apiResponse.Code != 0 {
		t.Fatalf(
			"expected success code 0, got %d",
			apiResponse.Code,
		)
	}

	var post model.Post

	err := app.Server.DB().
		Where("slug = ?", slug).
		First(&post).
		Error

	if err != nil {
		t.Fatalf(
			"query post failed: %v",
			err,
		)
	}

	otherUser := CreateEditorFixture(
		t,
		app,
	)

	token = GenerateTestToken(
		t,
		app,
		otherUser,
	)
	request := map[string]any{
		"title": "Hacked Title",
	}

	resp = DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPut,
		"/api/admin/v1/posts/"+post.ID,
		request,
		token,
	)

	AssertStatus(
		t,
		resp,
		http.StatusForbidden,
	)

	var updated model.Post

	err = app.Server.DB().
		Where("id = ?", post.ID).
		First(&updated).
		Error

	if err != nil {
		t.Fatalf(
			"query post failed: %v",
			err,
		)
	}
	originalTitle := post.Title
	if updated.Title != originalTitle {
		t.Fatalf(
			"post should not be modified, expected %q, got %q",
			originalTitle,
			updated.Title,
		)
	}
}
