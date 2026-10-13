package integration

import (
	"errors"
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/post/dto"
	"leslie-blog-server/internal/modules/post/model"
	postmodel "leslie-blog-server/internal/modules/post/model"
	tagmodel "leslie-blog-server/internal/modules/tag/model"
	"leslie-blog-server/internal/pkg/ulid"
	"net/http"
	"testing"

	"gorm.io/gorm"
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

	slug := "go-testing-" + ulid.New()
	body := dto.CreatePostRequest{
		Title:      "Go Testing",
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
		admin.Token,
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
		Where("slug = ?", slug).
		First(&post).
		Error

	if err != nil {
		t.Fatalf(
			"post should exist in database: %v",
			err,
		)
	}

	// Author 必须来自 JWT 对应的 Actor。
	if post.AuthorID != admin.User.ID {
		t.Fatalf(
			"expected author %s, got %s",
			admin.User.ID,
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

	t.Cleanup(func() {
		CleanupPostFixture(t, app.Server.DB(), &PostFixture{
			Post:     &post,
			Tags:     []*tagmodel.Tag{tag1, tag2},
			Category: category,
			User:     admin.User,
		})
	})

	t.Cleanup(func() {
		if err := app.Close(); err != nil {
			t.Logf(
				"close test application failed: %v",
				err,
			)
		}
	})
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
		viewer.Token,
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
		admin.Token,
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
		admin.Token,
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

	// Act
	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodGet,
		"/api/admin/v1/posts/not-exists",
		nil,
		admin.Token,
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
		user.Token,
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

	request := map[string]any{
		"title": "Hacked Title",
	}

	resp = DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPut,
		"/api/admin/v1/posts/"+post.ID,
		request,
		otherUser.Token,
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

type createPostActor string

const (
	createPostActorNone    createPostActor = "none"
	createPostActorAdmin   createPostActor = "admin"
	createPostActorViewer  createPostActor = "viewer"
	createPostActorInvalid createPostActor = "invalid"
)

type createPostTestCase struct {
	name  string
	actor createPostActor

	request dto.CreatePostRequest

	expectedStatus int
	expectedCode   int

	expectCreated bool
}

func TestPostAPI_Create(t *testing.T) {
	app := NewTestApplication(t)

	category := CreateCategoryFixture(
		t,
		app.Server.DB(),
	)

	admin := CreateAdminFixture(
		t,
		app,
	)

	viewer := CreateViewerFixture(
		t,
		app,
	)

	tokens := map[createPostActor]string{
		createPostActorAdmin:   admin.Token,
		createPostActorViewer:  viewer.Token,
		createPostActorInvalid: "invalid-token",
		createPostActorNone:    "",
	}

	tests := []createPostTestCase{
		{
			name:  "admin success",
			actor: createPostActorAdmin,

			request: dto.CreatePostRequest{
				Title:      "Go Testing",
				Slug:       "go-testing-admin-" + ulid.New(),
				CategoryID: category.ID,
				Content:    "# Go Testing",
			},

			expectedStatus: http.StatusOK,
			expectedCode:   0,
			expectCreated:  true,
		},

		{
			name:  "viewer forbidden",
			actor: createPostActorViewer,

			request: dto.CreatePostRequest{
				Title:      "Viewer Post",
				Slug:       "viewer-post-" + ulid.New(),
				CategoryID: category.ID,
				Content:    "# Viewer",
			},

			expectedStatus: http.StatusForbidden,
			expectedCode:   int(appErrors.ErrForbidden),
		},

		{
			name:  "invalid token",
			actor: createPostActorInvalid,

			request: dto.CreatePostRequest{
				Title:      "Invalid Token",
				Slug:       "invalid-token-post-" + ulid.New(),
				CategoryID: category.ID,
				Content:    "# Invalid",
			},

			expectedStatus: http.StatusUnauthorized,
			expectedCode:   int(appErrors.ErrUnauthorized),
		},

		{
			name:  "without token",
			actor: createPostActorNone,

			request: dto.CreatePostRequest{
				Title:      "No Token",
				Slug:       "no-token-post-" + ulid.New(),
				CategoryID: category.ID,
				Content:    "# No Token",
			},

			expectedStatus: http.StatusUnauthorized,
			expectedCode:   int(appErrors.ErrUnauthorized),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tokens[tt.actor]

			response := DoJSONRequest(
				t,
				app.Server.HTTPHandler(),
				http.MethodPost,
				"/api/admin/v1/posts",
				tt.request,
				token,
			)

			AssertStatus(
				t,
				response,
				tt.expectedStatus,
			)

			apiResponse := DecodeAPIResponse(
				t,
				response,
			)

			if apiResponse.Code != tt.expectedCode {
				t.Fatalf(
					"expected code %d, got %d",
					tt.expectedCode,
					apiResponse.Code,
				)
			}

			if tt.expectCreated {
				assertPostCreated(
					t,
					app.Server.DB(),
					tt.request.Slug,
				)
			} else {
				assertPostNotCreated(
					t,
					app.Server.DB(),
					tt.request.Slug,
				)
			}
		})
	}
}

func TestPostAPI_Create_InvalidCategory(t *testing.T) {
	app := NewTestApplication(t)
	admin := CreateAdminFixture(
		t,
		app,
	)
	categoryID := "01INVALIDCATEGORYID"

	request := dto.CreatePostRequest{
		Title:      "Invalid Category",
		Slug:       "invalid-category-post-" + ulid.New(),
		CategoryID: categoryID,
		Content:    "# Invalid Category",
	}

	resp := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts",
		request,
		admin.Token,
	)
	AssertStatus(
		t,
		resp,
		http.StatusBadRequest,
	)
	apiResponse := DecodeAPIResponse(
		t,
		resp,
	)
	if apiResponse.Code != int(appErrors.ErrCategoryNotFound) {
		t.Fatalf(
			"expected code %d, got %d",
			int(appErrors.ErrCategoryNotFound),
			apiResponse.Code,
		)
	}

	var count int64

	err := app.Server.DB().
		Model(&model.Post{}).
		Where("slug = ?", request.Slug).
		Count(&count).
		Error

	if err != nil {
		t.Fatalf(
			"count post failed: %v",
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

func TestPostAPI_Update_OwnerSuccess(t *testing.T) {
	app := NewTestApplication(t)
	editorA := CreateEditorFixture(
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
	post := CreatePostFixture(
		t,
		app.Server.DB(),
		editorA.User,
		category,
	)
	request := dto.UpdatePostRequest{
		Title:      "Updated Go Article",
		Slug:       post.Slug + "-updated",
		CategoryID: category.ID,
		TagIDs: []string{
			tag1.ID,
			tag2.ID,
		},
		Content:  "# Updated Content",
		Featured: true,
	}
	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPut,
		"/api/admin/v1/posts/"+post.ID,
		request,
		editorA.Token,
	)
	AssertStatus(
		t,
		response,
		http.StatusOK,
	)

	var updated postmodel.Post

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&updated).
		Error

	if err != nil {
		t.Fatalf(
			"query updated post failed: %v",
			err,
		)
	}
	if updated.Title != request.Title {
		t.Fatalf(
			"expected title %q, got %q",
			request.Title,
			updated.Title,
		)
	}

	if updated.Content != request.Content {
		t.Fatalf(
			"expected content %q, got %q",
			request.Content,
			updated.Content,
		)
	}
	if updated.CategoryID != category.ID {
		t.Fatalf(
			"expected category %q, got %q",
			category.ID,
			updated.CategoryID,
		)
	}
	if !updated.Featured {
		t.Fatal("expected post to be featured")
	}

	var relations []postmodel.PostTag

	err = app.Server.DB().
		Where("post_id = ?", post.ID).
		Find(&relations).
		Error

	if err != nil {
		t.Fatalf(
			"query post tags failed: %v",
			err,
		)
	}
	if len(relations) != 2 {
		t.Fatalf(
			"expected 2 post tags, got %d",
			len(relations),
		)
	}
	actualTagIDs := map[string]bool{}

	for _, relation := range relations {
		actualTagIDs[relation.TagID] = true
	}

	if !actualTagIDs[tag1.ID] {
		t.Fatalf(
			"expected tag %s",
			tag1.ID,
		)
	}

	if !actualTagIDs[tag2.ID] {
		t.Fatalf(
			"expected tag %s",
			tag2.ID,
		)
	}
}

func TestPostAPI_Update_AdminSuccess(t *testing.T) {
	app := NewTestApplication(t)

	admin := CreateAdminFixture(t, app)
	owner := CreateEditorFixture(t, app)

	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		owner.User,
		category,
	)

	request := dto.UpdatePostRequest{
		Title:      "Admin Updated Post",
		Slug:       post.Slug + "-admin-update",
		CategoryID: category.ID,
		Content:    "# Updated By Admin",
		Featured:   true,
	}

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPut,
		"/api/admin/v1/posts/"+post.ID,
		request,
		admin.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusOK,
	)
	var updated postmodel.Post

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&updated).
		Error

	if err != nil {
		t.Fatalf(
			"query updated post failed: %v",
			err,
		)
	}
	if updated.Title != request.Title {
		t.Fatalf(
			"expected title %q, got %q",
			request.Title,
			updated.Title,
		)
	}

	if updated.Content != request.Content {
		t.Fatalf(
			"expected content %q, got %q",
			request.Content,
			updated.Content,
		)
	}
	if updated.CategoryID != category.ID {
		t.Fatalf(
			"expected category %q, got %q",
			category.ID,
			updated.CategoryID,
		)
	}
	if !updated.Featured {
		t.Fatal("expected post to be featured")
	}
	if updated.Slug != request.Slug {
		t.Fatalf(
			"expected slug %q, got %q",
			request.Slug,
			updated.Slug,
		)
	}

}

func TestPostAPI_Update_ViewerForbidden(t *testing.T) {
	app := NewTestApplication(t)

	owner := CreateEditorFixture(t, app)
	viewer := CreateViewerFixture(t, app)

	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		owner.User,
		category,
	)

	request := dto.UpdatePostRequest{
		Title:   "Should Not Update",
		Content: "forbidden",
	}

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPut,
		"/api/admin/v1/posts/"+post.ID,
		request,
		viewer.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusForbidden,
	)
	var updated model.Post

	err := app.Server.DB().
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

