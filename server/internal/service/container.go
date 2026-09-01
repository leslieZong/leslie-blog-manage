// Package service 业务逻辑层
package service

import (
	"leslie-blog-server/config"
	"leslie-blog-server/internal/repository"
	"leslie-blog-server/pkg/utils"
)

// Container 服务容器，集中持有各域 service
type Container struct {
	Auth      *AuthService
	Post      *PostService
	Category  *CategoryService
	Project   *ProjectService
	TechStack *TechStackService
	Media     *MediaService
	Comment   *CommentService
	Settings  *SettingsService
	Search    *SearchService
	Stats     *StatsService
	GitHub    *GitHubService
}

// New 构造服务容器
func New(cfg *config.Config, repos *repository.Container, jwt *utils.JWT) *Container {
	return &Container{
		Auth:      NewAuthService(repos.User, jwt),
		Post:      NewPostService(repos.Post, repos.Category),
		Category:  NewCategoryService(repos.Category),
		Project:   NewProjectService(repos.Project),
		TechStack: NewTechStackService(repos.TechStack),
		Media:     NewMediaService(repos.Media, cfg),
		Comment:   NewCommentService(repos.Comment),
		Settings:  NewSettingsService(repos.Settings),
		Search:    NewSearchService(repos.Post, repos.Project),
		Stats:     NewStatsService(repos),
		GitHub:    NewGitHubService(cfg.GitHub),
	}
}
