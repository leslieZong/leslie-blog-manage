package repository

import (
	"context"
)

// PostTagRepository 负责文章与标签关联数据。
type PostTagRepository interface {

	// ReplaceTags 替换文章的全部 Tag。
	//
	// 例如：
	//
	// 原来：
	//
	// Post A
	// ├── Vue3
	// └── Vite
	//
	// 修改后：
	//
	// Post A
	// ├── Vue3
	// └── TypeScript
	//
	// Repository 会删除旧关系，
	// 然后创建新关系。
	ReplaceTags(
		ctx context.Context,
		postID string,
		tagIDs []string,
	) error

	// DeleteByPostID 删除文章的所有 Tag 关系。
	DeleteByPostID(
		ctx context.Context,
		postID string,
	) error
}
