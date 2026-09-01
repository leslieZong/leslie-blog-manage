package handler

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"
	"leslie-blog-server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct{ svc *service.PostService }

func NewPostHandler(svc *service.PostService) *PostHandler { return &PostHandler{svc: svc} }

// ===== Public =====

// List GET /api/v1/posts
func (h *PostHandler) List(c *gin.Context) {
	p := utils.ParsePage(c)
	f := repository.PostFilter{
		Keyword:    p.Keyword,
		CategoryID: parseUint(c.Query("categoryId")),
		Status:     statusPtr(p.Status),
	}
	loc := locale(c)
	// 公开侧强制只看已发布
	published := int8(1)
	f.Status = &published
	list, total, err := h.svc.ListPublished(f, loc, p.Page, p.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OKPage(c, list, total, p.Page, p.PageSize)
}

// GetBySlug GET /api/v1/posts/:slug
func (h *PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "缺少 slug")
		return
	}
	dto, err := h.svc.GetBySlug(slug, locale(c))
	if err != nil {
		response.Fail(c, response.CodeNotFound, err.Error())
		return
	}
	response.OK(c, dto)
}

// ===== Admin =====

// AdminList GET /api/v1/admin/posts
func (h *PostHandler) AdminList(c *gin.Context) {
	p := utils.ParsePage(c)
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		statusPtr = parseInt8(s)
	}
	f := repository.PostFilter{
		Keyword:    p.Keyword,
		CategoryID: parseUint(c.Query("categoryId")),
		Status:     statusPtr,
	}
	list, total, err := h.svc.ListAdmin(f, locale(c), p.Page, p.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OKPage(c, list, total, p.Page, p.PageSize)
}

// AdminGet GET /api/v1/admin/posts/:id
func (h *PostHandler) AdminGet(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	dto, err := h.svc.GetAdmin(id, locale(c))
	if err != nil {
		response.Fail(c, response.CodeNotFound, err.Error())
		return
	}
	response.OK(c, dto)
}

// AdminCreate POST /api/v1/admin/posts
func (h *PostHandler) AdminCreate(c *gin.Context) {
	var req model.PostReq
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

// AdminUpdate PUT /api/v1/admin/posts/:id
func (h *PostHandler) AdminUpdate(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	var req model.PostReq
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

// AdminDelete DELETE /api/v1/admin/posts/:id
func (h *PostHandler) AdminDelete(c *gin.Context) {
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

// Publish PATCH /api/v1/admin/posts/:id/publish
func (h *PostHandler) Publish(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	if err := h.svc.Publish(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已发布")
}

// Unpublish PATCH /api/v1/admin/posts/:id/unpublish
func (h *PostHandler) Unpublish(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	if err := h.svc.Unpublish(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已下线")
}

// ToggleTop PATCH /api/v1/admin/posts/:id/top
func (h *PostHandler) ToggleTop(c *gin.Context) {
	id, ok := paramID(c)
	if !ok {
		return
	}
	var body struct {
		IsTop bool `json:"isTop"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.svc.SetTop(id, body.IsTop); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "已更新")
}

// ===== helper =====

func parseUint(s string) uint {
	if s == "" {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

func statusPtr(s string) *int8 { return parseInt8(s) }
