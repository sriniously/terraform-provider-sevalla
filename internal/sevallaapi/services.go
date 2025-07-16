package sevallaapi

import (
	"context"
	"fmt"
)

// ApplicationService handles application-related API operations.
type ApplicationService struct {
	client *Client
}

// NewApplicationService creates a new ApplicationService instance with the provided client.
func NewApplicationService(client *Client) *ApplicationService {
	return &ApplicationService{client: client}
}

func (s *ApplicationService) List(ctx context.Context) ([]Application, error) {
	var apps []Application
	err := s.client.Get(ctx, "/applications", &apps)
	return apps, err
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

func (s *DatabaseService) List(ctx context.Context) ([]Database, error) {
	var dbs []Database
	err := s.client.Get(ctx, "/databases", &dbs)
	return dbs, err
}

func (s *DatabaseService) Get(ctx context.Context, id string) (*Database, error) {
	var db Database
	err := s.client.Get(ctx, fmt.Sprintf("/databases/%s", id), &db)
	return &db, err
}

func (s *DatabaseService) Create(ctx context.Context, req CreateDatabaseRequest) (*Database, error) {
	var db Database
	err := s.client.Post(ctx, "/databases", req, &db)
	return &db, err
}

func (s *DatabaseService) Update(ctx context.Context, id string, req UpdateDatabaseRequest) (*Database, error) {
	var db Database
	err := s.client.Put(ctx, fmt.Sprintf("/databases/%s", id), req, &db)
	return &db, err
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

func (s *StaticSiteService) List(ctx context.Context) ([]StaticSite, error) {
	var sites []StaticSite
	err := s.client.Get(ctx, "/static-sites", &sites)
	return sites, err
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

// ObjectStorageService handles object storage-related API operations.
type ObjectStorageService struct {
	client *Client
}

// NewObjectStorageService creates a new ObjectStorageService instance with the provided client.
func NewObjectStorageService(client *Client) *ObjectStorageService {
	return &ObjectStorageService{client: client}
}

func (s *ObjectStorageService) List(ctx context.Context) ([]ObjectStorage, error) {
	var buckets []ObjectStorage
	err := s.client.Get(ctx, "/object-storage", &buckets)
	return buckets, err
}

func (s *ObjectStorageService) Get(ctx context.Context, id string) (*ObjectStorage, error) {
	var bucket ObjectStorage
	err := s.client.Get(ctx, fmt.Sprintf("/object-storage/%s", id), &bucket)
	return &bucket, err
}

func (s *ObjectStorageService) Create(ctx context.Context, req CreateObjectStorageRequest) (*ObjectStorage, error) {
	var bucket ObjectStorage
	err := s.client.Post(ctx, "/object-storage", req, &bucket)
	return &bucket, err
}

func (s *ObjectStorageService) Update(
	ctx context.Context,
	id string,
	req UpdateObjectStorageRequest,
) (*ObjectStorage, error) {
	var bucket ObjectStorage
	err := s.client.Put(ctx, fmt.Sprintf("/object-storage/%s", id), req, &bucket)
	return &bucket, err
}

func (s *ObjectStorageService) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/object-storage/%s", id))
}

// PipelineService handles pipeline-related API operations.
type PipelineService struct {
	client *Client
}

// NewPipelineService creates a new PipelineService instance with the provided client.
func NewPipelineService(client *Client) *PipelineService {
	return &PipelineService{client: client}
}

func (s *PipelineService) List(ctx context.Context) ([]Pipeline, error) {
	var pipelines []Pipeline
	err := s.client.Get(ctx, "/pipelines", &pipelines)
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