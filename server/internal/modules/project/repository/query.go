package repository

import "leslie-blog-server/internal/pkg/pagination"

// ProjectListQuery 表示 Project 列表查询条件。
type ProjectListQuery struct {

	// 分页。
	pagination.Params

	// 搜索项目名称 / slug。
	Keyword string

	// 是否精选。
	//
	// nil  = 不筛选
	// true = 精选
	// false = 非精选
	Featured *bool

	// 状态。
	//
	// nil = 不筛选
	// 0   = 禁用
	// 1   = 启用
	Status *int8
}
