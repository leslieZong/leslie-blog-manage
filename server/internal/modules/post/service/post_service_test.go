package service

import (
	"context"
	appErrors "leslie-blog-server/internal/errors"
	categorymodel "leslie-blog-server/internal/modules/category/model"
	categoryRepository "leslie-blog-server/internal/modules/category/repository"
	"leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/modules/post/repository"
	"leslie-blog-server/internal/pkg/audit"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/cache"
	"leslie-blog-server/internal/pkg/ulid"
	"testing"
	"time"

	"gorm.io/gorm"
)

type mockPostRepository struct {
	// 用于模拟 FindByID 返回的数据。
	post *model.Post

	// 用于模拟 FindBySlug 返回的数据。
	slugPost *model.Post

	// 模拟 Repository 返回错误。
	err error

	// 记录 Create 是否被调用。
	createCalled bool

	// 记录 Update 是否被调用。
	updateCalled bool

	// 保存 Create 接收到的 Post。
	createdPost *model.Post

	// 保存 Update 接收到的 Post。
	updatedPost *model.Post

	// 保存 Delete 接收到的 ID。
	deletedID string

	// 保存 IncrementViewCount 接收到的 ID。
	incrementViewCountID string

	// 用于模拟 FindAll 返回的数据。
	posts []*model.Post
}

func (m *mockPostRepository) Create(
	ctx context.Context,
	post *model.Post,
) error {

	m.createCalled = true
	m.createdPost = post

	return m.err
}

func (m *mockPostRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	return m.post, m.err
}

func (m *mockPostRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.Post, error) {

	return m.slugPost, m.err
}

func (m *mockPostRepository) FindAll(
	ctx context.Context,
) ([]*model.Post, error) {

	return m.posts, m.err
}

func (m *mockPostRepository) Update(
	ctx context.Context,
	post *model.Post,
) error {

	m.updateCalled = true
	m.updatedPost = post

	return m.err
}

func (m *mockPostRepository) FindPublished(
	ctx context.Context,
) ([]*model.Post, error) {

	return m.posts, m.err
}

func (m *mockPostRepository) Delete(
	ctx context.Context,
	id string,
) error {

	m.deletedID = id

	return m.err
}

func (m *mockPostRepository) FindPage(
	ctx context.Context,
	query repository.PostListQuery,
) ([]*model.Post, int64, error) {

	return nil, 0, m.err
}

func (m *mockPostRepository) FindPublishedPage(
	ctx context.Context,
	query repository.PostListQuery,
) ([]*model.Post, int64, error) {

	return nil, 0, nil
}

func (m *mockPostRepository) FindFeaturedPage(
	ctx context.Context,
	query repository.PostListQuery,
) ([]*model.Post, int64, error) {

	return nil, 0, m.err
}

