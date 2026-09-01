package service

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
)

type SearchService struct {
	postRepo    *repository.PostRepo
	projectRepo *repository.ProjectRepo
}

func NewSearchService(postRepo *repository.PostRepo, projectRepo *repository.ProjectRepo) *SearchService {
	return &SearchService{postRepo: postRepo, projectRepo: projectRepo}
}

// Search 聚合检索：文章（已发布 + locale）+ 项目
func (s *SearchService) Search(keyword, locale string, page, pageSize int) (*model.SearchResult, error) {
	posts, total, err := s.postRepo.Search(keyword, locale, page, pageSize)
	if err != nil {
		return nil, err
	}
	// 扁平化文章
	postDTOs := make([]model.PostDTO, 0, len(posts))
	ps := NewPostService(s.postRepo, nil) // 仅用于复用 toDTO
	for i := range posts {
		postDTOs = append(postDTOs, ps.toDTO(posts[i], locale))
	}

	projects, err := s.projectRepo.ListAll()
	if err != nil {
		return nil, err
	}
	// 简单按 keyword 过滤项目名/描述
	filtered := make([]model.Project, 0, len(projects))
	for _, p := range projects {
		if contains(p.Name, keyword) || contains(p.Description, keyword) {
			filtered = append(filtered, p)
		}
	}
	return &model.SearchResult{
		Posts:     postDTOs,
		Projects:  filtered,
		Total:      total + int64(len(filtered)),
	}, nil
}

func contains(s, sub string) bool {
	return len(sub) == 0 || indexOf(s, sub) >= 0
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
