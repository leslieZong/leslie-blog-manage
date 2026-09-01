package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

// PageParams 从 query 解析分页参数
type PageParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   string // 前端可能传字符串或数字
}

// ParsePage 从 gin.Context 解析分页参数
func ParsePage(c *gin.Context) PageParams {
	p := PageParams{
		Page:     atoi(c.Query("page"), defaultPage),
		PageSize: atoi(c.Query("pageSize"), defaultPageSize),
		Keyword:  c.Query("keyword"),
		Status:   c.Query("status"),
	}
	if p.Page < 1 {
		p.Page = defaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = defaultPageSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
	return p
}

// Offset 计算偏移量
func (p PageParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func atoi(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// Locale 从 query 获取语言，默认 zh-CN
func Locale(c *gin.Context) string {
	loc := c.Query("locale")
	if loc == "" {
		loc = c.GetHeader("Accept-Language")
		if len(loc) >= 5 {
			loc = loc[:5]
		}
	}
	switch loc {
	case "zh-CN", "zh", "zh-cn":
		return "zh-CN"
	case "en-US", "en", "en-us":
		return "en-US"
	default:
		return "zh-CN"
	}
}
