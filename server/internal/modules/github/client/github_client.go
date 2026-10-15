package client

import (
	"context"
	"encoding/json"
	"fmt"
	"leslie-blog-server/internal/modules/github/dto"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	ListUserActivity(
		ctx context.Context,
		username string,
	) ([]*dto.Activity, error)
}

type githubClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGitHubClient(
	httpClient *http.Client,
	baseURL string,
) Client {
	return &githubClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

type githubEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
}

func convertEvents(
	events []githubEvent,
) []*dto.Activity {

	result := make(
		[]*dto.Activity,
		0,
		len(events),
	)

	for _, event := range events {
		result = append(
			result,
			&dto.Activity{
				ID:        event.ID,
				Type:      event.Type,
				Repo:      event.Repo.Name,
				CreatedAt: event.CreatedAt,
			},
		)
	}

	return result
}

func (c *githubClient) ListUserActivity(
	ctx context.Context,
	username string,
) ([]*dto.Activity, error) {

	url := fmt.Sprintf(
		"%s/users/%s/events/public",
		c.baseURL,
		url.PathEscape(username),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"github api returned status %d",
			resp.StatusCode,
		)
	}

	var events []githubEvent

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	return convertEvents(events), nil
}
