package service

import (
	"errors"
	"strings"

	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo     *repository.PostRepo
	categoryRepo *repository.CategoryRepo
}

func NewPostService(postRepo *repository.PostRepo, categoryRepo *repository.CategoryRepo) *PostService {
	return &PostService{postRepo: postRepo, categoryRepo: categoryRepo}
}

// ListAdmin 后台列表（扁平化为目标 locale）
func (s *PostService) ListAdmin(f repository.PostFilter, locale string, page, pageSize int) ([]model.PostDTO, int64, error) {
	posts, total, err := s.postRepo.List(f, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return s.toDTOs(posts, locale), total, nil
}

// GetAdmin 后台详情
func (s *PostService) GetAdmin(id uint, locale string) (*model.PostDTO, error) {
	p, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	dto := s.toDTO(*p, locale)
	return &dto, nil
}

// GetBySlug 公开详情（仅已发布 + 目标 locale 翻译），自增浏览量
func (s *PostService) GetBySlug(slug, locale string) (*model.PostDTO, error) {
	p, err := s.postRepo.FindBySlug(slug, locale)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	_ = s.postRepo.IncView(p.ID)
	p.ViewCount++
	dto := s.toDTO(*p, locale)
	return &dto, nil
}

// ListPublished 公开列表
func (s *PostService) ListPublished(f repository.PostFilter, locale string, page, pageSize int) ([]model.PostDTO, int64, error) {
	posts, total, err := s.postRepo.ListPublished(f, locale, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return s.toDTOs(posts, locale), total, nil
}

// Create 创建文章
func (s *PostService) Create(req model.PostReq) error {
	locale := normLocale(req.Locale)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = slugify(req.Title)
	}

	p := model.Post{
		Slug:           slug,
		Cover:          req.Cover,
		CategoryID:     req.CategoryID,
		Status:         req.Status,
		IsTop:          req.IsTop,
		SEOTitle:       req.SEOTitle,
		SEODescription: req.SEODescription,
		SEOKeywords:    req.SEOKeywords,
	}

	tr := model.PostTranslation{
		Locale:  locale,
		Title:   req.Title,
		Summary: req.Summary,
		Content: req.Content,
	}
	return s.postRepo.Create(&p, []model.PostTranslation{tr}, req.Tags)
}

// Update 更新文章（upsert 翻译）
func (s *PostService) Update(id uint, req model.PostReq) error {
	p, err := s.postRepo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	locale := normLocale(req.Locale)
	if s2 := strings.TrimSpace(req.Slug); s2 != "" {
		p.Slug = s2
	}
	p.Cover = req.Cover
	p.CategoryID = req.CategoryID
	p.Status = req.Status
	p.IsTop = req.IsTop
	p.SEOTitle = req.SEOTitle
	p.SEODescription = req.SEODescription
	p.SEOKeywords = req.SEOKeywords

	tr := model.PostTranslation{
		PostID:  p.ID,
		Locale:  locale,
		Title:   req.Title,
		Summary: req.Summary,
		Content: req.Content,
	}
	return s.postRepo.Update(p, []model.PostTranslation{tr}, req.Tags)
}

func (s *PostService) Delete(id uint) error { return s.postRepo.Delete(id) }

func (s *PostService) Publish(id uint) error            { return s.postRepo.UpdateStatus(id, 1) }
func (s *PostService) Unpublish(id uint) error          { return s.postRepo.UpdateStatus(id, 2) }
func (s *PostService) SetTop(id uint, isTop bool) error { return s.postRepo.SetTop(id, isTop) }

// Search 公开检索
func (s *PostService) Search(keyword, locale string, page, pageSize int) ([]model.PostDTO, int64, error) {
	posts, total, err := s.postRepo.Search(keyword, locale, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return s.toDTOs(posts, locale), total, nil
}

// ===== 映射 =====

func (s *PostService) toDTOs(posts []model.Post, locale string) []model.PostDTO {
	out := make([]model.PostDTO, 0, len(posts))
	for i := range posts {
		out = append(out, s.toDTO(posts[i], locale))
	}
	return out
}

// toDTO 把 model.Post（含 Translations）扁平化为目标 locale 的 DTO
func (s *PostService) toDTO(p model.Post, locale string) model.PostDTO {
	dto := model.PostDTO{
		ID:             p.ID,
		Slug:           p.Slug,
		Cover:          p.Cover,
		CategoryID:     p.CategoryID,
		CategoryName:   p.Category.Name,
		Status:         p.Status,
		IsTop:          p.IsTop,
		ViewCount:      p.ViewCount,
		CommentCount:   p.CommentCount,
		SEOTitle:       p.SEOTitle,
		SEODescription: p.SEODescription,
		SEOKeywords:    p.SEOKeywords,
		PublishedAt:    p.PublishedAt,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	// 取目标 locale 翻译，无则回退第一条
	if tr := pickTranslation(p.Translations, locale); tr != nil {
		dto.Title = tr.Title
		dto.Summary = tr.Summary
		dto.Content = tr.Content
	}
	// tags 转 []string
	if len(p.Tags) > 0 {
		dto.Tags = make([]string, 0, len(p.Tags))
		for _, t := range p.Tags {
			dto.Tags = append(dto.Tags, t.Name)
		}
	}
	return dto
}

func pickTranslation(ts []model.PostTranslation, locale string) *model.PostTranslation {
	for i := range ts {
		if ts[i].Locale == locale {
			return &ts[i]
		}
	}
	if len(ts) > 0 {
		return &ts[0]
	}
	return nil
}

func normLocale(l string) string {
	switch l {
	case "en-US", "en":
		return "en-US"
	default:
		return "zh-CN"
	}
}

// slugify 简易 slug 生成（中文保留，空格转 -）
func slugify(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "post"
	}
	out := strings.Builder{}
	for _, r := range s {
		switch {
		case r == ' ' || r == '_':
			out.WriteRune('-')
		default:
			out.WriteRune(r)
		}
	}
	return strings.ToLower(out.String())
}
