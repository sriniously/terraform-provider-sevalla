package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIntegrationFullStack tests a full-stack application with all resources.
func TestAccIntegrationFullStack(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationFullStackConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Application checks
					resource.TestCheckResourceAttr("sevalla_application.web_app", "name", "fullstack-web-app"),
					resource.TestCheckResourceAttr("sevalla_application.web_app", "instances", "2"),
					resource.TestCheckResourceAttr("sevalla_application.web_app", "memory", "1024"),
					resource.TestCheckResourceAttr("sevalla_application.web_app", "cpu", "500"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "domain"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "status"),

					// Database checks
					resource.TestCheckResourceAttr("sevalla_database.app_db", "name", "fullstack-postgres"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "version", "14"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "size", "medium"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "host"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "port"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "username"),

					// Cache database checks
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "name", "fullstack-redis"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "type", "redis"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "version", "7"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "size", "small"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_cache", "id"),

					// Static site checks
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "name", "fullstack-frontend"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "branch", "main"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "build_dir", "dist"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "build_cmd", "npm run build"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.frontend", "id"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.frontend", "domain"),

					// Object storage checks
					resource.TestCheckResourceAttr("sevalla_object_storage.app_storage", "name", "fullstack-storage"),
					resource.TestCheckResourceAttr("sevalla_object_storage.app_storage", "region", "us-east-1"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.app_storage", "id"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.app_storage", "endpoint"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.app_storage", "access_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.app_storage", "secret_key"),

					// Pipeline checks
					resource.TestCheckResourceAttr("sevalla_pipeline.app_pipeline", "name", "fullstack-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.app_pipeline", "branch", "main"),
					resource.TestCheckResourceAttr("sevalla_pipeline.app_pipeline", "auto_deploy", "true"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.app_pipeline", "id"),
					// Check that pipeline references the application
					resource.TestCheckResourceAttrPair("sevalla_pipeline.app_pipeline", "app_id", "sevalla_application.web_app", "id"),
				),
			},
		},
	})
}

// TestAccIntegrationAppWithDatabase tests an application with a database.
func TestAccIntegrationAppWithDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAppWithDatabaseConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Application checks
					resource.TestCheckResourceAttr("sevalla_application.api_app", "name", "api-with-db"),
					resource.TestCheckResourceAttr("sevalla_application.api_app", "instances", "1"),
					resource.TestCheckResourceAttr("sevalla_application.api_app", "memory", "512"),
					resource.TestCheckResourceAttr("sevalla_application.api_app", "cpu", "250"),
					resource.TestCheckResourceAttrSet("sevalla_application.api_app", "id"),

					// Database checks
					resource.TestCheckResourceAttr("sevalla_database.api_db", "name", "api-postgres"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "version", "14"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "size", "small"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "host"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "port"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "username"),

					// Check that environment variables are set correctly
					resource.TestCheckResourceAttr("sevalla_application.api_app", "environment.NODE_ENV", "production"),
					resource.TestCheckResourceAttr("sevalla_application.api_app", "environment.PORT", "3000"),
					// DATABASE_URL should be constructed from database attributes
					resource.TestCheckResourceAttrSet("sevalla_application.api_app", "environment.DATABASE_URL"),
				),
			},
		},
	})
}

