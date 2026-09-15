package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DefaultPage 默认页码。
const DefaultPage = 1

// DefaultPageSize 默认每页数量。
const DefaultPageSize = 20

// MaxPageSize 最大每页数量。
//
// 防止前端传：
//
//	pageSize=1000000
//
// 导致数据库一次查询大量数据。
const MaxPageSize = 100

// Params 表示分页参数。
type Params struct {
	Page     int
	PageSize int
}

// Parse 从 URL Query 中解析分页参数。
//
// 例如：
//
//	?page=2&pageSize=20
//
// 转换成：
//
//	Page = 2
//	PageSize = 20
func Parse(c *gin.Context) Params {

	page := DefaultPage
	pageSize := DefaultPageSize

	// 获取 page。
	//
	// 如果 URL 中不存在 page：
	//
	// c.DefaultQuery(...)
	//
	// 会使用默认值。
	pageString := c.DefaultQuery(
		"page",
		strconv.Itoa(DefaultPage),
	)

	// 将 string 转成 int。
	if value, err := strconv.Atoi(pageString); err == nil {

		if value > 0 {
			page = value
		}
	}

	// 获取 pageSize。
	pageSizeString := c.DefaultQuery(
		"pageSize",
		strconv.Itoa(DefaultPageSize),
	)

	if value, err := strconv.Atoi(pageSizeString); err == nil {

		if value > 0 {
			pageSize = value
		}
	}

	// 限制最大 pageSize。
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset 计算 SQL OFFSET。
//
// 第 1 页：
//
//	(1 - 1) * 20 = 0
//
// 第 2 页：
//
//	(2 - 1) * 20 = 20
//
// 第 3 页：
//
//	(3 - 1) * 20 = 40
func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// TotalPages 根据总数据量计算总页数。
func (p Params) TotalPages(total int64) int {

	if total <= 0 {
		return 0
	}

	return int(math.Ceil(
		float64(total) / float64(p.PageSize),
	))
}

// Result 是统一的分页返回结构。
//
// 例如：
//
//	{
//	    "list": [...],
//	    "total": 100,
//	    "page": 1,
//	    "pageSize": 20
//	}
type Result[T any] struct {
	List       []T   `json:"list"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

// NewResult 创建分页结果。
func NewResult[T any](
	list []T,
	params Params,
	total int64,
) Result[T] {

	return Result[T]{
		List:       list,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: params.TotalPages(total),
	}
}
