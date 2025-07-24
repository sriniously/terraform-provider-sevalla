package sevallaapi

// Note: Using int64 for timestamps instead of time.Time to match API responses

// Application represents a Sevalla application based on MKApplicationSchema.
type Application struct {
	App ApplicationDetails `json:"app"`
}

// ApplicationDetails represents the actual application data.
type ApplicationDetails struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	DisplayName          string               `json:"display_name"`
	Status               string               `json:"status"`
	CompanyID            string               `json:"company_id"`
	RepoURL              string               `json:"repo_url"`
	DefaultBranch        string               `json:"default_branch"`
	AutoDeploy           bool                 `json:"auto_deploy"`
	BuildPath            string               `json:"build_path"`
	BuildType            string               `json:"build_type"`
	NodeVersion          string               `json:"node_version,omitempty"`
	DockerfilePath       string               `json:"dockerfile_path,omitempty"`
	DockerComposeFile    string               `json:"docker_compose_file,omitempty"`
	StartCommand         string               `json:"start_command,omitempty"`
	InstallCommand       string               `json:"install_command,omitempty"`
	EnvironmentVariables []EnvVar             `json:"environment_variables,omitempty"`
	CreatedAt            int64                `json:"created_at"`
	UpdatedAt            int64                `json:"updated_at"`
	Deployments          []AppDeployment      `json:"deployments,omitempty"`
	Processes            []AppProcess         `json:"processes,omitempty"`
	InternalConnections  []InternalConnection `json:"internal_connections,omitempty"`
}

// ApplicationListItem represents an application in a list response.
type ApplicationListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

// ApplicationListResponse represents the response from the applications list endpoint.
type ApplicationListResponse struct {
	Company struct {
		Apps struct {
			Items []ApplicationListItem `json:"items"`
		} `json:"apps"`
	} `json:"company"`
}

// AppDeployment represents a deployment within an application.
type AppDeployment struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	Branch        string  `json:"branch"`
	RepoURL       string  `json:"repo_url"`
	CommitHash    string  `json:"commit_hash,omitempty"`
	CommitMessage *string `json:"commit_message"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at,omitempty"`
	BuildLogs     string  `json:"build_logs,omitempty"`
}

// AppProcess represents a process within an application.
type AppProcess struct {
	ID               string           `json:"id"`
	Key              string           `json:"key"`
	Type             string           `json:"type"`
	DisplayName      string           `json:"display_name"`
	ScalingStrategy  *ScalingStrategy `json:"scaling_strategy,omitempty"`
	ResourceTypeName string           `json:"resource_type_name"`
	Entrypoint       string           `json:"entrypoint"`
}

// Process represents a detailed application process.
type Process struct {
	Process ProcessDetails `json:"process"`
}

// ProcessDetails represents the actual process data.
type ProcessDetails struct {
	ID               string           `json:"id"`
	Type             string           `json:"type"`
	DisplayName      string           `json:"display_name"`
	ScalingStrategy  *ScalingStrategy `json:"scaling_strategy,omitempty"`
	ResourceTypeName string           `json:"resource_type_name"`
	Entrypoint       string           `json:"entrypoint"`
}

// ScalingStrategy represents the scaling configuration for a process.
type ScalingStrategy struct {
	Type   string                 `json:"type"`   // manual or horizontal
	Config map[string]interface{} `json:"config"` // Different configs based on type
}

// CreateApplicationRequest represents the request to create an application.
// Note: Application creation appears to be handled through deployments in the API.
type CreateApplicationRequest struct {
	CompanyID   string `json:"company_id"`
	DisplayName string `json:"display_name"`
	RepoURL     string `json:"repo_url"`
	Branch      string `json:"branch,omitempty"`
	// Add other fields as needed based on API documentation
}