// TestAccIntegrationMultiEnvironment tests multiple environments (dev, staging, prod).
func TestAccIntegrationMultiEnvironment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationMultiEnvironmentConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Development environment
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "name", "myapp-dev"),
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "instances", "1"),
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "memory", "512"),
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "environment.ENVIRONMENT", "development"),
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "environment.LOG_LEVEL", "debug"),
					resource.TestCheckResourceAttrSet("sevalla_application.dev_app", "id"),

					// Staging environment
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "name", "myapp-staging"),
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "instances", "2"),
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "memory", "1024"),
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "environment.ENVIRONMENT", "staging"),
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "environment.LOG_LEVEL", "info"),
					resource.TestCheckResourceAttrSet("sevalla_application.staging_app", "id"),

					// Production environment
					resource.TestCheckResourceAttr("sevalla_application.prod_app", "name", "myapp-production"),
					resource.TestCheckResourceAttr("sevalla_application.prod_app", "instances", "3"),
					resource.TestCheckResourceAttr("sevalla_application.prod_app", "memory", "2048"),
					resource.TestCheckResourceAttr("sevalla_application.prod_app",
						"environment.ENVIRONMENT", "production"),
					resource.TestCheckResourceAttr("sevalla_application.prod_app", "environment.LOG_LEVEL", "info"),
					resource.TestCheckResourceAttrSet("sevalla_application.prod_app", "id"),

					// Database checks for each environment
					resource.TestCheckResourceAttr("sevalla_database.dev_db", "name", "myapp-dev-db"),
					resource.TestCheckResourceAttr("sevalla_database.dev_db", "size", "small"),
					resource.TestCheckResourceAttr("sevalla_database.staging_db", "name", "myapp-staging-db"),
					resource.TestCheckResourceAttr("sevalla_database.staging_db", "size", "medium"),
					resource.TestCheckResourceAttr("sevalla_database.prod_db", "name", "myapp-production-db"),
					resource.TestCheckResourceAttr("sevalla_database.prod_db", "size", "large"),
				),
			},
		},
	})
}

// TestAccIntegrationDataSourcesWithResources tests data sources with actual resources.
func TestAccIntegrationDataSourcesWithResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationDataSourcesConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Resource checks
					resource.TestCheckResourceAttr("sevalla_application.ds_app", "name", "datasource-app"),
					resource.TestCheckResourceAttr("sevalla_database.ds_db", "name", "datasource-db"),
					resource.TestCheckResourceAttr("sevalla_static_site.ds_site", "name", "datasource-site"),
					resource.TestCheckResourceAttr("sevalla_object_storage.ds_storage", "name", "datasource-storage"),
					resource.TestCheckResourceAttr("sevalla_pipeline.ds_pipeline", "name", "datasource-pipeline"),

					// Data source checks - should match resource attributes
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "id",
						"data.sevalla_application.ds_app", "id"),
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "name",
						"data.sevalla_application.ds_app", "name"),
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "status",
						"data.sevalla_application.ds_app", "status"),

					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "id", "data.sevalla_database.ds_db", "id"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "name", "data.sevalla_database.ds_db", "name"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "type", "data.sevalla_database.ds_db", "type"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "host", "data.sevalla_database.ds_db", "host"),

					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "id",
						"data.sevalla_static_site.ds_site", "id"),
					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "name",
						"data.sevalla_static_site.ds_site", "name"),
					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "domain",
						"data.sevalla_static_site.ds_site", "domain"),

					resource.TestCheckResourceAttrPair("sevalla_object_storage.ds_storage", "id",
						"data.sevalla_object_storage.ds_storage", "id"),
					resource.TestCheckResourceAttrPair("sevalla_object_storage.ds_storage", "name",
						"data.sevalla_object_storage.ds_storage", "name"),
					resource.TestCheckResourceAttrPair("sevalla_object_storage.ds_storage", "endpoint",
						"data.sevalla_object_storage.ds_storage", "endpoint"),

					resource.TestCheckResourceAttrPair("sevalla_pipeline.ds_pipeline", "id",
						"data.sevalla_pipeline.ds_pipeline", "id"),
					resource.TestCheckResourceAttrPair("sevalla_pipeline.ds_pipeline", "name",
						"data.sevalla_pipeline.ds_pipeline", "name"),
					resource.TestCheckResourceAttrPair("sevalla_pipeline.ds_pipeline", "app_id",
						"data.sevalla_pipeline.ds_pipeline", "app_id"),
					resource.TestCheckResourceAttrPair("sevalla_pipeline.ds_pipeline", "branch",
						"data.sevalla_pipeline.ds_pipeline", "branch"),
					resource.TestCheckResourceAttrPair("sevalla_pipeline.ds_pipeline", "auto_deploy",
						"data.sevalla_pipeline.ds_pipeline", "auto_deploy"),
				),
			},
		},
	})
}

