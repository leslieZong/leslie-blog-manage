package handler

import (
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/audit/dto"
	"leslie-blog-server/internal/modules/audit/repository"
	"leslie-blog-server/internal/modules/audit/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service service.AuditService
}

func NewAuditHandler(
	service service.AuditService,
) *AuditHandler {
	return &AuditHandler{
		service: service,
	}
}

func (h *AuditHandler) List(c *gin.Context) error {
	query := repository.AuditLogListQuery{
		Params:     pagination.Parse(c),
		UserID:     c.Query("user_id"),
		Action:     c.Query("action"),
		Resource:   c.Query("resource"),
		ResourceID: c.Query("resource_id"),
		Result:     c.Query("result"),
	}

	auditLogs, total, err := h.service.List(
		c.Request.Context(),
		query,
	)

	if err != nil {
		return err
	}
	// 第三步：
	// 转换为 API 列表 DTO。
	res := dto.FromModelList(auditLogs)
	// 第四步：
	// 构造统一分页结果。
	result := pagination.NewResult(
		res,
		query.Params,
		total,
	)
	response.Success(c, result)
	return nil
}
func (h *AuditHandler) Record(c *gin.Context) error {
	var req dto.CreateAuditLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"invalid parameters",
			err,
		)
	}

	if err := h.service.Record(
		c.Request.Context(),
		req,
	); err != nil {
		return err
	}

	return nil
}