func (m *mockPostRepository) IncrementViewCount(
	ctx context.Context,
	id string,
) error {

	m.incrementViewCountID = id

	return m.err
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type mockCache struct {
	deletedKeys []string

	getValue string
	getErr   error

	setKey   string
	setValue string
	setTTL   time.Duration
}

func (m *mockCache) Get(
	ctx context.Context,
	key string,
) (string, error) {

	return m.getValue, m.getErr
}

func (m *mockCache) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {

	m.setKey = key
	m.setValue = value
	m.setTTL = ttl

	return nil
}

func (m *mockCache) Delete(
	ctx context.Context,
	key string,
) error {

	m.deletedKeys = append(
		m.deletedKeys,
		key,
	)

	return nil
}

type mockAuditRecorder struct {
	successEvents []audit.Event
	failureEvents []audit.Event
}

func (m *mockAuditRecorder) RecordSuccess(
	ctx context.Context,
	event audit.Event,
) error {

	m.successEvents = append(
		m.successEvents,
		event,
	)

	return nil
}

func (m *mockAuditRecorder) RecordFailure(
	ctx context.Context,
	event audit.Event,
	err error,
) error {

	m.failureEvents = append(
		m.failureEvents,
		event,
	)

	return nil
}

func newTestPost() *model.Post {
	return &model.Post{
		ID:       "post-001",
		Title:    "Go Testing",
		Slug:     "go-testing",
		AuthorID: "user-001",
		Status:   string(model.PostStatusDraft),
	}
}

func newTestActor() auth.Actor {
	return auth.Actor{
		UserID:  "user-001",
		IsAdmin: false,
	}
}

func TestPostService_GetByID(t *testing.T) {

	expected := newTestPost()

	repo := &mockPostRepository{
		post: expected,
	}

	service := &postService{
		repo: repo,
	}

	ctx := context.Background()

	actual, err := service.GetByID(
		ctx,
		expected.ID,
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

	if actual.ID != expected.ID {
		t.Fatalf(
			"expected id %s, got %s",
			expected.ID,
			actual.ID,
		)
	}
}

func TestPostService_GetByID_NotFound(t *testing.T) {

	repo := &mockPostRepository{
		err: gorm.ErrRecordNotFound,
	}

	service := &postService{
		repo: repo,
	}

	_, err := service.GetByID(
		context.Background(),
		"not-found",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr := appErrors.FromError(err)

	if appErr.Code != appErrors.ErrPostNotFound {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrPostNotFound,
			appErr.Code,
		)
	}
}

type mockCategoryRepository struct {
	category *categorymodel.Category
	err      error
}

func (m *mockCategoryRepository) FindByID(
	ctx context.Context,
	id string,
) (*categorymodel.Category, error) {
	return m.category, m.err
}
func (m *mockCategoryRepository) FindByName(
	ctx context.Context,
	name string,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func (m *mockCategoryRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func (m *mockCategoryRepository) FindAll(
	ctx context.Context,
) ([]*categorymodel.Category, error) {
	return nil, m.err
}

func (m *mockCategoryRepository) FindPublicAll(
	ctx context.Context,
) ([]*categorymodel.Category, error) {
	return nil, m.err
}

func (m *mockCategoryRepository) FindPublicBySlug(
	ctx context.Context,
	slug string,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func (m *mockCategoryRepository) FindPage(
	ctx context.Context,
	query categoryRepository.CategoryListQuery,
) ([]*categorymodel.Category, int64, error) {
	return nil, 0, m.err
}

func (m *mockCategoryRepository) Create(
	ctx context.Context,
	category *categorymodel.Category,
) error {
	return m.err
}
func (m *mockCategoryRepository) Update(
	ctx context.Context,
	category *categorymodel.Category,
) error {
	return m.err
}
func (m *mockCategoryRepository) Delete(
	ctx context.Context,
	id string,
) error {
	return m.err
}
func (m *mockCategoryRepository) CountPosts(
	ctx context.Context,
	categoryID string,
) (int64, error) {
	return 0, m.err
}

type mockCategoryService struct {
	category *categorymodel.Category
	err      error
}

func (m *mockCategoryService) Create(
	ctx context.Context,
	name string,
	slug string,
	description *string,
	sort int,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func (m *mockCategoryService) GetByID(
	ctx context.Context,
	id string,
) (*categorymodel.Category, error) {

	return m.category, m.err
}

func (m *mockCategoryService) List(
	ctx context.Context,
) ([]*categorymodel.Category, error) {
	return []*categorymodel.Category{m.category}, m.err
}

func (m *mockCategoryService) Update(
	ctx context.Context,
	id string,
	name string,
	slug string,
	description *string,
	sort int,
	status int8,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func (m *mockCategoryService) Delete(
	ctx context.Context,
	id string,
) error {
	return m.err
}

func (m *mockCategoryService) ListPage(
	ctx context.Context,
	query categoryRepository.CategoryListQuery,
) ([]*categorymodel.Category, int64, error) {
	return []*categorymodel.Category{m.category}, 0, m.err
}
func (m *mockCategoryService) ListPublic(
	ctx context.Context,
) ([]*categorymodel.Category, error) {
	return []*categorymodel.Category{m.category}, m.err
}

func (m *mockCategoryService) GetPublicBySlug(
	ctx context.Context,
	slug string,
) (*categorymodel.Category, error) {
	return m.category, m.err
}

func TestPostService_Create_CategoryNotFound(
	t *testing.T,
) {

	categoryRepo := &mockCategoryRepository{
		err: gorm.ErrRecordNotFound,
	}

	// 这里按照你当前 postService 实际依赖
	// 初始化其它 Mock。
	//
	// 重点不是构造函数形式，
	// 而是验证 Create 在 Category 不存在时失败。
	repo := &mockPostRepository{}
	cacheMock := &mockCache{}

	auditMock := &mockAuditRecorder{}
	postService := &postService{
		repo:          repo,
		cache:         cacheMock,
		auditRecorder: auditMock,
		categoryRepo:  categoryRepo,
	}

	_, err := postService.Create(
		context.Background(),
		ulid.New(),
		"title",
		"slug",
		"content",
		[]string{"category"},
		nil,
		"author",
		nil,
	)

	// 期望：
	//
	// Category 查询失败
	// ↓
	// Create 失败
	// ↓
	// 不应该创建 Post
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	appErr := appErrors.FromError(err)

	if appErr.Code != appErrors.ErrCategoryNotFound {
		t.Fatalf(
			"expected error %v, got %v",
			appErrors.ErrCategoryNotFound,
			appErr.Code,
		)
	}
	if repo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}
	if len(auditMock.successEvents) != 0 {
		t.Fatalf(
			"expected 0 audit event, got %d",
			len(auditMock.successEvents),
		)
	}
	if len(cacheMock.deletedKeys) != 0 {
		t.Fatalf(
			"expected 0 deleted cache key, got %d",
			len(cacheMock.deletedKeys),
		)
	}
}

func TestPostService_Publish_Author(
	t *testing.T,
) {

	post := newTestPost()

	repo := &mockPostRepository{
		post: post,
	}

	cacheMock := &mockCache{}

	auditMock := &mockAuditRecorder{}

	service := &postService{
		repo:          repo,
		cache:         cacheMock,
		auditRecorder: auditMock,
	}

	actor := auth.Actor{
		UserID:  post.AuthorID,
		IsAdmin: false,
	}

	_, err := service.Publish(
		context.Background(),
		actor,
		post.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if !repo.updateCalled {
		t.Fatal(
			"expected repository Update to be called",
		)
	}

	if repo.updatedPost.Status != string(
		model.PostStatusPublished,
	) {
		t.Fatalf(
			"expected status published, got %s",
			repo.updatedPost.Status,
		)
	}

	if len(cacheMock.deletedKeys) != 1 {
		t.Fatalf(
			"expected 1 deleted cache key, got %d",
			len(cacheMock.deletedKeys),
		)
	}

	if cacheMock.deletedKeys[0] != cache.KeyHome {
		t.Fatalf(
			"expected home cache key %s, got %s",
			cache.KeyHome,
			cacheMock.deletedKeys[0],
		)
	}
	if len(auditMock.successEvents) != 1 {
		t.Fatalf(
			"expected 1 audit event, got %d",
			len(auditMock.successEvents),
		)
	}

	event := auditMock.successEvents[0]

	if event.Action != audit.ActionPostPublish {
		t.Fatalf(
			"expected action %s, got %s",
			audit.ActionPostPublish,
			event.Action,
		)
	}

	if event.Resource != "post" {
		t.Fatalf(
			"expected resource post, got %s",
			event.Resource,
		)
	}

	if event.ResourceID != post.ID {
		t.Fatalf(
			"expected resource id %s, got %s",
			post.ID,
			event.ResourceID,
		)
	}
}

func TestPostService_Publish_NonAuthor(
	t *testing.T,
) {

	post := newTestPost()

	repo := &mockPostRepository{
		post: post,
	}
	cacheMock := &mockCache{}
	service := &postService{
		repo:  repo,
		cache: cacheMock,
	}

	actor := auth.Actor{
		UserID:  "another-user",
		IsAdmin: false,
	}

	_, err := service.Publish(
		context.Background(),
		actor,
		post.ID,
	)

	if err == nil {
		t.Fatal(
			"expected forbidden error, got nil",
		)
	}

	appErr := appErrors.FromError(err)

	if appErr.Code != appErrors.ErrPostNotOwner {
		t.Fatalf(
			"expected forbidden code %d, got %d",
			appErrors.ErrPostNotOwner,
			appErr.Code,
		)
	}

	// 最重要的一点：
	//
	// 权限失败以后，
	// 不应该执行数据库 Update。
	if repo.updateCalled {
		t.Fatal(
			"repository Update should not be called",
		)
	}
}

func TestPostService_Publish_Admin(
	t *testing.T,
) {

	post := newTestPost()

	repo := &mockPostRepository{
		post: post,
	}

	cacheMock := &mockCache{}
	auditMock := &mockAuditRecorder{}
	service := &postService{
		repo:          repo,
		cache:         cacheMock,
		auditRecorder: auditMock,
	}

	actor := auth.Actor{
		UserID:  "admin-001",
		IsAdmin: true,
	}

	_, err := service.Publish(
		context.Background(),
		actor,
		post.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected admin to publish, got %v",
			err,
		)
	}

	if !repo.updateCalled {
		t.Fatal(
			"expected Update to be called",
		)
	}
}

func TestPostService_Publish_Archived(
	t *testing.T,
) {

	post := newTestPost()

	post.Status = string(
		model.PostStatusArchived,
	)

	repo := &mockPostRepository{
		post: post,
	}

	service := &postService{
		repo: repo,
	}

	actor := auth.Actor{
		UserID:  post.AuthorID,
		IsAdmin: false,
	}

	_, err := service.Publish(
		context.Background(),
		actor,
		post.ID,
	)

	if err == nil {
		t.Fatal(
			"expected archived post to fail",
		)
	}

	appErr := appErrors.FromError(err)

	if appErr.Code != appErrors.ErrPostArchived {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrPostArchived,
			appErr.Code,
		)
	}

	if repo.updateCalled {
		t.Fatal(
			"archived post should not be updated",
		)
	}
}

func TestPostService_Archive_Author(
	t *testing.T,
) {

	post := newTestPost()

	post.Status = string(
		model.PostStatusPublished,
	)

	repo := &mockPostRepository{
		post: post,
	}

	cacheMock := &mockCache{}

	auditMock := &mockAuditRecorder{}

	service := &postService{
		repo:          repo,
		cache:         cacheMock,
		auditRecorder: auditMock,
	}

	actor := auth.Actor{
		UserID:  post.AuthorID,
		IsAdmin: false,
	}

	_, err := service.Archive(
		context.Background(),
		actor,
		post.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if repo.updatedPost.Status != string(
		model.PostStatusArchived,
	) {
		t.Fatalf(
			"expected archived, got %s",
			repo.updatedPost.Status,
		)
	}
}

func TestPostService_Delete_Author(
	t *testing.T,
) {

	post := newTestPost()

	repo := &mockPostRepository{
		post: post,
	}

	cacheMock := &mockCache{}

	service := &postService{
		repo:  repo,
		cache: cacheMock,
	}

	actor := auth.Actor{
		UserID:  post.AuthorID,
		IsAdmin: false,
	}

	err := service.Delete(
		context.Background(),
		actor,
		post.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if repo.deletedID != post.ID {
		t.Fatalf(
			"expected deleted id %s, got %s",
			post.ID,
			repo.deletedID,
		)
	}
}

func TestPostService_Delete_NonAuthor(
	t *testing.T,
) {

	post := newTestPost()

	repo := &mockPostRepository{
		post: post,
	}

	cacheMock := &mockCache{}

	service := &postService{
		repo:  repo,
		cache: cacheMock,
	}

	actor := auth.Actor{
		UserID:  "another-user",
		IsAdmin: false,
	}

	err := service.Delete(
		context.Background(),
		actor,
		post.ID,
	)

	if err == nil {
		t.Fatal(
			"expected forbidden error",
		)
	}

	appErr := appErrors.FromError(err)

	if appErr.Code != appErrors.ErrPostNotOwner {
		t.Fatalf(
			"expected error code %d, got %d",
			appErrors.ErrPostNotOwner,
			appErr.Code,
		)
	}

	if repo.deletedID != "" {
		t.Fatal(
			"Delete should not be called",
		)
	}
}

// Table-Driven Test : 可以一次测试多个场景
// TestPostService_Publish 测试发布文章。
func TestPostService_Publish(t *testing.T) {

	tests := []struct {
		name      string
		actor     auth.Actor
		post      *model.Post
		wantError bool
	}{
		{
			name: "author can publish draft",
			actor: auth.Actor{
				UserID:  "user-001",
				IsAdmin: false,
			},
			post: &model.Post{
				ID:       "post-001",
				AuthorID: "user-001",
				Status:   "draft",
			},
			wantError: false,
		},
		{
			name: "non-author cannot publish",
			actor: auth.Actor{
				UserID:  "user-999",
				IsAdmin: false,
			},
			post: &model.Post{
				ID:       "post-001",
				AuthorID: "user-001",
				Status:   "draft",
			},
			wantError: true,
		},
		{
			name: "admin can publish",
			actor: auth.Actor{
				UserID:  "admin-001",
				IsAdmin: true,
			},
			post: &model.Post{
				ID:       "post-001",
				AuthorID: "user-001",
				Status:   "draft",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			repo := &mockPostRepository{
				post: tt.post,
			}
			cache := &mockCache{}
			auditRecorder := &mockAuditRecorder{}

			service := &postService{
				repo:          repo,
				cache:         cache,
				auditRecorder: auditRecorder,
			}

			_, err := service.Publish(
				context.Background(),
				tt.actor,
				tt.post.ID,
			)

			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf(
					"expected no error, got %v",
					err,
				)
			}
		})
	}
}
