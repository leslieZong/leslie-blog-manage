package service

import (
	"context"
	"errors"
	appErrors "leslie-blog-server/internal/errors"
	categoryRepository "leslie-blog-server/internal/modules/category/repository"
	"leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/modules/post/repository"
	tagService "leslie-blog-server/internal/modules/tag/service"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/database"
	"leslie-blog-server/internal/pkg/ulid"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PostService 定义文章模块的业务能力。
//
// Handler 不应该直接依赖 Repository，
// 而应该依赖 Service。
type PostService interface {

	// Create 创建文章。
	Create(
		ctx context.Context,
		authorID string,
		title string,
		slug string,
		categoryID string,
		tagIDs []string,
		summary *string,
		content string,
		cover *string,
	) (*model.Post, error)

	// GetByID 根据文章 ID 获取文章。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.Post, error)

	// GetBySlug 根据文章 Slug 获取文章。
	GetBySlug(
		ctx context.Context,
		slug string,
	) (*model.Post, error)

	// 管理端分页查询文章
	List(
		ctx context.Context,
		query repository.PostListQuery,
	) ([]*model.Post, int64, error)

	// 公共页面分页查询文章
	ListPublished(
		ctx context.Context,
		query repository.PostListQuery,
	) ([]*model.Post, int64, error)

	// Public API 专用
	GetPublicByID(
		ctx context.Context,
		id string,
	) (*model.Post, error)

	GetPublicBySlug(
		ctx context.Context,
		slug string,
	) (*model.Post, error)

	// Update 更新文章。
	Update(
		ctx context.Context,
		actor auth.Actor,
		id string,
		title string,
		slug string,
		categoryID string,
		tagIDs []string,
		summary *string,
		content string,
		cover *string,
	) (*model.Post, error)

	// Delete 删除文章。
	//
	// 当前设计为软删除。
	Delete(
		ctx context.Context,
		actor auth.Actor,
		id string,
	) error

	// Publish 发布文章。
	Publish(
		ctx context.Context,
		actor auth.Actor,
		id string,
	) (*model.Post, error)

	// Archive 归档文章。
	Archive(
		ctx context.Context,
		actor auth.Actor,
		id string,
	) (*model.Post, error)

	// IncrementViewCount 增加文章阅读量。
	IncrementViewCount(
		ctx context.Context,
		id string,
	) error
}

// postService 是 PostService 的具体实现。
//
// 注意：
// 类型本身使用小写开头，说明它不需要暴露给其他 package。
// 外部只需要通过 PostService 接口使用它。
type postService struct {
	// repository 负责数据库操作。
	repo               repository.PostRepository
	postTagRepo        repository.PostTagRepository
	categoryRepo       categoryRepository.CategoryRepository
	tagService         tagService.TagService
	transactionManager *database.TransactionManager
}

func NewPostService(
	repo repository.PostRepository,
	postTagRepo repository.PostTagRepository,
	categoryRepo categoryRepository.CategoryRepository,
	tagService tagService.TagService,
	transactionManager *database.TransactionManager,
) PostService {

	return &postService{
		repo:               repo,
		postTagRepo:        postTagRepo,
		categoryRepo:       categoryRepo,
		tagService:         tagService,
		transactionManager: transactionManager,
	}
}

// normalizeIDs 归一化 ID 列表，去重并去空格。
func normalizeIDs(ids []string) []string {

	result := make([]string, 0, len(ids))

	seen := make(map[string]struct{})

	for _, id := range ids {

		id = strings.TrimSpace(id)

		if id == "" {
			continue
		}

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}

		result = append(result, id)
	}

	return result
}

