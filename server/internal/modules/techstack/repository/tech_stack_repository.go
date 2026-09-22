package repository

import (
	"context"

	"leslie-blog-server/internal/modules/techstack/model"
)

// TechStackRepository 定义 TechStack 数据访问能力。
//
// Repository 的职责只有一个：
//
// “如何从数据库读取 / 写入 TechStack”
//
// 它不应该关心：
//
// - HTTP
// - Gin
// - JWT
// - Casbin
// - 前端 JSON
// - HTTP Status Code
type TechStackRepository interface {

	// FindByID 根据 ID 查询 TechStack。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.TechStack, error)

	// FindBySlug 根据 slug 查询 TechStack。
	FindBySlug(
		ctx context.Context,
		slug string,
	) (*model.TechStack, error)

	// FindByName 根据名称查询 TechStack。
	FindByName(
		ctx context.Context,
		name string,
	) (*model.TechStack, error)

	// FindAll 查询所有 TechStack。
	FindAll(
		ctx context.Context,
	) ([]*model.TechStack, error)

	// FindPublic 查询所有公开的 TechStack。
	FindPublic(
		ctx context.Context,
	) ([]*model.TechStack, error)

	// Create 创建 TechStack。
	Create(
		ctx context.Context,
		techStack *model.TechStack,
	) error

	// Update 更新 TechStack。
	Update(
		ctx context.Context,
		techStack *model.TechStack,
	) error

	// Delete 软删除 TechStack。
	Delete(
		ctx context.Context,
		id string,
	) error

	// CountProjects 统计有多少 Project 正在使用这个 TechStack。
	CountProjects(
		ctx context.Context,
		techStackID string,
	) (int64, error)

	FindByIDs(
		ctx context.Context,
		ids []string,
	) ([]*model.TechStack, error)
}
