package audit

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/audit/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/httpx"
	"leslie-blog-server/internal/pkg/permission"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	auditHandler *handler.AuditHandler,
	enforcer *casbin.Enforcer,
) {
	auditLog := router.Group("/audit-logs")
	{
		auditLog.GET(
			"",
			middleware.Permission(enforcer, permission.AuditRead),
			httpx.Adapt(auditHandler.List),
		)
		auditLog.POST(
			"/record",
			middleware.Permission(enforcer, permission.AuditCreate),
			httpx.Adapt(auditHandler.Record),
		)
	}
}
