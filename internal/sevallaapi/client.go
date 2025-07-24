// Package sevallaapi provides a client for the Sevalla API.
package sevallaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.sevalla.com/v2"
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string

	// Services
	Applications *ApplicationService
	Databases    *DatabaseService
	StaticSites  *StaticSiteService
	Sites        *SiteService
	Pipelines    *PipelineService
	Deployments  *DeploymentService
	Company      *CompanyService
	Operations   *OperationService
}

type Config struct {
	BaseURL string
	Token   string
	Timeout time.Duration
}

// NewClient creates a new Sevalla API client with the provided configuration.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	client := &Client{
		BaseURL: config.BaseURL,
		HTTPClient: &http.Client{
			Timeout: config.Timeout,
		},
		Token: config.Token,
	}

	// Initialize services
	client.Applications = NewApplicationService(client)
	client.Databases = NewDatabaseService(client)
	client.StaticSites = NewStaticSiteService(client)
	client.Sites = NewSiteService(client)
	client.Pipelines = NewPipelineService(client)
	client.Deployments = NewDeploymentService(client)
	client.Company = NewCompanyService(client)
	client.Operations = NewOperationService(client)

	return client
}

func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	reqURL, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	const httpBadRequestThreshold = 400
	if resp.StatusCode >= httpBadRequestThreshold {
		return c.handleError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	const httpBadRequestThreshold = 400
	if resp.StatusCode >= httpBadRequestThreshold {
		return c.handleError(resp)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.makeRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	const httpBadRequestThreshold = 400
	if resp.StatusCode >= httpBadRequestThreshold {
		return c.handleError(resp)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	const httpBadRequestThreshold = 400
	if resp.StatusCode >= httpBadRequestThreshold {
		return c.handleError(resp)
	}

	return nil
}

func (c *Client) handleError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var errorResponse struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if errorResponse.Message != "" {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResponse.Message)
	}
	if errorResponse.Error != "" {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResponse.Error)
	}

	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

// Pipeline convenience methods.
func (c *Client) GetPipeline(ctx context.Context, id string) (*Pipeline, error) {
	return c.Pipelines.Get(ctx, id)
}

func (c *Client) CreatePipeline(ctx context.Context, req CreatePipelineRequest) (*Pipeline, error) {
	return c.Pipelines.Create(ctx, req)
}

func (c *Client) UpdatePipeline(ctx context.Context, id string, req UpdatePipelineRequest) (*Pipeline, error) {
	return c.Pipelines.Update(ctx, id, req)
}

func (c *Client) DeletePipeline(ctx context.Context, id string) error {
	return c.Pipelines.Delete(ctx, id)
}
