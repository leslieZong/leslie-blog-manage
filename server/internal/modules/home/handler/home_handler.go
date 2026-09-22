package handler

import (
	"leslie-blog-server/internal/modules/home/service"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	service service.HomeService
}

func NewHomeHandler(
	service service.HomeService,
) *HomeHandler {

	return &HomeHandler{
		service: service,
	}
}

func (h *HomeHandler) GetHome(
	c *gin.Context,
) {

	result, err :=
		h.service.GetHome(
			c.Request.Context(),
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
		result,
	)
}
