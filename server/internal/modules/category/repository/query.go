package repository

import "leslie-blog-server/internal/pkg/pagination"

// CategoryListQuery 是 Category 列表查询条件。
type CategoryListQuery struct {

	// 分页参数。
	pagination.Params

	// 搜索关键字。
	Keyword string

	// 状态。
	//
	// 1 = enabled
	// 0 = disabled
	Status *int8
}
