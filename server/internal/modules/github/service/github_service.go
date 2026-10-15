package service

import (
	"context"
	"encoding/json"
	"leslie-blog-server/internal/modules/github/client"
	"leslie-blog-server/internal/modules/github/dto"
	"leslie-blog-server/internal/pkg/cache"
	"time"
)

type GitHubService interface {
	GetPublicActivity(
		ctx context.Context,
	) (*dto.ActivityResponse, error)
}

type gitHubService struct {
	client   client.Client
	cache    cache.Cache
	username string
	ttl      time.Duration
}

// NewGitHubService 创建 GitHub 服务。
func NewGitHubService(
	client client.Client,
	cache cache.Cache,
	username string,
	ttl time.Duration,
) GitHubService {
	return &gitHubService{
		client:   client,
		cache:    cache,
		username: username,
		ttl:      ttl,
	}
}

func (s *gitHubService) GetPublicActivity(
	ctx context.Context,
) (*dto.ActivityResponse, error) {

	cached, err := s.cache.Get(
		ctx,
		cache.KeyGitHubActivity,
	)

	if err == nil {
		var response dto.ActivityResponse

		if json.Unmarshal(
			[]byte(cached),
			&response,
		) == nil {
			return &response, nil
		}
	}

	activities, err := s.client.ListUserActivity(
		ctx,
		s.username,
	)
	if err != nil {
		return nil, err
	}

	response := &dto.ActivityResponse{
		Items: activities,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	// Cache 写失败不能影响业务返回。
	_ = s.cache.Set(
		ctx,
		cache.KeyGitHubActivity,
		string(data),
		s.ttl,
	)

	return response, nil
}
