package handler

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{ svc *service.CategoryService }

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Query("keyword"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}

// Get GET /api/v1/categories/:id
func (h *CategoryHandler) Get(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	cat, err := h.svc.Get(id)
	if err != nil {
		response.Fail(c, response.CodeNotFound, err.Error())
		return
	}
	response.OK(c, cat)
}

// Create POST /api/v1/admin/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CategoryReq
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

// Update PUT /api/v1/admin/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	var req model.CategoryReq
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

// Delete DELETE /api/v1/admin/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
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
