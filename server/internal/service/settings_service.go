package service

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
)

type SettingsService struct{ repo *repository.SettingsRepo }

func NewSettingsService(repo *repository.SettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

// Get 返回站点设置（公开版隐藏社交敏感字段可选，这里直接返回）
func (s *SettingsService) Get() (*model.Settings, error) {
	return s.repo.Get()
}

// Update 更新设置
func (s *SettingsService) Update(req model.SettingsReq) error {
	cur, err := s.repo.Get()
	if err != nil {
		return err
	}
	cur.SiteName = req.SiteName
	cur.Logo = req.Logo
	cur.Description = req.Description
	cur.Keywords = req.Keywords
	cur.Author = req.Author
	cur.ICP = req.ICP
	if req.Social != nil {
		cur.SocialGithub = req.Social.Github
		cur.SocialEmail = req.Social.Email
		cur.SocialTwitter = req.Social.Twitter
	}
	return s.repo.Update(cur)
}
