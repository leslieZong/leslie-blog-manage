package handler

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/techstack/dto"
	"leslie-blog-server/internal/response"
)

func (h *TechStackHandler) PublicList(
	c *gin.Context,
) {

	techStacks, err := h.service.ListPublic(
		c.Request.Context(),
	)

	if err != nil {
		response.AppError(
			c,
			err,
		)

		return
	}

	items := dto.FromPublicModelList(techStacks)
	response.Success(
		c,
		items,
	)
}

func (h *TechStackHandler) PublicGetBySlug(
	c *gin.Context,
) {

	slug := c.Param("slug")

	techStack, err := h.service.GetPublicBySlug(
		c.Request.Context(),
		slug,
	)

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	response.Success(
		c,
		dto.FromPublicModel(techStack),
	)
}
