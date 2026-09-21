package repository

import "leslie-blog-server/internal/pkg/pagination"

// TagListQuery 表示 Tag 列表查询条件。
//
// 它同时包含：
//
// 1. 分页
// 2. 关键字搜索
// 3. 状态筛选
type TagListQuery struct {

	// 分页参数。
	pagination.Params

	// Keyword 用于搜索 Tag 名称或者 slug。
	Keyword string

	// Status：
	//
	// nil = 不筛选状态
	// 0   = 禁用
	// 1   = 启用
	//
	// 使用指针是为了区分：
	//
	// status 没传
	// 和
	// status=0
	Status *int8
}
