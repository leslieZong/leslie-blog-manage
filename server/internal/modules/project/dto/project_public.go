package dto

import "leslie-blog-server/internal/modules/project/model"

type PublicProjectResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Description *string `json:"description"`

	Cover *string `json:"cover"`

	GitHubURL *string `json:"githubUrl"`

	DemoURL *string `json:"demoUrl"`

	Featured bool `json:"featured"`
}

func FromPublicModel(
	project *model.Project,
) *PublicProjectResponse {

	if project == nil {
		return nil
	}

	return &PublicProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Slug:        project.Slug,
		Description: project.Description,
		Cover:       project.Cover,
		GitHubURL:   project.GitHubURL,
		DemoURL:     project.DemoURL,
		Featured:    project.Featured,
	}
}

func FromPublicModels(
	projects []*model.Project,
) []*PublicProjectResponse {

	res := make([]*PublicProjectResponse, 0, len(projects))

	for _, project := range projects {
		res = append(res, FromPublicModel(project))
	}

	return res
}
