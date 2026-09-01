package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type CommentRepo struct{ db *gorm.DB }

func NewCommentRepo(db *gorm.DB) *CommentRepo { return &CommentRepo{db: db} }

func (r *CommentRepo) List(keyword string, status *int8, page, pageSize int) ([]model.Comment, int64, error) {
	var total int64
	q := r.db.Model(&model.Comment{})
	if keyword != "" {
		q = q.Where("content LIKE ? OR author LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Comment
	if err := paginate(q.Order("id DESC"), page, pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	// 填充文章标题
	for i := range list {
		var title string
		r.db.Model(&model.PostTranslation{}).
			Where("post_id = ? AND locale = ?", list[i].PostID, "zh-CN").
			Limit(1).Pluck("title", &title)
		list[i].PostTitle = title
	}
	return list, total, nil
}

// ListByPost 取某文章已通过评论（公开侧）
func (r *CommentRepo) ListByPost(postID uint) ([]model.Comment, error) {
	var list []model.Comment
	if err := r.db.Where("post_id = ? AND status = ?", postID, int8(1)).
		Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CommentRepo) Create(c *model.Comment) error { return r.db.Create(c).Error }

func (r *CommentRepo) UpdateStatus(id uint, status int8) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CommentRepo) Delete(id uint) error { return r.db.Delete(&model.Comment{}, id).Error }
