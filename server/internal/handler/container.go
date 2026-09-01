// Package handler HTTP 控制器层（public + admin）
package handler

import (
	"strconv"

	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"
	"leslie-blog-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Container 控制器容器
type Container struct {
	Auth      *AuthHandler
	Post      *PostHandler
	Category  *CategoryHandler
	Project   *ProjectHandler
	TechStack *TechStackHandler
	Media     *MediaHandler
	Comment   *CommentHandler
	Settings  *SettingsHandler
	Search    *SearchHandler
	Stats     *StatsHandler
	GitHub    *GitHubHandler
}

// New 构造控制器容器
func New(svc *service.Container) *Container {
	return &Container{
		Auth:      NewAuthHandler(svc.Auth),
		Post:      NewPostHandler(svc.Post),
		Category:  NewCategoryHandler(svc.Category),
		Project:   NewProjectHandler(svc.Project),
		TechStack: NewTechStackHandler(svc.TechStack),
		Media:     NewMediaHandler(svc.Media),
		Comment:   NewCommentHandler(svc.Comment),
		Settings:  NewSettingsHandler(svc.Settings),
		Search:    NewSearchHandler(svc.Search),
		Stats:     NewStatsHandler(svc.Stats),
		GitHub:    NewGitHubHandler(svc.GitHub),
	}
}

// ===== 通用 helper =====

// paramID 从 :id 路径参数解析 uint
func paramID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

// currentUserID 从 context 取 jwt claims 中的 uid
func currentUserID(c *gin.Context) uint {
	if v, ok := c.Get("userID"); ok {
		if uid, ok := v.(uint); ok {
			return uid
		}
	}
	return 0
}

// locale 取 locale 参数
func locale(c *gin.Context) string {
	return utils.Locale(c)
}

// parseInt8 解析 status 等为 *int8
func parseInt8(s string) *int8 {
	if s == "" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	i := int8(n)
	return &i
}
