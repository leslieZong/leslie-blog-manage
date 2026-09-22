package dto

import (
	"leslie-blog-server/internal/modules/project/model"
	"time"

	techstackdto "leslie-blog-server/internal/modules/techstack/dto"
)

type CreateProjectRequest struct {
	Name string `json:"name"`

	Slug string `json:"slug"`

	Description *string `json:"description"`

	Cover *string `json:"cover"`

	GitHubURL *string `json:"githubUrl"`

	DemoURL *string `json:"demoUrl"`

	Featured bool `json:"featured"`

	Sort int `json:"sort"`

	// TechStack IDs。
	//
	// 前端只需要提交：
	//
	// [
	//   "01xxx",
	//   "01yyy"
	// ]
	//
	// 后端负责验证这些 TechStack 是否存在。
	TechStackIDs []string `json:"techStackIds"`
}

type UpdateProjectRequest struct {
	Name string `json:"name"`

	Slug string `json:"slug"`

	Description *string `json:"description"`

	Cover *string `json:"cover"`

	GitHubURL *string `json:"githubUrl"`

	DemoURL *string `json:"demoUrl"`

	Featured bool `json:"featured"`

	Status int8 `json:"status"`

	Sort         int      `json:"sort"`
	TechStackIDs []string `json:"techStackIds"`
}

type ProjectResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Description *string `json:"description"`

	Cover *string `json:"cover"`

	GitHubURL *string `json:"githubUrl"`

	DemoURL *string `json:"demoUrl"`

	Featured bool `json:"featured"`

	Status int8 `json:"status"`

	Sort int `json:"sort"`

	CreatedAt string `json:"createdAt"`

	UpdatedAt  string                            `json:"updatedAt"`
	TechStacks []*techstackdto.TechStackResponse `json:"techStacks"`
}

func FromModel(
	project *model.Project,
) *ProjectResponse {

	if project == nil {
		return nil
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Slug:        project.Slug,
		Description: project.Description,
		Cover:       project.Cover,
		GitHubURL:   project.GitHubURL,
		DemoURL:     project.DemoURL,
		Featured:    project.Featured,
		Status:      project.Status,
		Sort:        project.Sort,
		CreatedAt:   project.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   project.UpdatedAt.Format(time.RFC3339),
		TechStacks:  techstackdto.FromModelList(project.TechStacks),
	}
}

func FromModels(
	projects []*model.Project,
) []*ProjectResponse {

	if projects == nil {
		return nil
	}

	var res []*ProjectResponse

	for _, project := range projects {
		res = append(res, FromModel(project))
	}

	return res
}
