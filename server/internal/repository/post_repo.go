package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type PostRepo struct{ db *gorm.DB }

func NewPostRepo(db *gorm.DB) *PostRepo { return &PostRepo{db: db} }

// PostFilter 文章查询过滤
type PostFilter struct {
	Keyword    string
	CategoryID uint
	Status     *int8
}

// List 分页查询（admin：预加载全部翻译 + 标签 + 分类）
func (r *PostRepo) List(f PostFilter, page, pageSize int) ([]model.Post, int64, error) {
	var total int64
	q := r.applyFilter(r.db.Model(&model.Post{}), f)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var posts []model.Post
	q = r.applyFilter(r.db.Model(&model.Post{}), f).
		Preload("Category").
		Preload("Tags").
		Preload("Translations").
		Order("is_top DESC, created_at DESC")
	q = paginate(q, page, pageSize)
	if err := q.Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// ListPublished 公开列表：只取已发布，按 locale 取单条翻译
func (r *PostRepo) ListPublished(f PostFilter, locale string, page, pageSize int) ([]model.Post, int64, error) {
	published := int8(1)
	f.Status = &published
	posts, total, err := r.List(f, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	// 仅保留目标 locale 的翻译
	for i := range posts {
		posts[i].Translations = pickLocale(posts[i].Translations, locale)
	}
	return posts, total, nil
}

func (r *PostRepo) applyFilter(q *gorm.DB, f PostFilter) *gorm.DB {
	if f.Keyword != "" {
		q = q.Where("id IN (?)",
			r.db.Table("post_translations").Select("post_id").
				Where("title LIKE ? OR content LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%"),
		)
	}
	if f.CategoryID != 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	return q
}

// FindByID 详情（admin：全翻译；按 locale 单条用 FindByIDLocale）
func (r *PostRepo) FindByID(id uint) (*model.Post, error) {
	var p model.Post
	if err := r.db.
		Preload("Category").
		Preload("Tags").
		Preload("Translations").
		First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySlug 已发布文章 + 目标 locale 翻译（公开详情）
func (r *PostRepo) FindBySlug(slug, locale string) (*model.Post, error) {
	var p model.Post
	if err := r.db.
		Preload("Category").
		Preload("Tags").
		Preload("Translations", "locale = ?", locale).
		Where("slug = ? AND status = ?", slug, int8(1)).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Create 事务：创建文章 + 翻译 + 标签
func (r *PostRepo) Create(p *model.Post, translations []model.PostTranslation, tags []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		for i := range translations {
			translations[i].PostID = p.ID
		}
		if len(translations) > 0 {
			if err := tx.Create(&translations).Error; err != nil {
				return err
			}
		}
		if err := syncTags(tx, p.ID, tags); err != nil {
			return err
		}
		return nil
	})
}

// Update 事务：更新文章 + upsert 翻译 + 替换标签
func (r *PostRepo) Update(p *model.Post, translations []model.PostTranslation, tags []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(p).Error; err != nil {
			return err
		}
		for i := range translations {
			translations[i].PostID = p.ID
			// upsert：按 post_id+locale 存在则更新
			var exist model.PostTranslation
			if err := tx.Where("post_id = ? AND locale = ?", p.ID, translations[i].Locale).
				First(&exist).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if e := tx.Create(&translations[i]).Error; e != nil {
						return e
					}
					continue
				}
				return err
			}
			translations[i].ID = exist.ID
			if e := tx.Save(&translations[i]).Error; e != nil {
				return e
			}
		}
		return syncTags(tx, p.ID, tags)
	})
}

// Delete 事务：删除文章 + 翻译 + 标签
func (r *PostRepo) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", id).Delete(&model.PostTranslation{}).Error; err != nil {
			return err
		}
		if err := tx.Where("post_id = ?", id).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Post{}, id).Error
	})
}

// UpdateStatus 更新状态（发布时设置 published_at）
func (r *PostRepo) UpdateStatus(id uint, status int8) error {
	updates := map[string]interface{}{"status": status}
	if status == 1 {
		updates["published_at"] = gorm.Expr("COALESCE(published_at, NOW())")
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

// SetTop 设置/取消置顶
func (r *PostRepo) SetTop(id uint, isTop bool) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("is_top", isTop).Error
}

// IncView 浏览量 +1
func (r *PostRepo) IncView(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// Search 全文检索（标题 / 内容 / 标签），返回匹配文章
func (r *PostRepo) Search(keyword, locale string, page, pageSize int) ([]model.Post, int64, error) {
	return r.ListPublished(PostFilter{Keyword: keyword}, locale, page, pageSize)
}

// syncTags 用 tags 列表替换某文章的标签（先删后插）
func syncTags(tx *gorm.DB, postID uint, tags []string) error {
	if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(tags))
	var rows []model.PostTag
	for _, t := range tags {
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		rows = append(rows, model.PostTag{PostID: postID, Name: t})
	}
	if len(rows) > 0 {
		return tx.Create(&rows).Error
	}
	return nil
}

// pickLocale 取目标 locale 的翻译；若无则回退第一条
func pickLocale(ts []model.PostTranslation, locale string) []model.PostTranslation {
	for _, t := range ts {
		if t.Locale == locale {
			return []model.PostTranslation{t}
		}
	}
	if len(ts) > 0 {
		return []model.PostTranslation{ts[0]}
	}
	return ts
}
