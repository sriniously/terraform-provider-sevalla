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
					resource.TestCheckResourceAttr("sevalla_application.web_app", "display_name", "fullstack-web-app"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.web_app", "status"),

					// Database checks
					resource.TestCheckResourceAttr("sevalla_database.app_db", "display_name", "fullstack-postgres"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "version", "14"),
					resource.TestCheckResourceAttr("sevalla_database.app_db", "resource_type", "db1"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "name"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "internal_hostname"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_db", "internal_port"),

					// Cache database checks
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "display_name", "fullstack-redis"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "type", "redis"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "version", "7"),
					resource.TestCheckResourceAttr("sevalla_database.app_cache", "resource_type", "db1"),
					resource.TestCheckResourceAttrSet("sevalla_database.app_cache", "id"),

					// Static site checks
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "display_name", "fullstack-frontend"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "default_branch", "main"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "published_directory", "dist"),
					resource.TestCheckResourceAttr("sevalla_static_site.frontend", "build_command", "npm run build"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.frontend", "id"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.frontend", "hostname"),

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
					resource.TestCheckResourceAttr("sevalla_application.api_app", "display_name", "api-with-db"),
					resource.TestCheckResourceAttrSet("sevalla_application.api_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.api_app", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.api_app", "status"),

					// Database checks
					resource.TestCheckResourceAttr("sevalla_database.api_db", "display_name", "api-postgres"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "version", "14"),
					resource.TestCheckResourceAttr("sevalla_database.api_db", "resource_type", "db1"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "name"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "internal_hostname"),
					resource.TestCheckResourceAttrSet("sevalla_database.api_db", "internal_port"),
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
					resource.TestCheckResourceAttr("sevalla_application.dev_app", "display_name", "myapp-dev"),
					resource.TestCheckResourceAttrSet("sevalla_application.dev_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.dev_app", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.dev_app", "status"),

					// Staging environment
					resource.TestCheckResourceAttr("sevalla_application.staging_app", "display_name", "myapp-staging"),
					resource.TestCheckResourceAttrSet("sevalla_application.staging_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.staging_app", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.staging_app", "status"),

					// Production environment
					resource.TestCheckResourceAttr("sevalla_application.prod_app", "display_name", "myapp-production"),
					resource.TestCheckResourceAttrSet("sevalla_application.prod_app", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.prod_app", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.prod_app", "status"),

					// Database checks for each environment
					resource.TestCheckResourceAttr("sevalla_database.dev_db", "display_name", "myapp-dev-db"),
					resource.TestCheckResourceAttr("sevalla_database.dev_db", "resource_type", "db1"),
					resource.TestCheckResourceAttr("sevalla_database.staging_db", "display_name", "myapp-staging-db"),
					resource.TestCheckResourceAttr("sevalla_database.staging_db", "resource_type", "db2"),
					resource.TestCheckResourceAttr("sevalla_database.prod_db", "display_name", "myapp-production-db"),
					resource.TestCheckResourceAttr("sevalla_database.prod_db", "resource_type", "db3"),
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
					resource.TestCheckResourceAttr("sevalla_application.ds_app", "display_name", "datasource-app"),
					resource.TestCheckResourceAttr("sevalla_database.ds_db", "display_name", "datasource-db"),
					resource.TestCheckResourceAttr("sevalla_static_site.ds_site", "display_name", "datasource-site"),
					resource.TestCheckResourceAttr("sevalla_pipeline.ds_pipeline", "name", "datasource-pipeline"),

					// Data source checks - should match resource attributes
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "id",
						"data.sevalla_application.ds_app", "id"),
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "display_name",
						"data.sevalla_application.ds_app", "display_name"),
					resource.TestCheckResourceAttrPair("sevalla_application.ds_app", "status",
						"data.sevalla_application.ds_app", "status"),

					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "id", "data.sevalla_database.ds_db", "id"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "display_name", "data.sevalla_database.ds_db", "display_name"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "type", "data.sevalla_database.ds_db", "type"),
					resource.TestCheckResourceAttrPair("sevalla_database.ds_db", "internal_hostname", "data.sevalla_database.ds_db", "internal_hostname"),

					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "id",
						"data.sevalla_static_site.ds_site", "id"),
					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "display_name",
						"data.sevalla_static_site.ds_site", "display_name"),
					resource.TestCheckResourceAttrPair("sevalla_static_site.ds_site", "hostname",
						"data.sevalla_static_site.ds_site", "hostname"),

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
  display_name    = "fullstack-postgres"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db1"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "test-password"
  db_user         = "testuser"
}

