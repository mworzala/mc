package mojang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgentFormat = "mworzala/mc/%s"
)

var (
	mojangApiUrl     = "https://api.mojang.com/"
	sessionserverUrl = "https://sessionserver.mojang.com/"
	servicesApiUrl   = "https://api.minecraftservices.com/"
	profileApiUrl    = servicesApiUrl + "minecraft/profile"
)

type Client struct {
	baseUrl    string
	userAgent  string
	httpClient *http.Client
	timeout    time.Duration
}

func NewProfileClient(idVersion string) *Client {
	return &Client{
		baseUrl:    profileApiUrl,
		userAgent:  fmt.Sprintf(userAgentFormat, idVersion),
		httpClient: http.DefaultClient,
		timeout:    10 * time.Second,
	}
}

func get[T any](c *Client, ctx context.Context, endpoint string, headers http.Header) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", c.baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
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

func delete[T any](c *Client, ctx context.Context, endpoint string, headers http.Header) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", c.baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullUrl, nil)
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

func post[T any](c *Client, ctx context.Context, endpoint string, headers http.Header, body io.Reader) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", c.baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullUrl, body)
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

func put[T any](c *Client, ctx context.Context, endpoint string, headers http.Header, body io.Reader) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", c.baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fullUrl, body)
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