// Create 创建一篇文章。
func (s *postService) Create(
	ctx context.Context,
	authorID string,
	title string,
	slug string,
	categoryID string,
	tagIDs []string,
	summary *string,
	content string,
	cover *string,
) (*model.Post, error) {

	// 1. 验证 Actor
	// 2. 验证标题
	// 3. 验证 Slug
	// 4. 验证 Category
	// 5. 清理 TagIDs
	// 6. 验证 TagIDs
	// 7. 检查 Slug 唯一
	// 8. 创建 Post
	// 9. 开启事务
	// 10. 创建 Post
	// 11. 创建 PostTag
	// 12. Commit

	// -------------------------------------------------------
	// 1. 基础参数校验
	// -------------------------------------------------------

	// --------------------------------------------------------
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post title is empty",
		)

	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post slug is empty",
		)
	}

	if content == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post content is empty",
		)
	}

	if authorID == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post author ID is empty",
		)
	}

	// ------------------------------------------------
	// 关键业务规则：
	//
	// 文章必须属于一个存在且启用的分类。
	// ------------------------------------------------
	if err := s.validateCategory(
		ctx,
		categoryID,
	); err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 2. 检查 slug 是否已经存在
	// -------------------------------------------------------

	existingPost, err := s.repo.FindBySlug(ctx, slug)

	if err == nil && existingPost != nil {
		return nil, appErrors.New(
			appErrors.ErrPostSlugExists,
			http.StatusConflict,
			"post slug already exists",
		)
	}

	// 没有找到记录，这是正常情况，可以继续创建。
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to check post slug",
			err,
		)
	}

	// 归一化 Tag ID 列表
	tagIDs = normalizeIDs(tagIDs)
	// 检查 Tag 是否存在且启用
	if _, err := s.tagService.FindByIDs(
		ctx,
		tagIDs,
	); err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 3. 创建 Post Model
	// -------------------------------------------------------

	post := &model.Post{
		ID:         ulid.New(),
		Title:      title,
		Slug:       slug,
		Summary:    summary,
		CategoryID: categoryID,
		Content:    content,
		Cover:      cover,

		// 创建文章默认是草稿。
		Status: string(model.PostStatusDraft),

		// 作者必须来自当前登录用户。
		//
		// 不能让前端提交 author_id。
		AuthorID: authorID,

		// 草稿没有发布时间。
		PublishedAt: nil,

		// 新文章阅读量从 0 开始。
		ViewCount: 0,
	}

	// -------------------------------------------------------
	// 4. 保存到数据库
	// -------------------------------------------------------
	err = s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			postRepo := repository.NewPostRepository(tx)

			postTagRepo := repository.NewPostTagRepository(tx)

			// 创建 Post
			if err := postRepo.Create(ctx, post); err != nil {
				return err
			}

			// 创建 Tag 关联
			if err := postTagRepo.ReplaceTags(
				ctx,
				post.ID,
				tagIDs,
			); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 5. 返回创建后的文章
	// -------------------------------------------------------
	// 刷新文章，确保关联的 Tag 刷新。
	post, err = s.repo.FindByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetByID 根据 ID 获取文章。
func (s *postService) GetByID(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post ID 不能为空",
		)
	}

	post, err := s.repo.FindByID(ctx, id)

	// =========================================================
	// 第三步：判断文章是否不存在
	// =========================================================
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query post",
			err,
		)
	}

	return post, nil
}

// GetBySlug 根据 slug 获取文章。
func (s *postService) GetBySlug(
	ctx context.Context,
	slug string,
) (*model.Post, error) {

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post slug 不能为空",
		)
	}

	post, err := s.repo.FindBySlug(ctx, slug)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query post",
			err,
		)
	}

	return post, nil
}

// List 获取文章列表。
func (s *postService) List(
	ctx context.Context,
	query repository.PostListQuery,
) ([]*model.Post, int64, error) {

	// Service 层不负责直接操作数据库。
	//
	// 它只负责：
	// 1. 接收业务查询条件
	// 2. 调用 Repository
	// 3. 返回查询结果
	//
	// 真正的 SQL 查询由 Repository 完成。
	posts, total, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *postService) ListPublished(
	ctx context.Context,
	query repository.PostListQuery,
) ([]*model.Post, int64, error) {
	query.Status = string(model.PostStatusPublished)
	posts, total, err := s.repo.FindPublishedPage(ctx, query)

	if err != nil {
		return nil, 0, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query published posts",
			err,
		)
	}

	return posts, total, nil
}

// GetPublicByID
//
// Public API 根据 ID 查询文章。
//
// 与 Admin 的 GetByID 最大区别：
//
// Admin:
//
//	可以查看 draft / published / archived
//
// Public:
//
//	只能查看 published
//
// 因此这里必须进行状态检查。
func (s *postService) GetPublicByID(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post id 不能为空",
		)
	}

	post, err := s.repo.FindByID(ctx, id)

	// 数据库没有找到
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	// 数据库发生其他错误
	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query post",
			err,
		)
	}

	// Public API 只允许 published
	if post.Status != string(model.PostStatusPublished) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	// 防止软删除文章被公开访问
	if post.DeletedAt.Valid {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	return post, nil
}

