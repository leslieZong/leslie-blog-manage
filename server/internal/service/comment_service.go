package service

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
)

type CommentService struct{ repo *repository.CommentRepo }

func NewCommentService(repo *repository.CommentRepo) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) List(keyword string, status *int8, page, pageSize int) ([]model.Comment, int64, error) {
	return s.repo.List(keyword, status, page, pageSize)
}

func (s *CommentService) ListByPost(postID uint) ([]model.Comment, error) {
	return s.repo.ListByPost(postID)
}

func (s *CommentService) Create(req model.CommentReq, ip string) error {
	c := model.Comment{
		PostID:   req.PostID,
		Author:   req.Author,
		Email:    req.Email,
		Content:  req.Content,
		ParentID: req.Parent,
		Status:   0, // 待审
		IP:       ip,
	}
	return s.repo.Create(&c)
}

func (s *CommentService) Approve(id uint) error { return s.repo.UpdateStatus(id, 1) }
func (s *CommentService) Reject(id uint) error  { return s.repo.UpdateStatus(id, 2) }
func (s *CommentService) Delete(id uint) error { return s.repo.Delete(id) }
