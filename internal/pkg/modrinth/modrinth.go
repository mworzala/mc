package modrinth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	userAgentFormat = "mworzala/mc-cli/%s"

	fabricApiProjectId = "P7dR8mSH"
)

var (
	baseUrl     = "https://api.modrinth.com/v2"
	getVersions = fmt.Sprintf("%s/project/%%s/version", baseUrl)
)

type Client struct {
	userAgent  string
	httpClient *http.Client
	timeout    time.Duration
}

func NewClient(idVersion string) *Client {
	return &Client{
		userAgent:  fmt.Sprintf(userAgentFormat, idVersion),
		httpClient: http.DefaultClient,
		timeout:    10 * time.Second,
	}
}

func (c *Client) GetVersions(ctx context.Context, projectId, loader, gameVersion string) ([]*Version, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	endpoint, err := url.Parse(fmt.Sprintf(getVersions, projectId))
	if err != nil {
		return nil, err
	}
	endpoint.Query().Set("loaders", fmt.Sprintf("[%s]", loader))
	endpoint.Query().Set("game_versions", fmt.Sprintf("[%s]", gameVersion))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result []*Version
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return result, nil
}
