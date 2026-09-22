package dto

import (
	"leslie-blog-server/internal/modules/techstack/model"
	"time"
)

type CreateTechStackRequest struct {
	Name string `json:"name"`

	Slug string `json:"slug"`

	Icon *string `json:"icon"`

	Description *string `json:"description"`

	OfficialURL *string `json:"officialUrl"`

	Category *string `json:"category"`

	Sort int `json:"sort"`
}

type UpdateTechStackRequest struct {
	Name string `json:"name"`

	Slug string `json:"slug"`

	Icon *string `json:"icon"`

	Description *string `json:"description"`

	OfficialURL *string `json:"officialUrl"`

	Category *string `json:"category"`

	Sort int `json:"sort"`

	Status int8 `json:"status"`
}

type TechStackResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Icon *string `json:"icon"`

	Description *string `json:"description"`

	OfficialURL *string `json:"officialUrl"`

	Category *string `json:"category"`

	Sort int `json:"sort"`

	Status int8 `json:"status"`

	CreatedAt string `json:"createdAt"`

	UpdatedAt string `json:"updatedAt"`
}

func FromModel(
	techStack *model.TechStack,
) *TechStackResponse {

	if techStack == nil {
		return nil
	}

	return &TechStackResponse{
		ID:          techStack.ID,
		Name:        techStack.Name,
		Slug:        techStack.Slug,
		Icon:        techStack.Icon,
		Description: techStack.Description,
		OfficialURL: techStack.OfficialURL,
		Category:    techStack.Category,
		Sort:        techStack.Sort,
		Status:      techStack.Status,
		CreatedAt:   techStack.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   techStack.UpdatedAt.Format(time.RFC3339),
	}
}

func FromModelList(
	techStacks []*model.TechStack,
) []*TechStackResponse {

	if techStacks == nil {
		return nil
	}

	techStackResponses := make([]*TechStackResponse, 0, len(techStacks))

	for _, techStack := range techStacks {
		response := FromModel(techStack)
		if response != nil {
			techStackResponses = append(techStackResponses, response)
		}
	}

	return techStackResponses
}