# Redis cache for sessions
resource "sevalla_database" "app_cache" {
  display_name    = "fullstack-redis"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db1"
  type            = "redis"
  version         = "7"
  db_name         = "redis"
  db_password     = "test-password"
}


# Backend API application
resource "sevalla_application" "web_app" {
  display_name   = "fullstack-web-app"
  company_id     = "` + testAccCompanyID() + `"
  repo_url       = "https://github.com/example/fullstack-backend"
  default_branch = "main"
  auto_deploy    = true
}

# Frontend static site
resource "sevalla_static_site" "frontend" {
  display_name     = "fullstack-frontend"
  company_id       = "` + testAccCompanyID() + `"
  repo_url         = "https://github.com/example/fullstack-frontend"
  default_branch   = "main"
  auto_deploy      = true
  build_command    = "npm run build"
  published_directory = "dist"
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
  display_name    = "api-postgres"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db1"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "test-password"
  db_user         = "testuser"
}

# API application
resource "sevalla_application" "api_app" {
  display_name   = "api-with-db"
  company_id     = "` + testAccCompanyID() + `"
  repo_url       = "https://github.com/example/api-app"
  default_branch = "main"
  auto_deploy    = true
}
`
}

func testAccIntegrationMultiEnvironmentConfig() string {
	return providerConfig + `
# Development environment
resource "sevalla_database" "dev_db" {
  display_name    = "myapp-dev-db"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db1"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "dev-password"
  db_user         = "testuser"
}

resource "sevalla_application" "dev_app" {
  display_name   = "myapp-dev"
  company_id     = "` + testAccCompanyID() + `"
  repo_url       = "https://github.com/example/myapp"
  default_branch = "develop"
  auto_deploy    = true
}

# Staging environment
resource "sevalla_database" "staging_db" {
  display_name    = "myapp-staging-db"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db2"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "staging-password"
  db_user         = "testuser"
}

resource "sevalla_application" "staging_app" {
  display_name   = "myapp-staging"
  company_id     = "` + testAccCompanyID() + `"
  repo_url       = "https://github.com/example/myapp"
  default_branch = "staging"
  auto_deploy    = true
}

# Production environment
resource "sevalla_database" "prod_db" {
  display_name    = "myapp-production-db"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db3"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "production-password"
  db_user         = "testuser"
}

resource "sevalla_application" "prod_app" {
  display_name   = "myapp-production"
  company_id     = "` + testAccCompanyID() + `"
  repo_url       = "https://github.com/example/myapp"
  default_branch = "main"
  auto_deploy    = true
}
`
}

func testAccIntegrationDataSourcesConfig() string {
	return providerConfig + `
# Create resources
resource "sevalla_application" "ds_app" {
  display_name  = "datasource-app"
  company_id    = "` + testAccCompanyID() + `"
  repo_url      = "https://github.com/example/datasource-app"
  auto_deploy   = true
}

resource "sevalla_database" "ds_db" {
  display_name    = "datasource-db"
  company_id      = "` + testAccCompanyID() + `"
  location        = "us-central1"
  resource_type   = "db1"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "test-password"
  db_user         = "testuser"
}

resource "sevalla_static_site" "ds_site" {
  display_name     = "datasource-site"
  company_id       = "` + testAccCompanyID() + `"
  repo_url         = "https://github.com/example/datasource-site"
  default_branch   = "main"
  auto_deploy      = true
  build_command    = "npm run build"
  published_directory = "dist"
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


data "sevalla_pipeline" "ds_pipeline" {
  id = sevalla_pipeline.ds_pipeline.id
}
`
}
