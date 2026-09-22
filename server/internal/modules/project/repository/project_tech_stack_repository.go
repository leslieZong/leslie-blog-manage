package repository

import (
	"context"

	"leslie-blog-server/internal/modules/project/model"
)

// ProjectTechStackRepository
//
// 负责 Project 与 TechStack 之间的关系数据。
//
// 例如：
//
// Project A
//
//	↓
//
// ProjectTechStack
//
//	↓
//
// Vue
// Go
// MySQL
type ProjectTechStackRepository interface {

	// ReplaceTechStacks
	//
	// 将一个 Project 的所有 TechStack 关系，
	// 替换成新的关系集合。
	//
	// 例如：
	//
	// 原来：
	// Vue
	// Go
	// MySQL
	//
	// 新的：
	// Vue
	// TypeScript
	// Go
	//
	// 执行：
	//
	// 删除旧关系
	// +
	// 创建新关系
	ReplaceTechStacks(
		ctx context.Context,
		projectID string,
		techStackIDs []string,
	) error

	// FindByProjectID 查询一个 Project 使用的所有 TechStack。
	FindByProjectID(
		ctx context.Context,
		projectID string,
	) ([]*model.ProjectTechStack, error)
}