// GetPublicBySlug
//
// 根据 slug 获取博客前台文章。
//
// 例如：
//
// /api/v1/posts/slug/go-start
//
// 只有 published 状态的文章才允许返回。
func (s *postService) GetPublicBySlug(
	ctx context.Context,
	slug string,
) (*model.Post, error) {

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post slug cannot be empty",
		)
	}

	post, err := s.repo.FindBySlug(ctx, slug)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query post",
			err,
		)
	}

	// 非 published 状态对 Public API 不可见
	if post.Status != string(model.PostStatusPublished) {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	// 已软删除的数据也不能公开
	if post.DeletedAt.Valid {
		return nil, appErrors.New(
			appErrors.ErrPostNotFound,
			http.StatusNotFound,
			"post not found",
		)
	}

	return post, nil
}

func canModifyPost(
	post *model.Post,
	actor auth.Actor,
) bool {

	// 管理员可以操作所有文章。
	if actor.IsAdmin {
		return true
	}

	// 普通编辑只能操作自己创建的文章。
	return post.AuthorID == actor.UserID
}

// Update 更新文章内容。
func (s *postService) Update(
	ctx context.Context,
	actor auth.Actor,
	id string,
	title string,
	slug string,
	categoryID string,
	tagIDs []string,
	summary *string,
	content string,
	cover *string,
) (*model.Post, error) {

	// --------------------------------------------------
	// 1. 查询文章
	// --------------------------------------------------

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("文章 ID 不能为空")
	}
	post, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	// --------------------------------------------------
	// 2. 权限：检查当前用户是否可以修改这篇文章
	// --------------------------------------------------
	//
	// RBAC Middleware 已经检查：
	//
	//     user -> post:update
	//
	// 这里检查的是第二层：
	//
	//     这个用户能不能修改“这一篇”文章？
	//
	// Admin：
	//     可以修改任何文章
	//
	// Editor：
	//     只能修改自己的文章

	if !canModifyPost(post, actor) {
		return nil, appErrors.New(
			appErrors.ErrPostNotOwner,
			http.StatusForbidden,
			"you do not have permission to modify this post",
		)
	}

	// --------------------------------------------------
	// 3. 基础参数校验
	// --------------------------------------------------

	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, errors.New("文章标题不能为空")
	}

	if slug == "" {
		return nil, errors.New("文章 slug 不能为空")
	}

	if content == "" {
		return nil, errors.New("文章内容不能为空")
	}

	// --------------------------------------------------
	// 4. 校验 Category
	// --------------------------------------------------

	if err := s.validateCategory(
		ctx,
		categoryID,
	); err != nil {
		return nil, err
	}
	// 归一化 Tag ID 列表
	tagIDs = normalizeIDs(tagIDs)
	// 检查 Tag 是否存在且启用
	if _, err := s.tagService.FindByIDs(
		ctx,
		tagIDs,
	); err != nil {
		return nil, err
	}

	// --------------------------------------------------
	// 5. 检查 slug 是否与其他文章冲突
	// --------------------------------------------------

	existingPost, err := s.repo.FindBySlug(ctx, slug)

	if err == nil && existingPost != nil {

		// 如果找到的文章不是当前文章，
		// 说明 slug 已经被其他文章占用。
		if existingPost.ID != post.ID {
			return nil, errors.New("文章 slug 已存在")
		}
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// --------------------------------------------------
	// 6. 修改允许修改的字段
	// --------------------------------------------------

	post.Title = title
	post.Slug = slug
	post.Summary = summary
	post.Content = content
	post.Cover = cover
	post.CategoryID = categoryID

	// 注意：
	//
	// 这里没有修改：
	//
	// post.AuthorID
	// post.Status
	// post.PublishedAt
	// post.ViewCount
	//
	// 因为这些字段不是普通 Update 应该直接修改的。

	// --------------------------------------------------
	// 7. 保存
	// --------------------------------------------------

	err = s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			postRepo := repository.NewPostRepository(tx)

			postTagRepo := repository.NewPostTagRepository(tx)

			// 创建 Post
			if err := postRepo.Update(ctx, post); err != nil {
				return err
			}

			// 创建 Tag 关联
			if err := postTagRepo.ReplaceTags(
				ctx,
				post.ID,
				tagIDs,
			); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	// 刷新文章，确保关联的 Tag 刷新。
	post, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Publish 发布文章。
func (s *postService) Publish(
	ctx context.Context,
	actor auth.Actor,
	id string,
) (*model.Post, error) {

	// -------------------------------------------------------
	// 1. 参数校验
	// -------------------------------------------------------

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"post ID 不能为空",
		)
	}

	// -------------------------------------------------------
	// 2. 查询文章
	// -------------------------------------------------------

	post, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query post",
			err,
		)
	}

	// 第一层：资源所有权。
	if !canModifyPost(post, actor) {
		return nil, appErrors.New(
			appErrors.ErrPostNotOwner,
			http.StatusForbidden,
			"you do not have permission to publish this post",
		)
	}

	// -------------------------------------------------------
	// 3. 检查当前文章状态
	//
	// 当前设计：
	//
	// draft     → 可以发布
	// published → 不允许重复发布
	// archived  → 不允许直接发布
	// -------------------------------------------------------

	status := model.PostStatus(post.Status)
	if !status.CanPublish() {
		if status == model.PostStatusArchived {
			return nil, appErrors.New(
				appErrors.ErrPostArchived,
				http.StatusConflict,
				"archived post cannot be published",
			)
		}

		if status == model.PostStatusPublished {
			return nil, appErrors.New(
				appErrors.ErrPostAlreadyPublished,
				http.StatusConflict,
				"post is already published",
			)
		}
	}

	// -------------------------------------------------------
	// 4. 执行发布
	// -------------------------------------------------------

	now := time.Now()

	post.Status = string(model.PostStatusPublished)

	post.PublishedAt = &now

	// -------------------------------------------------------
	// 5. 保存
	// -------------------------------------------------------

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}
func (s *postService) Archive(
	ctx context.Context,
	actor auth.Actor,
	id string,
) (*model.Post, error) {

	post, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(
				appErrors.ErrPostNotFound,
				http.StatusNotFound,
				"post not found",
			)
		}

		return nil, err
	}

	// Ownership。
	if !canModifyPost(post, actor) {
		return nil, appErrors.New(
			appErrors.ErrPostNotOwner,
			http.StatusForbidden,
			"you do not have permission to archive this post",
		)
	}

	// 只有 published 状态才能归档。
	if !model.PostStatus(post.Status).CanArchive() {

		if post.Status == string(model.PostStatusArchived) {
			return nil, appErrors.New(
				appErrors.ErrPostArchived,
				http.StatusBadRequest,
				"post is already archived",
			)
		}

		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"only published post can be archived",
		)
	}

	post.Status = string(model.PostStatusArchived)

	if err := s.repo.Update(
		ctx,
		post,
	); err != nil {
		return nil, err
	}

	return post, nil
}

