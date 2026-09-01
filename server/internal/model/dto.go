package model

import "time"

// ===== Auth =====

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
}

type LoginResp struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

// ===== Post =====

// PostReq 创建/更新文章请求（content 来自指定 locale 的翻译）
type PostReq struct {
	Title          string   `json:"title" binding:"required"`
	Slug           string   `json:"slug"`
	Summary        string   `json:"summary"`
	Content        string   `json:"content" binding:"required"`
	Cover          string   `json:"cover"`
	CategoryID     uint     `json:"categoryId"`
	Tags           []string `json:"tags"`
	Status         int8     `json:"status"`
	IsTop          bool     `json:"isTop"`
	Locale         string   `json:"locale,omitempty"` // 默认 zh-CN
	SEOTitle       string   `json:"seoTitle,omitempty"`
	SEODescription string   `json:"seoDescription,omitempty"`
	SEOKeywords    string   `json:"seoKeywords,omitempty"`
}

// PostDTO 文章响应（翻译扁平化，tags 转 []string）
type PostDTO struct {
	ID             uint       `json:"id"`
	Slug           string     `json:"slug"`
	Title          string     `json:"title"`
	Summary        string     `json:"summary"`
	Content        string     `json:"content,omitempty"`
	Cover          string     `json:"cover"`
	CategoryID     uint       `json:"categoryId"`
	CategoryName   string     `json:"categoryName,omitempty"`
	Tags           []string   `json:"tags,omitempty"`
	Status         int8       `json:"status"`
	IsTop          bool       `json:"isTop"`
	ViewCount      uint       `json:"viewCount"`
	CommentCount   uint       `json:"commentCount"`
	SEOTitle       string     `json:"seoTitle,omitempty"`
	SEODescription string     `json:"seoDescription,omitempty"`
	SEOKeywords    string     `json:"seoKeywords,omitempty"`
	PublishedAt    *time.Time `json:"publishedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// ===== Category =====

type CategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    uint   `json:"parentId"`
	SortOrder   int    `json:"sortOrder"`
}

// ===== Project =====

type ProjectReq struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Cover       string   `json:"cover"`
	DemoURL     string   `json:"demoUrl"`
	RepoURL     string   `json:"repoUrl"`
	TechStack   []string `json:"techStack"`
	Status      int8     `json:"status"`
	SortOrder   int      `json:"sortOrder"`
}

// ===== TechStack =====

type TechStackReq struct {
	Name        string `json:"name" binding:"required"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Level       int    `json:"level"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}

// ===== Comment =====

type CommentReq struct {
	PostID  uint   `json:"postId" binding:"required"`
	Author  string `json:"author" binding:"required"`
	Email   string `json:"email"`
	Content string `json:"content" binding:"required"`
	Parent  uint   `json:"parentId"`
}

// ===== Settings =====

type SettingsReq struct {
	SiteName   string         `json:"siteName"`
	Logo       string         `json:"logo"`
	Description string        `json:"description"`
	Keywords   string         `json:"keywords"`
	Author     string         `json:"author"`
	ICP        string         `json:"icp"`
	Social     *SocialLinks   `json:"social,omitempty"`
}

type SocialLinks struct {
	Github  string `json:"github"`
	Email   string `json:"email"`
	Twitter string `json:"twitter"`
}

// ===== Search =====

type SearchResult struct {
	Posts    []PostDTO `json:"posts"`
	Projects []Project `json:"projects,omitempty"`
	Total    int64     `json:"total"`
}

// ===== Stats =====

type Stats struct {
	PostCount      int64 `json:"postCount"`
	PublishedCount int64 `json:"publishedCount"`
	CategoryCount  int64 `json:"categoryCount"`
	CommentCount   int64 `json:"commentCount"`
	ProjectCount   int64 `json:"projectCount"`
	TechStackCount  int64 `json:"techStackCount"`
	MediaCount     int64 `json:"mediaCount"`
	ViewCount      int64 `json:"viewCount"`
}

type RecentItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Status    int8      `json:"status,omitempty"`
}
