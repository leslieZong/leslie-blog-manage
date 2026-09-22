package dto

import "leslie-blog-server/internal/modules/techstack/model"

// PublicTechStackResponse 是 Blog 前台使用的 TechStack 数据。
//
// 注意：
//
// Public API 不需要暴露：
//
// - status
// - sort
// - createdAt
// - updatedAt
//
// 这些属于管理数据。
type PublicTechStackResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Icon *string `json:"icon"`

	Description *string `json:"description"`

	OfficialURL *string `json:"officialUrl"`

	Category *string `json:"category"`
}

func FromPublicModel(
	techStack *model.TechStack,
) *PublicTechStackResponse {

	if techStack == nil {
		return nil
	}

	return &PublicTechStackResponse{
		ID:          techStack.ID,
		Name:        techStack.Name,
		Slug:        techStack.Slug,
		Icon:        techStack.Icon,
		Description: techStack.Description,
		OfficialURL: techStack.OfficialURL,
		Category:    techStack.Category,
	}
}

func FromPublicModelList(
	techStacks []*model.TechStack,
) []*PublicTechStackResponse {

	if techStacks == nil {
		return nil
	}

	techStackResponses := make([]*PublicTechStackResponse, 0, len(techStacks))

	for _, techStack := range techStacks {
		response := FromPublicModel(techStack)
		if response != nil {
			techStackResponses = append(techStackResponses, response)
		}
	}

	return techStackResponses
}
