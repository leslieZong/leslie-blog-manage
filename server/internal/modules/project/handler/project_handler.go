package handler

import (
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/project/dto"
	"leslie-blog-server/internal/modules/project/repository"
	"leslie-blog-server/internal/modules/project/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/pkg/utils"
	"leslie-blog-server/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(
	projectService service.ProjectService,
) *ProjectHandler {

	return &ProjectHandler{
		service: projectService,
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {

	var req dto.CreateProjectRequest

	// ---------------------------------------------
	// 1. JSON → DTO
	// ---------------------------------------------

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	// ---------------------------------------------
	// 2. 调用 Service
	// ---------------------------------------------

	project, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		response.AppError(c, err)

		return
	}

	// ---------------------------------------------
	// 3. Model → DTO
	// ---------------------------------------------

	response.Success(
		c,
		dto.FromModel(project),
	)
}

func (h *ProjectHandler) GetByID(
	c *gin.Context,
) {

	id := c.Param("id")

	project, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {

		response.AppError(c, err)

		return
	}

	response.Success(
		c,
		dto.FromModel(project),
	)
}

func (h *ProjectHandler) Update(
	c *gin.Context,
) {

	id := c.Param("id")

	var req dto.UpdateProjectRequest

	// JSON 参数解析。
	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	// Service。
	project, err := h.service.Update(
		c.Request.Context(),
		id,
		req,
	)

	if err != nil {

		response.AppError(c, err)

		return
	}

	response.Success(
		c,
		dto.FromModel(project),
	)
}

func (h *ProjectHandler) Delete(
	c *gin.Context,
) {

	id := c.Param("id")

	err := h.service.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {

		response.AppError(c, err)

		return
	}

	response.Success(
		c,
		nil,
	)
}

func (h *ProjectHandler) List(
	c *gin.Context,
) {
	query := repository.ProjectListQuery{
		Params:  pagination.Parse(c),
		Keyword: c.Query("keyword"),
	}

	// ---------------------------------------------
	// Featured
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

	// ---------------------------------------------
	// Status
	// ---------------------------------------------
	statusValue := c.Query("status")
	status, err := utils.ParseStatus(statusValue)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid status",
		)
		return
	}

	query.Status = status

	// ---------------------------------------------
	// Service
	// ---------------------------------------------
	ctx := c.Request.Context()
	res, total, err := h.service.ListPage(
		ctx,
		query,
	)

	if err != nil {
		response.AppError(c, err)

		return
	}
	// 将 Model 列表转换成 DTO 列表。
	list := dto.FromModels(res)

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