// Delete 删除文章。
//
// 当前系统使用软删除。
// Repository 内部会设置 deleted_at。
func (s *postService) Delete(
	ctx context.Context,
	actor auth.Actor,
	id string,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return errors.New("文章 ID 不能为空")
	}

	// 先确认文章存在。
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	// 权限第二层检查。
	if !canModifyPost(post, actor) {
		return appErrors.New(
			appErrors.ErrPostNotOwner,
			http.StatusForbidden,
			"you do not have permission to delete this post",
		)
	}

	// 执行软删除。
	return s.repo.Delete(ctx, id)
}

// IncrementViewCount 增加文章阅读量。
func (s *postService) IncrementViewCount(
	ctx context.Context,
	id string,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return errors.New("文章 ID 不能为空")
	}

	return s.repo.IncrementViewCount(ctx, id)
}

func (s *postService) validateCategory(
	ctx context.Context,
	categoryID string,
) error {

	// 检查 Category
	// --------------------------------------------------------

	categoryID = strings.TrimSpace(categoryID)

	if categoryID == "" {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category id cannot be empty",
		)
	}

	category, err := s.categoryRepo.FindByID(
		ctx,
		categoryID,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category not found",
		)
	}

	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query category",
			err,
		)
	}

	// 分类被禁用时，不允许新文章使用。
	if category.Status != 1 {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category is disabled",
		)
	}

	return nil
}
