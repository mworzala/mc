package mojang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mworzala/mc/internal/pkg/util"
)

var (
	mojangApiUrl     = "https://api.mojang.com"
	sessionserverUrl = "https://sessionserver.mojang.com"
	servicesApiUrl   = "https://api.minecraftservices.com"
	profileApiUrl    = servicesApiUrl + "/minecraft/profile"
)

type Client struct {
	baseUrl      string
	userAgent    string
	accountToken string
	httpClient   *http.Client
	timeout      time.Duration
}

func NewProfileClient(idVersion string, accountToken string) *Client {
	return &Client{
		baseUrl:      profileApiUrl,
		userAgent:    util.MakeUserAgent(idVersion),
		accountToken: accountToken,
		httpClient:   http.DefaultClient,
		timeout:      10 * time.Second,
	}
}

func do[T any](c *Client, ctx context.Context, method string, url string, headers http.Header, body io.Reader) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	req.Header.Set("User-Agent", c.userAgent)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		var errorRes unauthorizedError
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
		return nil, fmt.Errorf("401 %s: %s", errorRes.Path, errorRes.ErrorMessage)
	} else if res.StatusCode == http.StatusBadRequest {
		var errorRes badRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
		return nil, fmt.Errorf("400 %s: %s", errorRes.Path, errorRes.Error)
	} else if res.StatusCode == http.StatusNotFound {
		var errorRes notFoundError
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
		return nil, fmt.Errorf("404 %s: %s", errorRes.Path, errorRes.ErrorMessage)
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

func get[T any](c *Client, ctx context.Context, endpoint string, headers http.Header) (*T, error) {
	url := c.baseUrl + endpoint
	return do[T](c, ctx, http.MethodGet, url, headers, nil)
}

func delete[T any](c *Client, ctx context.Context, endpoint string, headers http.Header) (*T, error) {
	url := c.baseUrl + endpoint
	return do[T](c, ctx, http.MethodDelete, url, headers, nil)
}

func post[T any](c *Client, ctx context.Context, endpoint string, headers http.Header, body io.Reader) (*T, error) {
	url := c.baseUrl + endpoint
	return do[T](c, ctx, http.MethodPost, url, headers, body)
}

func put[T any](c *Client, ctx context.Context, endpoint string, headers http.Header, body io.Reader) (*T, error) {
	url := c.baseUrl + endpoint
	return do[T](c, ctx, http.MethodPut, url, headers, body)
}

type unauthorizedError struct {
	Path             string `json:"path"`
	ErrorType        string `json:"errorType"`
	Error            string `json:"error"`
	ErrorMessage     string `json:"errorMessage"`
	DeveloperMessage string `json:"developerMessage"`
}

type notFoundError struct {
	Path             string `json:"path"`
	ErrorType        string `json:"errorType"`
	Error            string `json:"error"`
	ErrorMessage     string `json:"errorMessage"`
	DeveloperMessage string `json:"developerMessage"`
}

type badRequestError struct {
	Path      string `json:"path"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}
