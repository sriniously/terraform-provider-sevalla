# Full-Stack Application Example
# This example demonstrates a complete web application stack with:
# - Frontend (React/Vue static site)
# - Backend API (Node.js/Python application)
# - Database (PostgreSQL)
# - Cache (Redis)
# - File storage (Object Storage)
# - CI/CD pipelines

terraform {
  required_providers {
    sevalla = {
      source = "sriniously/sevalla"
    }
  }
}

provider "sevalla" {
  token = var.sevalla_token
}

# PostgreSQL database for application data
resource "sevalla_database" "main_db" {
  name     = "${var.app_name}-postgres"
  type     = "postgresql"
  version  = "14"
  size     = "medium"
  password = var.db_password
}

# Redis cache for session storage and caching
resource "sevalla_database" "cache" {
  name = "${var.app_name}-redis"
  type = "redis"
  size = "small"
}

# Object storage for user uploads and assets
resource "sevalla_object_storage" "uploads" {
  name   = "${var.app_name}-uploads"
  region = var.region
}

# Object storage for static assets (CDN)
resource "sevalla_object_storage" "assets" {
  name   = "${var.app_name}-assets"
  region = var.region
}

# Backend API application
resource "sevalla_application" "api" {
  name        = "${var.app_name}-api"
  description = "Backend API for ${var.app_name}"

  repository = var.api_repo_url

  branch        = var.api_branch
  build_command = "npm ci && npm run build"
  start_command = "npm start"

  environment = {
    NODE_ENV = "production"
    PORT     = "3000"

    # Database connection
    DATABASE_URL = "postgresql://${sevalla_database.main_db.username}:${var.db_password}@${sevalla_database.main_db.host}:${sevalla_database.main_db.port}/${sevalla_database.main_db.name}"

    # Redis connection
    REDIS_URL = "redis://${sevalla_database.cache.host}:${sevalla_database.cache.port}"

    # Object storage configuration
    UPLOADS_BUCKET_ENDPOINT = sevalla_object_storage.uploads.endpoint
    UPLOADS_ACCESS_KEY      = sevalla_object_storage.uploads.access_key
    UPLOADS_SECRET_KEY      = sevalla_object_storage.uploads.secret_key

    ASSETS_BUCKET_ENDPOINT = sevalla_object_storage.assets.endpoint
    ASSETS_ACCESS_KEY      = sevalla_object_storage.assets.access_key
    ASSETS_SECRET_KEY      = sevalla_object_storage.assets.secret_key

    # Application configuration
    JWT_SECRET   = var.jwt_secret
    API_BASE_URL = "https://api.${var.domain}"
    CORS_ORIGIN  = "https://${var.domain}"
  }

  instances = var.api_instances
  memory    = 1024
  cpu       = 500

  domain = "api.${var.domain}"
}

# Frontend static site
resource "sevalla_static_site" "frontend" {
  name = "${var.app_name}-frontend"

  repository = var.frontend_repo_url

  branch    = var.frontend_branch
  build_dir = "dist"
  build_cmd = "npm ci && REACT_APP_API_URL=https://api.${var.domain} npm run build"

  domain = var.domain
}

# CI/CD pipeline for API with automatic deployment
resource "sevalla_pipeline" "api_pipeline" {
  name        = "${var.app_name}-api-pipeline"
  app_id      = sevalla_application.api.id
  branch      = var.api_branch
  auto_deploy = true
}

# CI/CD pipeline for frontend with automatic deployment
resource "sevalla_pipeline" "frontend_pipeline" {
  name        = "${var.app_name}-frontend-pipeline"
  app_id      = sevalla_static_site.frontend.id
  branch      = var.frontend_branch
  auto_deploy = true
}

# Variables
variable "sevalla_token" {
  description = "Sevalla API token"
  type        = string
  sensitive   = true
}

variable "app_name" {
  description = "Application name prefix"
  type        = string
  default     = "myapp"
}

variable "domain" {
  description = "Primary domain for the application"
  type        = string
}

variable "region" {
  description = "Deployment region"
  type        = string
  default     = "us-east-1"
}

variable "api_repo_url" {
  description = "GitHub repository URL for the API"
  type        = string
}

variable "frontend_repo_url" {
  description = "GitHub repository URL for the frontend"
  type        = string
}

variable "api_branch" {
  description = "Git branch for API deployment"
  type        = string
  default     = "main"
}

variable "frontend_branch" {
  description = "Git branch for frontend deployment"
  type        = string
  default     = "main"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
  validation {
    condition     = length(var.db_password) >= 12
    error_message = "Database password must be at least 12 characters long."
  }
}

variable "jwt_secret" {
  description = "JWT secret for authentication"
  type        = string
  sensitive   = true
  validation {
    condition     = length(var.jwt_secret) >= 32
    error_message = "JWT secret must be at least 32 characters long."
  }
}

variable "api_instances" {
  description = "Number of API instances"
  type        = number
  default     = 2
  validation {
    condition     = var.api_instances >= 1 && var.api_instances <= 10
    error_message = "API instances must be between 1 and 10."
  }
}

# Outputs
output "application_urls" {
  description = "URLs for accessing the application"
  value = {
    frontend = "https://${var.domain}"
    api      = "https://api.${var.domain}"
  }
}

output "database_connection" {
  description = "Database connection details"
  value = {
    host     = sevalla_database.main_db.host
    port     = sevalla_database.main_db.port
    username = sevalla_database.main_db.username
    database = sevalla_database.main_db.name
  }
  sensitive = true
}

output "cache_connection" {
  description = "Redis cache connection details"
  value = {
    host = sevalla_database.cache.host
    port = sevalla_database.cache.port
  }
}

output "object_storage" {
  description = "Object storage endpoints and credentials"
  value = {
    uploads = {
      endpoint   = sevalla_object_storage.uploads.endpoint
      access_key = sevalla_object_storage.uploads.access_key
    }
    assets = {
      endpoint   = sevalla_object_storage.assets.endpoint
      access_key = sevalla_object_storage.assets.access_key
    }
  }
  sensitive = true
}

output "resource_ids" {
  description = "Resource IDs for reference"
  value = {
    api_app_id       = sevalla_application.api.id
    frontend_site_id = sevalla_static_site.frontend.id
    database_id      = sevalla_database.main_db.id
    cache_id         = sevalla_database.cache.id
    uploads_bucket   = sevalla_object_storage.uploads.id
    assets_bucket    = sevalla_object_storage.assets.id
  }
}