// Configuration functions for integration tests

func testAccIntegrationFullStackConfig() string {
	return providerConfig + `
# PostgreSQL database for the application
resource "sevalla_database" "app_db" {
  name     = "fullstack-postgres"
  type     = "postgresql"
  version  = "14"
  size     = "medium"
  password = "test-password"
}

# Redis cache for sessions
resource "sevalla_database" "app_cache" {
  name     = "fullstack-redis"
  type     = "redis"
  version  = "7"
  size     = "small"
  password = "test-password"
}

# Object storage for uploads
resource "sevalla_object_storage" "app_storage" {
  name   = "fullstack-storage"
  region = "us-east-1"
}

# Backend API application
resource "sevalla_application" "web_app" {
  name        = "fullstack-web-app"
  description = "Full-stack web application backend"
  
  repository {
    url    = "https://github.com/example/fullstack-backend"
    type   = "github"
    branch = "main"
  }
  
  branch        = "main"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV     = "production"
    PORT         = "3000"
    DATABASE_URL = "postgresql://${sevalla_database.app_db.username}:${sevalla_database.app_db.password}@${sevalla_database.app_db.host}:${sevalla_database.app_db.port}/${sevalla_database.app_db.name}"
    REDIS_URL    = "redis://:${sevalla_database.app_cache.password}@${sevalla_database.app_cache.host}:${sevalla_database.app_cache.port}"
    S3_BUCKET    = sevalla_object_storage.app_storage.name
    S3_ENDPOINT  = sevalla_object_storage.app_storage.endpoint
    S3_ACCESS_KEY = sevalla_object_storage.app_storage.access_key
    S3_SECRET_KEY = sevalla_object_storage.app_storage.secret_key
  }
  
  instances = 2
  memory    = 1024
  cpu       = 500
}

# Frontend static site
resource "sevalla_static_site" "frontend" {
  name      = "fullstack-frontend"
  branch    = "main"
  build_dir = "dist"
  build_cmd = "npm run build"
  
  repository {
    url    = "https://github.com/example/fullstack-frontend"
    type   = "github"
    branch = "main"
  }
}

# CI/CD pipeline
resource "sevalla_pipeline" "app_pipeline" {
  name        = "fullstack-pipeline"
  app_id      = sevalla_application.web_app.id
  branch      = "main"
  auto_deploy = true
}
`
}

func testAccIntegrationAppWithDatabaseConfig() string {
	return providerConfig + `
# PostgreSQL database
resource "sevalla_database" "api_db" {
  name     = "api-postgres"
  type     = "postgresql"
  version  = "14"
  size     = "small"
  password = "test-password"
}

# API application
resource "sevalla_application" "api_app" {
  name        = "api-with-db"
  description = "API application with PostgreSQL database"
  
  repository {
    url    = "https://github.com/example/api-app"
    type   = "github"
    branch = "main"
  }
  
  branch        = "main"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV     = "production"
    PORT         = "3000"
    DATABASE_URL = "postgresql://${sevalla_database.api_db.username}:${sevalla_database.api_db.password}@${sevalla_database.api_db.host}:${sevalla_database.api_db.port}/${sevalla_database.api_db.name}"
  }
  
  instances = 1
  memory    = 512
  cpu       = 250
}
`
}

