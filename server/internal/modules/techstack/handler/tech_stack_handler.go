package handler

import (
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/techstack/dto"
	"leslie-blog-server/internal/modules/techstack/service"
	"leslie-blog-server/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TechStackHandler struct {
	service service.TechStackService
}

func NewTechStackHandler(
	techStackService service.TechStackService,
) *TechStackHandler {

	return &TechStackHandler{
		service: techStackService,
	}
}

func (h *TechStackHandler) Create(
	c *gin.Context,
) {

	var req dto.CreateTechStackRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	techStack, err := h.service.Create(
		c.Request.Context(),
		req,
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
		dto.FromModel(techStack),
	)
}

func (h *TechStackHandler) GetByID(
	c *gin.Context,
) {

	id := c.Param("id")

	techStack, err := h.service.GetByID(
		c.Request.Context(),
		id,
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
		dto.FromModel(techStack),
	)
}

func (h *TechStackHandler) List(
	c *gin.Context,
) {

	techStacks, err := h.service.List(
		c.Request.Context(),
	)

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}
	items := dto.FromModelList(techStacks)

	response.Success(
		c,
		items,
	)
}

func (h *TechStackHandler) Update(
	c *gin.Context,
) {

	id := c.Param("id")

	var req dto.UpdateTechStackRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	techStack, err := h.service.Update(
		c.Request.Context(),
		id,
		req,
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
		dto.FromModel(techStack),
	)
}

func (h *TechStackHandler) Delete(
	c *gin.Context,
) {

	id := c.Param("id")

	if err := h.service.Delete(
		c.Request.Context(),
		id,
	); err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	response.Success(
		c,
		gin.H{
			"status": "ok",
		},
	)
}
