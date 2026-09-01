package handler

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"
	"leslie-blog-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct{ svc *service.ProjectService }

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// List GET /api/v1/projects  |  GET /api/v1/admin/projects
func (h *ProjectHandler) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := h.svc.List(p.Keyword, parseInt8(p.Status), p.Page, p.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OKPage(c, list, total, p.Page, p.PageSize)
}

// ListAll GET /api/v1/projects/all（不分页）
func (h *ProjectHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}

// Get GET /api/v1/admin/projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	p, err := h.svc.Get(id)
	if err != nil {
		response.Fail(c, response.CodeNotFound, err.Error())
		return
	}
	response.OK(c, p)
}

// Create POST /api/v1/admin/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	var req model.ProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Create(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已创建")
}

// Update PUT /api/v1/admin/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	var req model.ProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Update(id, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已更新")
}

// Delete DELETE /api/v1/admin/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已删除")
}