// CreateDeploymentRequest represents the request to create a deployment.
type CreateDeploymentRequest struct {
	Branch        string `json:"branch,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
}

// UpdateApplicationRequest represents the request to update an application.
type UpdateApplicationRequest struct {
	DisplayName          *string      `json:"display_name,omitempty"`
	BuildPath            *string      `json:"build_path,omitempty"`
	BuildType            *BuildType   `json:"build_type,omitempty"`
	DefaultBranch        *string      `json:"default_branch,omitempty"`
	AutoDeploy           *bool        `json:"auto_deploy,omitempty"`
	NodeVersion          *NodeVersion `json:"node_version,omitempty"`
	DockerfilePath       *string      `json:"dockerfile_path,omitempty"`
	DockerComposeFile    *string      `json:"docker_compose_file,omitempty"`
	PackConfig           *PackConfig  `json:"pack_config,omitempty"`
	EnvironmentVariables []EnvVar     `json:"environment_variables,omitempty"`
	StartCommand         *string      `json:"start_command,omitempty"`
	InstallCommand       *string      `json:"install_command,omitempty"`
}

// PackConfig represents configuration for pack-based builds.
type PackConfig struct {
	Builder string `json:"builder"`
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Database represents a Sevalla database from the detailed view.
type Database struct {
	Database DatabaseDetails `json:"database"`
}

// DatabaseDetails represents the actual database data.
type DatabaseDetails struct {
	ID                       string               `json:"id"`
	Name                     string               `json:"name"`
	DisplayName              string               `json:"display_name"`
	Status                   string               `json:"status"`
	CreatedAt                int64                `json:"created_at"`
	MemoryLimit              int                  `json:"memory_limit"`
	CPULimit                 int                  `json:"cpu_limit"`
	StorageSize              int                  `json:"storage_size"`
	Type                     string               `json:"type"`
	Version                  string               `json:"version"`
	Cluster                  DatabaseCluster      `json:"cluster"`
	ResourceTypeName         string               `json:"resource_type_name"`
	InternalHostname         *string              `json:"internal_hostname"`
	InternalPort             *string              `json:"internal_port"`
	InternalConnections      []DatabaseConnection `json:"internal_connections"`
	Data                     DatabaseData         `json:"data"`
	ExternalConnectionString string               `json:"external_connection_string"`
	ExternalHostname         *string              `json:"external_hostname"`
	ExternalPort             *string              `json:"external_port"`
}

// DatabaseListItem represents a database in a list response.
type DatabaseListItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	DisplayName      string `json:"display_name"`
	Status           string `json:"status"`
	UpdatedAt        int64  `json:"updated_at"`
	Type             string `json:"type"`
	Version          string `json:"version"`
	ResourceTypeName string `json:"resource_type_name"`
}

// DatabaseCluster represents the cluster information for a database.
type DatabaseCluster struct {
	ID          string `json:"id"`
	Location    string `json:"location"`
	DisplayName string `json:"display_name"`
}

// DatabaseConnection represents an internal connection to a database.
type DatabaseConnection struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// DatabaseData represents the database credentials and configuration.
type DatabaseData struct {
	DBName         string  `json:"db_name"`
	DBPassword     string  `json:"db_password"`
	DBRootPassword *string `json:"db_root_password"`
	DBUser         *string `json:"db_user"`
}

// CreateDatabaseRequest represents the request to create a database.
type CreateDatabaseRequest struct {
	CompanyID    string `json:"company_id"`
	Location     string `json:"location"`
	ResourceType string `json:"resource_type"` // db1, db2, ..., db9
	DisplayName  string `json:"display_name"`
	DBName       string `json:"db_name"`
	DBPassword   string `json:"db_password"`
	DBUser       string `json:"db_user,omitempty"` // Optional for Redis, required for others
	Type         string `json:"type"`              // postgresql, redis, mariadb, mysql
	Version      string `json:"version"`
}

// UpdateDatabaseRequest represents the request to update a database.
type UpdateDatabaseRequest struct {
	DisplayName  *string `json:"display_name,omitempty"`
	ResourceType *string `json:"resource_type,omitempty"`
	// Add other updateable fields based on API specification
}

// StaticSite represents a Sevalla static site from the detailed view.
type StaticSite struct {
	StaticSite StaticSiteDetails `json:"static_site"`
}

// StaticSiteDetails represents the actual static site data.
type StaticSiteDetails struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	DisplayName        string                 `json:"display_name"`
	Status             string                 `json:"status"`
	RepoURL            string                 `json:"repo_url"`
	DefaultBranch      string                 `json:"default_branch"`
	AutoDeploy         bool                   `json:"auto_deploy"`
	RemoteRepositoryID string                 `json:"remote_repository_id"`
	GitRepositoryID    string                 `json:"git_repository_id"`
	GitType            string                 `json:"git_type"`
	Hostname           string                 `json:"hostname"`
	BuildCommand       *string                `json:"build_command"`
	CreatedAt          int64                  `json:"created_at"`
	UpdatedAt          int64                  `json:"updated_at"`
	Deployments        []StaticSiteDeployment `json:"deployments,omitempty"`
}

// StaticSiteListItem represents a static site in a list response.
type StaticSiteListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

// StaticSiteDeployment represents a deployment within a static site.
type StaticSiteDeployment struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	RepoURL       string  `json:"repo_url"`
	Branch        string  `json:"branch"`
	CommitMessage *string `json:"commit_message"`
	CreatedAt     int64   `json:"created_at"`
}

// CreateStaticSiteRequest represents the request to create a static site.
// Note: Static site creation appears to be handled through deployments in the API.
type CreateStaticSiteRequest struct {
	CompanyID   string  `json:"company_id"`
	DisplayName string  `json:"display_name"`
	RepoURL     string  `json:"repo_url"`
	Branch      *string `json:"branch,omitempty"`
	// Add other fields as needed based on API documentation
}

// UpdateStaticSiteRequest represents the request to update a static site.
type UpdateStaticSiteRequest struct {
	DisplayName        *string `json:"display_name,omitempty"`
	AutoDeploy         *bool   `json:"auto_deploy,omitempty"`
	DefaultBranch      *string `json:"default_branch,omitempty"`
	BuildCommand       *string `json:"build_command,omitempty"`
	NodeVersion        *string `json:"node_version,omitempty"`        // 16.20.0|18.16.0|20.2.0
	PublishedDirectory *string `json:"published_directory,omitempty"` // dist
}

// Site represents a WordPress site from the detailed view.
type Site struct {
	Site SiteDetails `json:"site"`
}

// SiteDetails represents the actual site data.
type SiteDetails struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	DisplayName  string        `json:"display_name"`
	CompanyID    string        `json:"company_id"`
	Status       string        `json:"status"`
	Environments []Environment `json:"environments"`
}

// SiteListItem represents a site in a list response.
type SiteListItem struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Status      string      `json:"status"`
	SiteLabels  []SiteLabel `json:"siteLabels"`
}

// SiteLabel represents a label attached to a site.
type SiteLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Environment represents a site environment.
type Environment struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	DisplayName   string   `json:"display_name"`
	IsPremium     bool     `json:"is_premium"`
	IsBlocked     bool     `json:"is_blocked"`
	Domains       []Domain `json:"domains"`
	PrimaryDomain Domain   `json:"primaryDomain"`
}

// Domain represents a domain attached to an environment.
type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// CreateSiteRequest represents the request to create a WordPress site.
type CreateSiteRequest struct {
	CompanyID   string `json:"company_id"`
	DisplayName string `json:"display_name"`
	// Add other fields as needed based on API documentation
}

// UpdateSiteRequest represents the request to update a WordPress site.
type UpdateSiteRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	// Add other updateable fields based on API specification
}

// CompanyUsers represents the response from the company users endpoint.
type CompanyUsers struct {
	Company struct {
		Users []CompanyUser `json:"users"`
	} `json:"company"`
}

// CompanyUser represents a user within a company.
type CompanyUser struct {
	User UserDetails `json:"user"`
}

// UserDetails represents the actual user data.
type UserDetails struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	FullName string `json:"full_name"`
}

// OperationResponse represents a response for asynchronous operations.
type OperationResponse struct {
	OperationID string `json:"operation_id"`
	Message     string `json:"message"`
	Status      int    `json:"status"`
}

// Operation represents the status of an ongoing operation.
type Operation struct {
	ID          string      `json:"id"`
	Status      string      `json:"status"` // pending, running, completed, failed
	Type        string      `json:"type"`   // create_site, delete_database, etc.
	ResourceID  string      `json:"resource_id,omitempty"`
	Progress    int         `json:"progress"` // 0-100
	Message     string      `json:"message"`
	CreatedAt   int64       `json:"created_at"`
	CompletedAt *int64      `json:"completed_at,omitempty"`
	Error       *string     `json:"error,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// StatusResponse represents a standard API status response.
type StatusResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

// Deployment represents a deployment - this might need adjustment based on actual API.
type Deployment struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	Branch        string `json:"branch"`
	CommitHash    string `json:"commit_hash,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	CreatedAt     int64  `json:"created_at"`
}

// Pipeline represents a deployment pipeline.
type Pipeline struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"display_name"`
	Stages      []PipelineStage `json:"stages"`
}

