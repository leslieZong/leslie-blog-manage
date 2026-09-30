package service

import (
	"context"
	"leslie-blog-server/internal/modules/audit/dto"
	"leslie-blog-server/internal/modules/audit/model"
	"leslie-blog-server/internal/modules/audit/repository"
	"leslie-blog-server/internal/pkg/ulid"
)

type AuditService interface {
	Record(
		ctx context.Context,
		log dto.CreateAuditLogRequest,
	) error
	List(
		ctx context.Context,
		query repository.AuditLogListQuery,
	) ([]*model.AuditLog, int64, error)
}

type auditService struct {
	repo repository.AuditLogRepository
}

func NewAuditService(
	repo repository.AuditLogRepository,
) AuditService {

	return &auditService{
		repo: repo,
	}
}

func (s *auditService) Record(
	ctx context.Context,
	log dto.CreateAuditLogRequest,
) error {
	auditLog := &model.AuditLog{
		ID:           ulid.New(),
		UserID:       log.UserID,
		Action:       log.Action,
		Resource:     log.Resource,
		ResourceID:   log.ResourceID,
		RequestID:    &log.RequestID,
		IP:           &log.IP,
		UserAgent:    &log.UserAgent,
		Result:       log.Result,
		ErrorMessage: log.ErrorMessage,
	}

	return s.repo.Create(
		ctx,
		auditLog,
	)
}

func (s *auditService) List(
	ctx context.Context,
	query repository.AuditLogListQuery,
) ([]*model.AuditLog, int64, error) {
	list, total, err := s.repo.FindPage(
		ctx,
		query,
	)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
