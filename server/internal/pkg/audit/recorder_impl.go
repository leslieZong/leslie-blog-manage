package audit

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/audit/dto"
	"leslie-blog-server/internal/modules/audit/service"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/logger"
	"leslie-blog-server/internal/pkg/requestmeta"
)

// recorder 是 AuditRecorder 的具体实现。
type recorder struct {

	// auditService 负责真正保存审计日志。
	auditService service.AuditService

	// logger 用于记录 Audit 本身发生的异常。
	logger *logger.Logger
}

// NewRecorder 创建 Audit Recorder。
func NewRecorder(
	auditService service.AuditService,
	logger *logger.Logger,
) Recorder {

	return &recorder{
		auditService: auditService,
		logger:       logger,
	}
}
func (r *recorder) RecordSuccess(
	ctx context.Context,
	event Event,
) error {

	return r.record(
		ctx,
		event,
		ResultSuccess,
		nil,
	)
}
func (r *recorder) RecordFailure(
	ctx context.Context,
	event Event,
	err error,
) error {

	if err == nil {
		err = errors.New("unknown business error")
	}

	return r.record(
		ctx,
		event,
		ResultFailed,
		err,
	)
}
func (r *recorder) record(
	ctx context.Context,
	event Event,
	result string,
	businessErr error,
) error {
	metadata, ok := requestmeta.FromContext(ctx)
	if !ok {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"request metadata not found",
		)
	}

	UserID := auth.GetContextUserID(ctx)
	RequestID := logger.RequestID(ctx)
	// ----------------------------------------
	// 构造 Audit DTO
	// ----------------------------------------

	req := dto.CreateAuditLogRequest{
		UserID:     &UserID,
		Action:     event.Action,
		Resource:   event.Resource,
		ResourceID: &event.ResourceID,
		RequestID:  RequestID,
		IP:         metadata.IP,
		UserAgent:  metadata.UserAgent,
		Result:     result,
	}

	// ----------------------------------------
	// 失败原因
	// ----------------------------------------

	if businessErr != nil {
		message := businessErr.Error()

		req.ErrorMessage = &message
	}

	// ----------------------------------------
	// 保存 Audit Log
	// ----------------------------------------

	if err := r.auditService.Record(
		ctx,
		req,
	); err != nil {

		// ----------------------------------------
		// Audit 自己失败
		// ----------------------------------------
		//
		// 注意：
		//
		// Audit 是辅助能力。
		//
		// 不应该因为 Audit 写失败，
		// 把已经成功的业务操作重新变成失败。
		// ----------------------------------------

		logger.ErrorContext(
			r.logger,
			ctx,
			"failed to record audit log",
			slog.String(
				"action",
				event.Action,
			),
			slog.String(
				"resource",
				event.Resource,
			),
			slog.String(
				"resource_id",
				event.ResourceID,
			),
			slog.Any(
				"error",
				err,
			),
		)

		return err
	}

	return nil
}
