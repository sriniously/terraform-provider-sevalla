package sevallaapi

import (
	"context"
	"fmt"
	"time"
)

// ApplicationService handles application-related API operations.
type ApplicationService struct {
	client *Client
}

// NewApplicationService creates a new ApplicationService instance with the provided client.
func NewApplicationService(client *Client) *ApplicationService {
	return &ApplicationService{client: client}
}

func (s *ApplicationService) List(ctx context.Context, companyID string) ([]ApplicationListItem, error) {
	var response ApplicationListResponse
	url := fmt.Sprintf("/applications?company=%s", companyID)
	err := s.client.Get(ctx, url, &response)
	return response.Company.Apps.Items, err
}

func (s *ApplicationService) Get(ctx context.Context, id string) (*Application, error) {
	var app Application
	err := s.client.Get(ctx, fmt.Sprintf("/applications/%s", id), &app)
	return &app, err
}

func (s *ApplicationService) Create(ctx context.Context, req CreateApplicationRequest) (*Application, error) {
	var app Application
	err := s.client.Post(ctx, "/applications", req, &app)
	return &app, err
}

func (s *ApplicationService) Update(
	ctx context.Context,
	id string,
	req UpdateApplicationRequest,
) (*Application, error) {
	var app Application
	err := s.client.Put(ctx, fmt.Sprintf("/applications/%s", id), req, &app)
	return &app, err
}

func (s *ApplicationService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/applications/%s", id))
}

// DatabaseService handles database-related API operations.
type DatabaseService struct {
	client *Client
}

// NewDatabaseService creates a new DatabaseService instance with the provided client.
func NewDatabaseService(client *Client) *DatabaseService {
	return &DatabaseService{client: client}
}

func (s *DatabaseService) List(ctx context.Context, companyID string) ([]DatabaseListItem, error) {
	var response DatabaseListResponse
	url := fmt.Sprintf("/databases?company=%s", companyID)
	err := s.client.Get(ctx, url, &response)
	if err != nil {
		return nil, err
	}
	return response.Company.Databases.Items, nil
}

func (s *DatabaseService) Get(ctx context.Context, id string) (*Database, error) {
	var db Database
	// Based on OpenAPI spec, the database GET endpoint requires internal and external query parameters
	url := fmt.Sprintf("/databases/%s?internal=true&external=true", id)
	err := s.client.Get(ctx, url, &db)
	return &db, err
}

func (s *DatabaseService) Create(ctx context.Context, req CreateDatabaseRequest) (*Database, error) {
	// The create endpoint only returns the database ID
	var createResp struct {
		Database struct {
			ID string `json:"id"`
		} `json:"database"`
	}
	err := s.client.Post(ctx, "/databases", req, &createResp)
	if err != nil {
		return nil, err
	}
	
	// Retry getting the database details up to 3 times with a short delay
	// as the database might not be immediately available after creation
	var db *Database
	for i := 0; i < 3; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		db, err = s.Get(ctx, createResp.Database.ID)
		if err == nil {
			break
		}
	}
	
	return db, err
}

func (s *DatabaseService) Update(ctx context.Context, id string, req UpdateDatabaseRequest) (*Database, error) {
	// The update endpoint returns limited information
	var updateResp struct {
		Database struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
			Status      string `json:"status"`
		} `json:"database"`
	}
	err := s.client.Put(ctx, fmt.Sprintf("/databases/%s", id), req, &updateResp)
	if err != nil {
		return nil, err
	}
	// Fetch the full database details
	return s.Get(ctx, id)
}

func (s *DatabaseService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/databases/%s", id))
}

// StaticSiteService handles static site-related API operations.
type StaticSiteService struct {
	client *Client
}

// NewStaticSiteService creates a new StaticSiteService instance with the provided client.
func NewStaticSiteService(client *Client) *StaticSiteService {
	return &StaticSiteService{client: client}
}

func (s *StaticSiteService) List(ctx context.Context, companyID string) ([]StaticSiteListItem, error) {
	var response StaticSiteListResponse
	url := fmt.Sprintf("/static-sites?company=%s", companyID)
	err := s.client.Get(ctx, url, &response)
	if err != nil {
		return nil, err
	}
	return response.Company.StaticSites.Items, nil
}

func (s *StaticSiteService) Get(ctx context.Context, id string) (*StaticSite, error) {
	var site StaticSite
	err := s.client.Get(ctx, fmt.Sprintf("/static-sites/%s", id), &site)
	return &site, err
}

func (s *StaticSiteService) Create(ctx context.Context, req CreateStaticSiteRequest) (*StaticSite, error) {
	var site StaticSite
	err := s.client.Post(ctx, "/static-sites", req, &site)
	return &site, err
}

