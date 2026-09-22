package repository

import (
	"context"
	"leslie-blog-server/internal/modules/project/model"

	"gorm.io/gorm"
)

type projectTechStackRepository struct {
	db *gorm.DB
}

// NewProjectTechStackRepository 创建关系 Repository。
func NewProjectTechStackRepository(
	db *gorm.DB,
) ProjectTechStackRepository {

	return &projectTechStackRepository{
		db: db,
	}
}

func (r *projectTechStackRepository) ReplaceTechStacks(
	ctx context.Context,
	projectID string,
	techStackIDs []string,
) error {

	// -------------------------------------------------
	// 第一步：删除当前 Project 的所有 TechStack 关系
	// -------------------------------------------------

	if err := r.db.
		WithContext(ctx).
		Where(
			"project_id = ?",
			projectID,
		).
		Delete(
			&model.ProjectTechStack{},
		).
		Error; err != nil {

		return err
	}

	// -------------------------------------------------
	// 第二步：如果没有新的 TechStack，直接结束。
	//
	// []：
	//
	// 表示这个 Project 不使用任何 TechStack。
	// -------------------------------------------------

	if len(techStackIDs) == 0 {
		return nil
	}

	// -------------------------------------------------
	// 第三步：构建新的关系数据
	// -------------------------------------------------

	relations := make(
		[]*model.ProjectTechStack,
		0,
		len(techStackIDs),
	)

	for index, techStackID := range techStackIDs {

		relations = append(
			relations,
			&model.ProjectTechStack{

				ProjectID: projectID,

				TechStackID: techStackID,

				// 数组下标从 0 开始，
				// 我们的 sort 从 1 开始。
				Sort: index + 1,
			},
		)
	}

	// -------------------------------------------------
	// 第四步：批量插入关系
	// -------------------------------------------------

	return r.db.
		WithContext(ctx).
		Create(
			&relations,
		).
		Error
}

func (r *projectTechStackRepository) FindByProjectID(
	ctx context.Context,
	projectID string,
) ([]*model.ProjectTechStack, error) {

	var relations []*model.ProjectTechStack

	err := r.db.
		WithContext(ctx).
		Where(
			"project_id = ?",
			projectID,
		).
		Order("sort ASC").
		Find(&relations).
		Error

	if err != nil {
		return nil, err
	}

	return relations, nil
}
