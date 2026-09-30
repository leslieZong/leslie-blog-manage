package service

import (
	"context"
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/audit/dto"
	"leslie-blog-server/internal/modules/audit/model"
	"leslie-blog-server/internal/modules/audit/repository"
	"leslie-blog-server/internal/pkg/requestmeta"
	"leslie-blog-server/internal/pkg/ulid"
	"net/http"
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
	metadata, ok := requestmeta.FromContext(ctx)
	if !ok {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"request metadata not found",
		)
	}
	auditLog := &model.AuditLog{
		ID:           ulid.New(),
		UserID:       &metadata.UserID,
		Action:       log.Action,
		Resource:     log.Resource,
		ResourceID:   log.ResourceID,
		RequestID:    &metadata.RequestID,
		IP:           &metadata.IP,
		UserAgent:    &metadata.UserAgent,
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
