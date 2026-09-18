package dto

import "leslie-blog-server/internal/modules/tag/model"

// SimpleTagResponse 是文章中嵌入的 Tag 简化信息。
//
// 注意：
//
// 文章详情不需要把 Tag 的所有管理字段全部返回。
//
// 只需要：
//
// id
// name
// slug
type SimpleTagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// FromSimpleModel Model → SimpleTagResponse。
func FromSimpleModel(
	tag *model.Tag,
) *SimpleTagResponse {

	if tag == nil {
		return nil
	}

	return &SimpleTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}
}

// FromSimpleModelList Model → SimpleTagResponse 列表。
func FromSimpleModelList(
	tags []*model.Tag,
) []*SimpleTagResponse {

	if tags == nil {
		return nil
	}

	simpleTags := make(
		[]*SimpleTagResponse,
		0,
		len(tags),
	)

	for _, tag := range tags {
		simpleTags = append(simpleTags, FromSimpleModel(tag))
	}

	return simpleTags
}
