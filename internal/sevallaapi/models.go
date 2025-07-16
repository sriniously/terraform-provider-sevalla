package sevallaapi

import "time"

// Application represents a Sevalla application.
type Application struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Domain       string            `json:"domain,omitempty"`
	Repository   *Repository       `json:"repository,omitempty"`
	Branch       string            `json:"branch,omitempty"`
	BuildCommand string            `json:"build_command,omitempty"`
	StartCommand string            `json:"start_command,omitempty"`
	Environment  map[string]string `json:"environment,omitempty"`
	Instances    int               `json:"instances,omitempty"`
	Memory       int               `json:"memory,omitempty"`
	CPU          int               `json:"cpu,omitempty"`
	Status       string            `json:"status,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// Repository represents a source code repository.
type Repository struct {
	URL    string `json:"url"`
	Type   string `json:"type"` // github, gitlab, bitbucket
	Branch string `json:"branch,omitempty"`
}

// CreateApplicationRequest represents the request to create an application.
type CreateApplicationRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Repository   *Repository       `json:"repository,omitempty"`
	Branch       string            `json:"branch,omitempty"`
	BuildCommand string            `json:"build_command,omitempty"`
	StartCommand string            `json:"start_command,omitempty"`
	Environment  map[string]string `json:"environment,omitempty"`
	Instances    int               `json:"instances,omitempty"`
	Memory       int               `json:"memory,omitempty"`
	CPU          int               `json:"cpu,omitempty"`
}

// UpdateApplicationRequest represents the request to update an application.
type UpdateApplicationRequest struct {
	Name         *string            `json:"name,omitempty"`
	Description  *string            `json:"description,omitempty"`
	Repository   *Repository        `json:"repository,omitempty"`
	Branch       *string            `json:"branch,omitempty"`
	BuildCommand *string            `json:"build_command,omitempty"`
	StartCommand *string            `json:"start_command,omitempty"`
	Environment  *map[string]string `json:"environment,omitempty"`
	Instances    *int               `json:"instances,omitempty"`
	Memory       *int               `json:"memory,omitempty"`
	CPU          *int               `json:"cpu,omitempty"`
}

// Database represents a Sevalla database.
type Database struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // postgresql, mysql, mariadb, redis
	Version   string    `json:"version"`
	Size      string    `json:"size"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateDatabaseRequest represents the request to create a database.
type CreateDatabaseRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // postgresql, mysql, mariadb, redis
	Version  string `json:"version,omitempty"`
	Size     string `json:"size,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateDatabaseRequest represents the request to update a database.
type UpdateDatabaseRequest struct {
	Name     *string `json:"name,omitempty"`
	Size     *string `json:"size,omitempty"`
	Password *string `json:"password,omitempty"`
}

// StaticSite represents a Sevalla static site.
type StaticSite struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Domain     string      `json:"domain,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
	Branch     string      `json:"branch,omitempty"`
	BuildDir   string      `json:"build_dir,omitempty"`
	BuildCmd   string      `json:"build_cmd,omitempty"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// CreateStaticSiteRequest represents the request to create a static site.
type CreateStaticSiteRequest struct {
	Name       string      `json:"name"`
	Repository *Repository `json:"repository,omitempty"`
	Branch     string      `json:"branch,omitempty"`
	BuildDir   string      `json:"build_dir,omitempty"`
	BuildCmd   string      `json:"build_cmd,omitempty"`
}

// UpdateStaticSiteRequest represents the request to update a static site.
type UpdateStaticSiteRequest struct {
	Name       *string     `json:"name,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
	Branch     *string     `json:"branch,omitempty"`
	BuildDir   *string     `json:"build_dir,omitempty"`
	BuildCmd   *string     `json:"build_cmd,omitempty"`
}

// ObjectStorage represents a Sevalla object storage bucket.
type ObjectStorage struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	Size      int64     `json:"size"`
	Objects   int       `json:"objects"`
	Endpoint  string    `json:"endpoint"`
	AccessKey string    `json:"access_key"`
	SecretKey string    `json:"secret_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateObjectStorageRequest represents the request to create object storage.
type CreateObjectStorageRequest struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
}

// UpdateObjectStorageRequest represents the request to update object storage.
type UpdateObjectStorageRequest struct {
	Name *string `json:"name,omitempty"`
}

// Deployment represents a deployment.
type Deployment struct {
	ID         string    `json:"id"`
	AppID      string    `json:"app_id"`
	Status     string    `json:"status"`
	Branch     string    `json:"branch"`
	CommitHash string    `json:"commit_hash"`
	CommitMsg  string    `json:"commit_message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Pipeline represents a deployment pipeline.
type Pipeline struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	AppID      string    `json:"app_id"`
	Branch     string    `json:"branch"`
	AutoDeploy bool      `json:"auto_deploy"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreatePipelineRequest represents the request to create a pipeline.
type CreatePipelineRequest struct {
	Name       string `json:"name"`
	AppID      string `json:"app_id"`
	Branch     string `json:"branch"`
	AutoDeploy bool   `json:"auto_deploy"`
}

// UpdatePipelineRequest represents the request to update a pipeline.
type UpdatePipelineRequest struct {
	Name       *string `json:"name,omitempty"`
	Branch     *string `json:"branch,omitempty"`
	AutoDeploy *bool   `json:"auto_deploy,omitempty"`
}
