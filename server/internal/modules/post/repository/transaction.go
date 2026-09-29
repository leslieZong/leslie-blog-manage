package repository

import (
	"gorm.io/gorm"
)

// Repositories 保存 Post 模块事务中需要使用的 Repository。
//
// 所有 Repository 都使用同一个 tx。
type Repositories struct {
	Post    PostRepository
	PostTag PostTagRepository
}

// NewRepositories 根据数据库连接创建 Repository 集合。
//
// db 可以是普通 *gorm.DB，
// 也可以是事务中的 tx。
//
// 因此这个方法非常重要：
//
// 普通操作：
//
//	NewRepositories(db)
//
// 事务操作：
//
//	NewRepositories(tx)
//
// 两者 Repository 的代码完全一样。
func NewRepositories(
	db *gorm.DB,
) *Repositories {

	return &Repositories{
		Post: NewPostRepository(db),

		PostTag: NewPostTagRepository(db),
	}
}
