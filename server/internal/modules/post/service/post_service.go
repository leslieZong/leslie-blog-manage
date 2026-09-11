package service

import (
	"context"
	"errors"
	"leslie-blog-server/internal/modules/post/model"
	"leslie-blog-server/internal/modules/post/repository"
	"leslie-blog-server/internal/pkg/ulid"
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

	// List 获取文章列表。
	List(
		ctx context.Context,
	) ([]*model.Post, error)

	// Update 更新文章。
	Update(
		ctx context.Context,
		id string,
		title string,
		slug string,
		summary *string,
		content string,
		cover *string,
	) (*model.Post, error)

	// Delete 删除文章。
	//
	// 当前设计为软删除。
	Delete(
		ctx context.Context,
		id string,
	) error

	// Publish 发布文章。
	Publish(
		ctx context.Context,
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
	repo repository.PostRepository
}

func NewPostService(
	repo repository.PostRepository,
) PostService {

	return &postService{
		repo: repo,
	}
}

// Create 创建一篇文章。
func (s *postService) Create(
	ctx context.Context,
	authorID string,
	title string,
	slug string,
	summary *string,
	content string,
	cover *string,
) (*model.Post, error) {

	// -------------------------------------------------------
	// 1. 基础参数校验
	// -------------------------------------------------------

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

	if authorID == "" {
		return nil, errors.New("文章作者不能为空")
	}

	// -------------------------------------------------------
	// 2. 检查 slug 是否已经存在
	// -------------------------------------------------------

	existingPost, err := s.repo.FindBySlug(ctx, slug)

	if err == nil && existingPost != nil {
		return nil, errors.New("文章 slug 已存在")
	}

	// 没有找到记录，这是正常情况，可以继续创建。
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// -------------------------------------------------------
	// 3. 创建 Post Model
	// -------------------------------------------------------

	post := &model.Post{
		ID:      ulid.New(),
		Title:   title,
		Slug:    slug,
		Summary: summary,
		Content: content,
		Cover:   cover,

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

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 5. 返回创建后的文章
	// -------------------------------------------------------

	return post, nil
}

// GetByID 根据 ID 获取文章。
func (s *postService) GetByID(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("文章 ID 不能为空")
	}

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
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
		return nil, errors.New("文章 slug 不能为空")
	}

	post, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// List 获取文章列表。
func (s *postService) List(
	ctx context.Context,
) ([]*model.Post, error) {

	posts, err := s.repo.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Update 更新文章内容。
func (s *postService) Update(
	ctx context.Context,
	id string,
	title string,
	slug string,
	summary *string,
	content string,
	cover *string,
) (*model.Post, error) {

	// -------------------------------------------------------
	// 1. 参数校验
	// -------------------------------------------------------

	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	content = strings.TrimSpace(content)

	if id == "" {
		return nil, errors.New("文章 ID 不能为空")
	}

	if title == "" {
		return nil, errors.New("文章标题不能为空")
	}

	if slug == "" {
		return nil, errors.New("文章 slug 不能为空")
	}

	if content == "" {
		return nil, errors.New("文章内容不能为空")
	}

	// -------------------------------------------------------
	// 2. 查询文章
	// -------------------------------------------------------

	post, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 3. 检查 slug 是否被其他文章占用
	// -------------------------------------------------------

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

	// -------------------------------------------------------
	// 4. 修改允许修改的字段
	// -------------------------------------------------------

	post.Title = title
	post.Slug = slug
	post.Summary = summary
	post.Content = content
	post.Cover = cover

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

	// -------------------------------------------------------
	// 5. 保存
	// -------------------------------------------------------

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// Publish 发布文章。
func (s *postService) Publish(
	ctx context.Context,
	id string,
) (*model.Post, error) {

	// -------------------------------------------------------
	// 1. 参数校验
	// -------------------------------------------------------

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("文章 ID 不能为空")
	}

	// -------------------------------------------------------
	// 2. 查询文章
	// -------------------------------------------------------

	post, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------
	// 3. 检查当前文章状态
	// -------------------------------------------------------

	status := model.PostStatus(post.Status)

	if status == model.PostStatusArchived {
		return nil, errors.New("已归档文章不能直接发布")
	}

	if status == model.PostStatusPublished {
		return nil, errors.New("文章已经发布")
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

// Delete 删除文章。
//
// 当前系统使用软删除。
// Repository 内部会设置 deleted_at。
func (s *postService) Delete(
	ctx context.Context,
	id string,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return errors.New("文章 ID 不能为空")
	}

	// 先确认文章存在。
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
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
