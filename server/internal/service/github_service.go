package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"leslie-blog-server/config"
)

type GitHubService struct {
	username string
	token    string
	client   *http.Client
}

func NewGitHubService(cfg config.GitHubConfig) *GitHubService {
	return &GitHubService{
		username: cfg.Username,
		token:    cfg.Token,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *GitHubService) do(path string, v interface{}) error {
	if s.username == "" {
		return errors.New("github username 未配置")
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com"+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api %d: %s", resp.StatusCode, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// Profile 用户信息
func (s *GitHubService) Profile() (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := s.do(fmt.Sprintf("/users/%s", s.username), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Repos 仓库列表
func (s *GitHubService) Repos(perPage int) ([]map[string]interface{}, error) {
	if perPage <= 0 || perPage > 100 {
		perPage = 30
	}
	var m []map[string]interface{}
	path := fmt.Sprintf("/users/%s/repos?sort=updated&per_page=%d&type=owner", s.username, perPage)
	if err := s.do(path, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Contributions 最近公开事件（GitHub REST 无贡献图，用 events 近似）
func (s *GitHubService) Contributions(perPage int) ([]map[string]interface{}, error) {
	if perPage <= 0 || perPage > 100 {
		perPage = 30
	}
	var m []map[string]interface{}
	path := fmt.Sprintf("/users/%s/events/public?per_page=%d", s.username, perPage)
	if err := s.do(path, &m); err != nil {
		return nil, err
	}
	return m, nil
}
