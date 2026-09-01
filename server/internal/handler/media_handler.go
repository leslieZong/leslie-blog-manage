package handler

import (
	"io"

	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"
	"leslie-blog-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct{ svc *service.MediaService }

func NewMediaHandler(svc *service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

// List GET /api/v1/admin/media
func (h *MediaHandler) List(c *gin.Context) {
	p := utils.ParsePage(c)
	list, total, err := h.svc.List(p.Keyword, c.Query("type"), p.Page, p.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OKPage(c, list, total, p.Page, p.PageSize)
}

// Upload POST /api/v1/admin/media/upload  (multipart form, field "file")
func (h *MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "未收到文件: "+err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.Save(file.Filename, content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, m)
}

// Delete DELETE /api/v1/admin/media/:id
func (h *MediaHandler) Delete(c *gin.Context) {
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
