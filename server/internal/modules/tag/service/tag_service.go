package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/tag/model"
	"leslie-blog-server/internal/modules/tag/repository"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

// TagService 定义 Tag 业务能力。
//
// Handler 不直接调用 Repository。
//
// 正确关系：
//
// Handler
//
//	↓
//
// Service
//
//	↓
//
// Repository
//
//	↓
//
// MySQL
type TagService interface {

	// 创建 Tag。
	Create(
		ctx context.Context,
		name string,
		slug string,
		description *string,
	) (*model.Tag, error)

	// 根据 ID 查询。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.Tag, error)

	// 查询全部 Tag。
	List(
		ctx context.Context,
	) ([]*model.Tag, error)

	// 更新 Tag。
	Update(
		ctx context.Context,
		id string,
		name string,
		slug string,
		description *string,
		status int8,
	) (*model.Tag, error)

	// 删除 Tag。
	Delete(
		ctx context.Context,
		id string,
	) error

	// 根据多个 ID 查询。
	//
	// 这个方法主要给 Post Service 使用。
	FindByIDs(
		ctx context.Context,
		ids []string,
	) ([]*model.Tag, error)
}

// tagService 是 TagService 的具体实现。
type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建 Tag Service。
func NewTagService(
	tagRepo repository.TagRepository,
) TagService {

	return &tagService{
		tagRepo: tagRepo,
	}
}

// Create 创建 Tag。
func (s *tagService) Create(
	ctx context.Context,
	name string,
	slug string,
	description *string,
) (*model.Tag, error) {

	// -----------------------------------------------------
	// 1. 基础数据清洗
	// -----------------------------------------------------

	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	// -----------------------------------------------------
	// 2. 基础业务校验
	// -----------------------------------------------------

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"tag name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"tag slug is required",
		)
	}

	// -----------------------------------------------------
	// 3. 检查名称是否重复
	// -----------------------------------------------------

	existing, err := s.tagRepo.FindByName(ctx, name)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tag name already exists",
		)
	}

	if err != nil && !isRecordNotFound(err) {
		return nil, err
	}

	// -----------------------------------------------------
	// 4. 检查 slug 是否重复
	// -----------------------------------------------------

	existing, err = s.tagRepo.FindBySlug(ctx, slug)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tag slug already exists",
		)
	}

	if err != nil && !isRecordNotFound(err) {
		return nil, err
	}

	// -----------------------------------------------------
	// 5. 创建 Model
	// -----------------------------------------------------

	tag := &model.Tag{
		ID:          ulid.New(),
		Name:        name,
		Slug:        slug,
		Description: description,
		Status:      1,
	}

	// -----------------------------------------------------
	// 6. 保存数据库
	// -----------------------------------------------------

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetByID 根据 ID 查询 Tag。
func (s *tagService) GetByID(
	ctx context.Context,
	id string,
) (*model.Tag, error) {
	// -----------------------------------------------------
	// 1. 校验 ID
	// -----------------------------------------------------

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tag id is required",
		)
	}

	return s.tagRepo.FindByID(ctx, id)
}

// List 查询所有 Tag。
func (s *tagService) List(
	ctx context.Context,
) ([]*model.Tag, error) {

	return s.tagRepo.FindAll(ctx)
}

// Update 更新 Tag。
func (s *tagService) Update(
	ctx context.Context,
	id string,
	name string,
	slug string,
	description *string,
	status int8,
) (*model.Tag, error) {

	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	// -----------------------------------------------------
	// 1. 查询原 Tag
	// -----------------------------------------------------

	tag, err := s.tagRepo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------
	// 2. 校验名称
	// -----------------------------------------------------

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tag name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tag slug is required",
		)
	}

	// -----------------------------------------------------
	// 3. 修改字段
	// -----------------------------------------------------

	tag.Name = name
	tag.Slug = slug
	tag.Description = description
	tag.Status = status

	// -----------------------------------------------------
	// 4. 保存
	// -----------------------------------------------------

	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete 删除 Tag。
func (s *tagService) Delete(
	ctx context.Context,
	id string,
) error {

	// -----------------------------------------------------
	// 1. 确认 Tag 存在
	// -----------------------------------------------------

	_, err := s.tagRepo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	// -----------------------------------------------------
	// 2. 执行软删除
	// -----------------------------------------------------

	return s.tagRepo.Delete(ctx, id)
}

// FindByIDs 查询多个 Tag。
func (s *tagService) FindByIDs(
	ctx context.Context,
	ids []string,
) ([]*model.Tag, error) {

	// -----------------------------------------------------
	// 1. 去掉空字符串
	// -----------------------------------------------------

	cleanIDs := make([]string, 0, len(ids))

	seen := make(map[string]struct{})

	for _, id := range ids {

		id = strings.TrimSpace(id)

		if id == "" {
			continue
		}

		// -------------------------------------------------
		// 防止重复 Tag ID
		// -------------------------------------------------

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}

		cleanIDs = append(cleanIDs, id)
	}

	// 没有 Tag。
	if len(cleanIDs) == 0 {
		return []*model.Tag{}, nil
	}

	// -----------------------------------------------------
	// 2. 查询数据库
	// -----------------------------------------------------

	tags, err := s.tagRepo.FindByIDs(ctx, cleanIDs)

	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------
	// 3. 判断是否全部存在
	// -----------------------------------------------------

	if len(tags) != len(cleanIDs) {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"one or more tags do not exist",
		)
	}

	return tags, nil
}

// isRecordNotFound 判断 GORM 是否返回“记录不存在”。
func isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
