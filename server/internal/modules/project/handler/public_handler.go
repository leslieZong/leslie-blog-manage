package handler

import (
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/project/dto"
	"leslie-blog-server/internal/modules/project/repository"
	"leslie-blog-server/internal/modules/project/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectPublicHandler struct {
	service service.ProjectService
}

func NewProjectPublicHandler(
	projectService service.ProjectService,
) *ProjectPublicHandler {

	return &ProjectPublicHandler{
		service: projectService,
	}
}

func (h *ProjectPublicHandler) List(
	c *gin.Context,
) {

	query := repository.ProjectListQuery{
		Params:  pagination.Parse(c),
		Keyword: c.Query("keyword"),
	}

	// ---------------------------------------------
	// featured
	// ---------------------------------------------

	if value := c.Query("featured"); value != "" {

		featured, err := strconv.ParseBool(
			value,
		)

		if err != nil {

			response.Error(
				c,
				http.StatusBadRequest,
				appErrors.ErrInvalidParams,
				"invalid featured",
			)

			return
		}

		query.Featured = &featured
	}

	// 注意：
	//
	// 这里故意没有解析 status。
	//
	// Public API 不允许用户指定 status。

	res, total, err := h.service.ListPublic(
		c.Request.Context(),
		query,
	)

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}
	list := dto.FromPublicModels(
		res,
	)
	result := pagination.NewResult(
		list,
		query.Params,
		total,
	)

	response.Success(
		c,
		result,
	)
}

func (h *ProjectPublicHandler) GetByID(
	c *gin.Context,
) {

	id := c.Param("id")

	project, err := h.service.GetPublicByID(
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
		dto.FromPublicModel(project),
	)
}

func (h *ProjectPublicHandler) GetBySlug(
	c *gin.Context,
) {

	slug := c.Param("slug")

	project, err := h.service.GetPublicBySlug(
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
		dto.FromPublicModel(project),
	)
}
