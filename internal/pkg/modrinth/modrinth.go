package modrinth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	userAgentFormat = "mworzala/mc/%s"

	fabricApiProjectId = "P7dR8mSH"
)

var (
	prodUrl    = "https://api.modrinth.com/v2"
	stagingUrl = "https://staging-api.modrinth.com/v2"
)

type Client struct {
	baseUrl    string
	userAgent  string
	httpClient *http.Client
	timeout    time.Duration
}

func NewClient(idVersion string) *Client {
	return &Client{
		baseUrl:    prodUrl,
		userAgent:  fmt.Sprintf(userAgentFormat, idVersion),
		httpClient: http.DefaultClient,
		timeout:    10 * time.Second,
	}
}

func NewStagingClient() *Client {
	return &Client{
		baseUrl:    stagingUrl,
		userAgent:  fmt.Sprintf(userAgentFormat, "dev"),
		httpClient: http.DefaultClient,
		timeout:    10 * time.Second,
	}
}

func get[T any](c *Client, ctx context.Context, endpoint string, params url.Values) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s?%s", c.baseUrl, endpoint, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var errorRes badRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}

		return nil, fmt.Errorf("400 %s: %s", errorRes.Error, errorRes.Description)
	} else if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("500 internal server error")
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from server: %d", res.StatusCode)
	}

	var result T
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &result, nil
}

// Some common error types

type badRequestError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