func testAccIntegrationMultiEnvironmentConfig() string {
	return providerConfig + `
# Development environment
resource "sevalla_database" "dev_db" {
  name     = "myapp-dev-db"
  type     = "postgresql"
  version  = "14"
  size     = "small"
  password = "dev-password"
}

resource "sevalla_application" "dev_app" {
  name        = "myapp-dev"
  description = "Development application"
  
  repository {
    url    = "https://github.com/example/myapp"
    type   = "github"
    branch = "develop"
  }
  
  branch        = "develop"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    ENVIRONMENT  = "development"
    LOG_LEVEL    = "debug"
    PORT         = "3000"
    DATABASE_URL = "postgresql://${sevalla_database.dev_db.username}:${sevalla_database.dev_db.password}@${sevalla_database.dev_db.host}:${sevalla_database.dev_db.port}/${sevalla_database.dev_db.name}"
  }
  
  instances = 1
  memory    = 512
  cpu       = 250
}

# Staging environment
resource "sevalla_database" "staging_db" {
  name     = "myapp-staging-db"
  type     = "postgresql"
  version  = "14"
  size     = "medium"
  password = "staging-password"
}

resource "sevalla_application" "staging_app" {
  name        = "myapp-staging"
  description = "Staging application"
  
  repository {
    url    = "https://github.com/example/myapp"
    type   = "github"
    branch = "staging"
  }
  
  branch        = "staging"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    ENVIRONMENT  = "staging"
    LOG_LEVEL    = "info"
    PORT         = "3000"
    DATABASE_URL = "postgresql://${sevalla_database.staging_db.username}:${sevalla_database.staging_db.password}@${sevalla_database.staging_db.host}:${sevalla_database.staging_db.port}/${sevalla_database.staging_db.name}"
  }
  
  instances = 2
  memory    = 1024
  cpu       = 500
}

# Production environment
resource "sevalla_database" "prod_db" {
  name     = "myapp-production-db"
  type     = "postgresql"
  version  = "14"
  size     = "large"
  password = "production-password"
}

resource "sevalla_application" "prod_app" {
  name        = "myapp-production"
  description = "Production application"
  
  repository {
    url    = "https://github.com/example/myapp"
    type   = "github"
    branch = "main"
  }
  
  branch        = "main"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    ENVIRONMENT  = "production"
    LOG_LEVEL    = "info"
    PORT         = "3000"
    DATABASE_URL = "postgresql://${sevalla_database.prod_db.username}:${sevalla_database.prod_db.password}@${sevalla_database.prod_db.host}:${sevalla_database.prod_db.port}/${sevalla_database.prod_db.name}"
  }
  
  instances = 3
  memory    = 2048
  cpu       = 1000
}
`
}

func testAccIntegrationDataSourcesConfig() string {
	return providerConfig + `
# Create resources
resource "sevalla_application" "ds_app" {
  name        = "datasource-app"
  description = "Application for data source testing"
  
  repository {
    url    = "https://github.com/example/datasource-app"
    type   = "github"
    branch = "main"
  }
  
  branch        = "main"
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV = "production"
    PORT     = "3000"
  }
  
  instances = 1
  memory    = 512
  cpu       = 250
}

resource "sevalla_database" "ds_db" {
  name     = "datasource-db"
  type     = "postgresql"
  version  = "14"
  size     = "small"
  password = "test-password"
}

resource "sevalla_static_site" "ds_site" {
  name      = "datasource-site"
  branch    = "main"
  build_dir = "dist"
  build_cmd = "npm run build"
  
  repository {
    url    = "https://github.com/example/datasource-site"
    type   = "github"
    branch = "main"
  }
}

resource "sevalla_object_storage" "ds_storage" {
  name   = "datasource-storage"
  region = "us-east-1"
}

resource "sevalla_pipeline" "ds_pipeline" {
  name        = "datasource-pipeline"
  app_id      = sevalla_application.ds_app.id
  branch      = "main"
  auto_deploy = true
}

# Data sources
data "sevalla_application" "ds_app" {
  id = sevalla_application.ds_app.id
}

data "sevalla_database" "ds_db" {
  id = sevalla_database.ds_db.id
}

data "sevalla_static_site" "ds_site" {
  id = sevalla_static_site.ds_site.id
}

data "sevalla_object_storage" "ds_storage" {
  id = sevalla_object_storage.ds_storage.id
}

data "sevalla_pipeline" "ds_pipeline" {
  id = sevalla_pipeline.ds_pipeline.id
}
`
}