func TestPostAPI_Delete_DeleteSuccess(t *testing.T) {
	app := NewTestApplication(t)
	admin := CreateAdminFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())
	post := CreatePostFixture(
		t,
		app.Server.DB(),
		admin.User,
		category,
	)
	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodDelete,
		"/api/admin/v1/posts/"+post.ID,
		nil,
		admin.Token,
	)
	AssertStatus(
		t,
		response,
		http.StatusOK,
	)

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&model.Post{}).
		Error

	if err == nil {
		t.Fatalf(
			"post should be deleted, but not found",
		)
	}
	if !errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		t.Fatalf(
			"expected record not found after delete, got %v",
			err,
		)
	}
	if len(post.Tags) != 0 {
		t.Fatalf(
			"expected no tags after delete, got %d",
			len(post.Tags),
		)
	}
}

func TestPostAPI_Publish_OwnerSuccess(t *testing.T) {
	app := NewTestApplication(t)

	editor := CreateEditorFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		editor.User,
		category,
	)
	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		editor.Token,
	)
	AssertStatus(
		t,
		response,
		http.StatusOK,
	)

	var actual postmodel.Post

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&actual).
		Error

	if err != nil {
		t.Fatalf(
			"query post failed: %v",
			err,
		)
	}
	if actual.Status != string(postmodel.PostStatusPublished) {
		t.Fatalf(
			"expected status %q, got %q",
			string(postmodel.PostStatusPublished),
			actual.Status,
		)
	}
	if actual.PublishedAt == nil {
		t.Fatal("expected published_at to be set")
	}

	// if !actual.CreatedAt.Equal(post.CreatedAt) {
	// 	t.Fatal("publish should not change created_at")
	// }
}

