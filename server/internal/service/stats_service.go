package service

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
)

type StatsService struct{ repos *repository.Container }

func NewStatsService(repos *repository.Container) *StatsService {
	return &StatsService{repos: repos}
}

// Stats 汇总统计
func (s *StatsService) Stats() (*model.Stats, error) {
	var st model.Stats
	if err := s.count(&st.PostCount, s.repos.Post, nil); err != nil {
		return nil, err
	}
	published := int8(1)
	if err := s.count(&st.PublishedCount, s.repos.Post, &published); err != nil {
		return nil, err
	}
	// 各表行数
	st.CategoryCount = s.tableCount("categories")
	st.ProjectCount = s.tableCount("projects")
	st.TechStackCount = s.tableCount("tech_stack")
	st.MediaCount = s.tableCount("media")
	st.CommentCount = s.tableCount("comments")
	// 总浏览量：sum(view_count)
	var views *int64
	row := repository.DB.Raw("SELECT COALESCE(SUM(view_count),0) FROM posts").Row()
	if err := row.Scan(&views); err == nil && views != nil {
		st.ViewCount = *views
	}
	return &st, nil
}

// RecentPosts 最近文章（默认 locale 标题）
func (s *StatsService) RecentPosts(limit int) ([]model.RecentItem, error) {
	posts, _, err := s.repos.Post.List(repository.PostFilter{}, 1, limit)
	if err != nil {
		return nil, err
	}
	out := make([]model.RecentItem, 0, len(posts))
	for _, p := range posts {
		title := ""
		if tr := pickTranslation(p.Translations, "zh-CN"); tr != nil {
			title = tr.Title
		}
		out = append(out, model.RecentItem{
			ID: p.ID, Title: title, CreatedAt: p.CreatedAt, Status: p.Status,
		})
	}
	return out, nil
}

// RecentComments 最近评论
func (s *StatsService) RecentComments(limit int) ([]model.RecentItem, error) {
	comments, _, err := s.repos.Comment.List("", nil, 1, limit)
	if err != nil {
		return nil, err
	}
	out := make([]model.RecentItem, 0, len(comments))
	for _, c := range comments {
		out = append(out, model.RecentItem{
			ID: c.ID, Title: c.Content, CreatedAt: c.CreatedAt, Status: c.Status,
		})
	}
	return out, nil
}

// count 复用 postRepo 的 List total
func (s *StatsService) count(dst *int64, repo *repository.PostRepo, status *int8) error {
	_, total, err := repo.List(repository.PostFilter{Status: status}, 1, 1)
	if err != nil {
		return err
	}
	*dst = total
	return nil
}

func (s *StatsService) tableCount(table string) int64 {
	var n int64
	if repository.DB != nil {
		repository.DB.Table(table).Count(&n)
	}
	return n
}