// PipelineStage represents a stage within a pipeline.
type PipelineStage struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // preview or standard
}

// CreatePipelineRequest represents the request to create a pipeline.
type CreatePipelineRequest struct {
	DisplayName string `json:"display_name"`
	// Add other fields as needed based on API documentation
}

// UpdatePipelineRequest represents the request to update a pipeline.
type UpdatePipelineRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	// Add other updateable fields based on API specification
}

// InternalConnection represents a connection between resources.
type InternalConnection struct {
	ID         string `json:"id"`
	TargetType string `json:"target_type"` // appResource, dbResource, envResource
	TargetID   string `json:"target_id"`
	CreatedAt  int64  `json:"created_at"`
}

// CreateInternalConnectionRequest represents the request to create an internal connection.
type CreateInternalConnectionRequest struct {
	TargetType string `json:"target_type"` // appResource, dbResource, envResource
	TargetID   string `json:"target_id"`
}

// CDNStatus represents CDN configuration status.
type CDNStatus struct {
	IsTurnedOn bool `json:"isTurnedOn"`
}

// EdgeCachingStatus represents edge caching configuration status.
type EdgeCachingStatus struct {
	IsTurnedOn bool `json:"isTurnedOn"`
}

// ClearCacheResponse represents the response from clearing cache.
type ClearCacheResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ApplicationMetrics represents application analytics data.
type ApplicationMetrics struct {
	Timeframe []string  `json:"timeframe"`
	Data      []float64 `json:"data"`
}

