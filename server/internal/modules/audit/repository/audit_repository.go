package repository

import (
	"context"
	"leslie-blog-server/internal/modules/audit/model"
	"leslie-blog-server/internal/pkg/pagination"
)

type AuditLogListQuery struct {
	// 分页参数。
	pagination.Params

	UserID     string
	Action     string
	Resource   string
	ResourceID string
	Result     string
}

type AuditLogRepository interface {

	// Create 创建一条审计日志。
	Create(
		ctx context.Context,
		log *model.AuditLog,
	) error

	// FindPage 分页查询审计日志。
	FindPage(
		ctx context.Context,
		query AuditLogListQuery,
	) ([]*model.AuditLog, int64, error)
}