func TestPostAPI_Publish_AlreadyPublished(t *testing.T) {
	app := NewTestApplication(t)

	editor := CreateEditorFixture(t, app)

	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		editor.User,
		category,
	)

	// 第一次 Publish
	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		editor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusOK,
	)

	// 第二次 Publish
	response = DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		editor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusConflict,
	)
}

func TestPostAPI_Publish_OtherEditorForbidden(t *testing.T) {
	app := NewTestApplication(t)

	owner := CreateEditorFixture(t, app)
	otherEditor := CreateEditorFixture(t, app)

	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		owner.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		otherEditor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusForbidden,
	)

	var actual postmodel.Post

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&actual).
		Error

	if err != nil {
		t.Fatalf(
			"query post failed: %v",
			err,
		)
	}

	if actual.Status != string(postmodel.PostStatusDraft) {
		t.Fatalf(
			"post should remain draft, got %q",
			actual.Status,
		)
	}
}

func TestPostAPI_Publish_AdminSuccess(t *testing.T) {
	app := NewTestApplication(t)

	admin := CreateAdminFixture(t, app)
	owner := CreateEditorFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		owner.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		admin.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusOK,
	)
}

func TestPostAPI_Publish_ViewerForbidden(t *testing.T) {
	app := NewTestApplication(t)

	owner := CreateEditorFixture(t, app)
	viewer := CreateViewerFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePostFixture(
		t,
		app.Server.DB(),
		owner.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		viewer.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusForbidden,
	)
}

func TestPostAPI_Archive_OwnerSuccess(t *testing.T) {
	app := NewTestApplication(t)

	editor := CreateEditorFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreatePublishedPostFixture(
		t,
		app.Server.DB(),
		editor.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/archive",
		nil,
		editor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusOK,
	)

	var actual postmodel.Post

	err := app.Server.DB().
		Where("id = ?", post.ID).
		First(&actual).
		Error

	if err != nil {
		t.Fatalf(
			"query post failed: %v",
			err,
		)
	}

	if actual.Status != string(postmodel.PostStatusArchived) {
		t.Fatalf(
			"expected archived, got %q",
			actual.Status,
		)
	}
	if actual.PublishedAt == nil {
		t.Fatal("archive should not clear published_at")
	}

}

func TestPostAPI_Publish_ArchivedForbidden(t *testing.T) {
	app := NewTestApplication(t)

	editor := CreateEditorFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreateArchivedPostFixture(
		t,
		app.Server.DB(),
		editor.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/publish",
		nil,
		editor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusConflict,
	)
}

func TestPostAPI_Archive_AlreadyArchived(t *testing.T) {
	app := NewTestApplication(t)

	editor := CreateEditorFixture(t, app)
	category := CreateCategoryFixture(t, app.Server.DB())

	post := CreateArchivedPostFixture(
		t,
		app.Server.DB(),
		editor.User,
		category,
	)

	response := DoJSONRequest(
		t,
		app.Server.HTTPHandler(),
		http.MethodPost,
		"/api/admin/v1/posts/"+post.ID+"/archive",
		nil,
		editor.Token,
	)

	AssertStatus(
		t,
		response,
		http.StatusBadRequest,
	)
}