// BandwidthMetrics represents bandwidth usage metrics.
type BandwidthMetrics struct {
	Timeframe []string  `json:"timeframe"`
	Data      []float64 `json:"data"`
	Unit      string    `json:"unit"` // e.g., "bytes", "MB", "GB"
}

// BuildTimeMetrics represents build time analytics.
type BuildTimeMetrics struct {
	Timeframe []string  `json:"timeframe"`
	Data      []float64 `json:"data"`
	Unit      string    `json:"unit"` // e.g., "seconds", "minutes"
}

// RuntimeMetrics represents runtime performance metrics.
type RuntimeMetrics struct {
	Timeframe []string  `json:"timeframe"`
	Data      []float64 `json:"data"`
	Unit      string    `json:"unit"` // e.g., "ms", "seconds"
}

// HTTPRequestMetrics represents HTTP request analytics.
type HTTPRequestMetrics struct {
	Timeframe []string `json:"timeframe"`
	Data      []int64  `json:"data"`
}

// MetricsQuery represents query parameters for metrics endpoints.
type MetricsQuery struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD format
	EndDate   string `json:"end_date"`   // YYYY-MM-DD format
	Interval  string `json:"interval"`   // hour, day, week, month
}

// DatabaseListResponse represents the response from the databases list endpoint.
// Based on CompanyDatabasesSchema from the OpenAPI spec.
type DatabaseListResponse struct {
	Company struct {
		Databases struct {
			Items []DatabaseListItem `json:"items"`
		} `json:"databases"`
	} `json:"company"`
}

// StaticSiteListResponse represents the response from the static sites list endpoint.
type StaticSiteListResponse struct {
	Company struct {
		StaticSites struct {
			Items []StaticSiteListItem `json:"items"`
		} `json:"static_sites"`
	} `json:"company"`
}

// SiteListResponse represents the response from the sites list endpoint.
// Based on getSites-Response from the OpenAPI spec.
type SiteListResponse struct {
	Company struct {
		Sites []SiteListItem `json:"sites"`
	} `json:"company"`
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Errors  interface{} `json:"errors,omitempty"`
}

// AuthValidationResponse represents the response from the authentication endpoint.
type AuthValidationResponse struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	ExpiresAt int64  `json:"expires_at"`
	KeyID     string `json:"key_id"`
}

// ResourceType represents the available database resource types.
type ResourceType string

const (
	ResourceTypeDB1 ResourceType = "db1"
	ResourceTypeDB2 ResourceType = "db2"
	ResourceTypeDB3 ResourceType = "db3"
	ResourceTypeDB4 ResourceType = "db4"
	ResourceTypeDB5 ResourceType = "db5"
	ResourceTypeDB6 ResourceType = "db6"
	ResourceTypeDB7 ResourceType = "db7"
	ResourceTypeDB8 ResourceType = "db8"
	ResourceTypeDB9 ResourceType = "db9"
)

// DatabaseType represents the available database types.
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeRedis      DatabaseType = "redis"
	DatabaseTypeMariaDB    DatabaseType = "mariadb"
	DatabaseTypeMySQL      DatabaseType = "mysql"
)

// BuildType represents the available build types for applications.
type BuildType string

const (
	BuildTypeDockerfile BuildType = "dockerfile"
	BuildTypePack       BuildType = "pack"
	BuildTypeNixpacks   BuildType = "nixpacks"
)

// NodeVersion represents the available Node.js versions.
type NodeVersion string

const (
	NodeVersion16 NodeVersion = "16.20.0"
	NodeVersion18 NodeVersion = "18.16.0"
	NodeVersion20 NodeVersion = "20.2.0"
)

// ApplicationStatus represents the possible application states.
type ApplicationStatus string

const (
	ApplicationStatusDeploying ApplicationStatus = "deploying"
	ApplicationStatusDeployed  ApplicationStatus = "deployed"
	ApplicationStatusFailed    ApplicationStatus = "failed"
	ApplicationStatusStopped   ApplicationStatus = "stopped"
)

// DatabaseStatus represents the possible database states.
type DatabaseStatus string

const (
	DatabaseStatusCreating DatabaseStatus = "creating"
	DatabaseStatusActive   DatabaseStatus = "active"
	DatabaseStatusFailed   DatabaseStatus = "failed"
	DatabaseStatusDeleting DatabaseStatus = "deleting"
)

// DeploymentStatus represents the possible deployment states.
type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusRunning    DeploymentStatus = "running"
	DeploymentStatusSuccessful DeploymentStatus = "successful"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusCanceled   DeploymentStatus = "canceled"
)