func (s *StaticSiteService) Update(ctx context.Context, id string, req UpdateStaticSiteRequest) (*StaticSite, error) {
	var site StaticSite
	err := s.client.Put(ctx, fmt.Sprintf("/static-sites/%s", id), req, &site)
	return &site, err
}

func (s *StaticSiteService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/static-sites/%s", id))
}

// PipelineService handles pipeline-related API operations.
type PipelineService struct {
	client *Client
}

// NewPipelineService creates a new PipelineService instance with the provided client.
func NewPipelineService(client *Client) *PipelineService {
	return &PipelineService{client: client}
}

func (s *PipelineService) List(ctx context.Context, companyID string) ([]Pipeline, error) {
	var pipelines []Pipeline
	url := fmt.Sprintf("/pipelines?company=%s", companyID)
	err := s.client.Get(ctx, url, &pipelines)
	return pipelines, err
}

func (s *PipelineService) Get(ctx context.Context, id string) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.Get(ctx, fmt.Sprintf("/pipelines/%s", id), &pipeline)
	return &pipeline, err
}

func (s *PipelineService) Create(ctx context.Context, req CreatePipelineRequest) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.Post(ctx, "/pipelines", req, &pipeline)
	return &pipeline, err
}

func (s *PipelineService) Update(ctx context.Context, id string, req UpdatePipelineRequest) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.Put(ctx, fmt.Sprintf("/pipelines/%s", id), req, &pipeline)
	return &pipeline, err
}

func (s *PipelineService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/pipelines/%s", id))
}

// DeploymentService handles deployment-related API operations.
type DeploymentService struct {
	client *Client
}

// NewDeploymentService creates a new DeploymentService instance with the provided client.
func NewDeploymentService(client *Client) *DeploymentService {
	return &DeploymentService{client: client}
}

func (s *DeploymentService) List(ctx context.Context, appID string) ([]Deployment, error) {
	var deployments []Deployment
	err := s.client.Get(ctx, fmt.Sprintf("/applications/%s/deployments", appID), &deployments)
	return deployments, err
}

func (s *DeploymentService) Get(ctx context.Context, appID, deploymentID string) (*Deployment, error) {
	var deployment Deployment
	err := s.client.Get(ctx, fmt.Sprintf("/applications/%s/deployments/%s", appID, deploymentID), &deployment)
	return &deployment, err
}

// SiteService handles WordPress site-related API operations.
type SiteService struct {
	client *Client
}

// NewSiteService creates a new SiteService instance with the provided client.
func NewSiteService(client *Client) *SiteService {
	return &SiteService{client: client}
}

func (s *SiteService) List(ctx context.Context, companyID string) ([]SiteListItem, error) {
	var response SiteListResponse
	url := fmt.Sprintf("/sites?company=%s", companyID)
	err := s.client.Get(ctx, url, &response)
	if err != nil {
		return nil, err
	}
	return response.Company.Sites, nil
}

func (s *SiteService) Get(ctx context.Context, id string) (*Site, error) {
	var site Site
	err := s.client.Get(ctx, fmt.Sprintf("/sites/%s", id), &site)
	return &site, err
}

func (s *SiteService) Create(ctx context.Context, req CreateSiteRequest) (*OperationResponse, error) {
	var opResp OperationResponse
	err := s.client.Post(ctx, "/sites", req, &opResp)
	return &opResp, err
}

func (s *SiteService) Update(ctx context.Context, id string, req UpdateSiteRequest) (*Site, error) {
	var site Site
	err := s.client.Put(ctx, fmt.Sprintf("/sites/%s", id), req, &site)
	return &site, err
}

func (s *SiteService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/sites/%s", id))
}

// CompanyService handles company-related API operations.
type CompanyService struct {
	client *Client
}

// NewCompanyService creates a new CompanyService instance with the provided client.
func NewCompanyService(client *Client) *CompanyService {
	return &CompanyService{client: client}
}

func (s *CompanyService) GetUsers(ctx context.Context, companyID string) (*CompanyUsers, error) {
	var users CompanyUsers
	err := s.client.Get(ctx, fmt.Sprintf("/company/%s/users", companyID), &users)
	return &users, err
}

// OperationService handles operation-related API operations.
type OperationService struct {
	client *Client
}

// NewOperationService creates a new OperationService instance with the provided client.
func NewOperationService(client *Client) *OperationService {
	return &OperationService{client: client}
}

func (s *OperationService) GetStatus(ctx context.Context, operationID string) (*Operation, error) {
	var op Operation
	err := s.client.Get(ctx, fmt.Sprintf("/operations/%s", operationID), &op)
	return &op, err
}
