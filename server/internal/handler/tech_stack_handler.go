package handler

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type TechStackHandler struct{ svc *service.TechStackService }

func NewTechStackHandler(svc *service.TechStackService) *TechStackHandler {
	return &TechStackHandler{svc: svc}
}

// List GET /api/v1/tech-stack
func (h *TechStackHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Query("keyword"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}

// Create POST /api/v1/admin/tech-stack
func (h *TechStackHandler) Create(c *gin.Context) {
	var req model.TechStackReq
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

// Update PUT /api/v1/admin/tech-stack/:id
func (h *TechStackHandler) Update(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	var req model.TechStackReq
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

// Delete DELETE /api/v1/admin/tech-stack/:id
func (h *TechStackHandler) Delete(c *gin.Context) {
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
