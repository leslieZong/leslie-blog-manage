package repository

import (
	"context"
	"leslie-blog-server/internal/modules/audit/model"

	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{
		db: db,
	}
}

func (r *auditLogRepository) Create(
	ctx context.Context,
	log *model.AuditLog,
) error {

	return r.db.
		WithContext(ctx).
		Create(log).
		Error
}
func (r *auditLogRepository) FindPage(
	ctx context.Context,
	params AuditLogListQuery,
) ([]*model.AuditLog, int64, error) {
	var (
		auditLogs []*model.AuditLog
		total     int64
	)
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})
	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}
	if params.Resource != "" {
		query = query.Where("resource = ?", params.Resource)
	}
	if params.ResourceID != "" {
		query = query.Where("resource_id = ?", params.ResourceID)
	}
	if params.Result != "" {
		query = query.Where("result = ?", params.Result)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC").Offset(params.Offset())
	err = query.Find(&auditLogs).Error
	if err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}
