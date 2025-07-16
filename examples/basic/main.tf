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

# Create a simple application
resource "sevalla_application" "example" {
  name        = "example-app"
  description = "Example application"
  
  repository {
    url  = "https://github.com/user/example-app"
    type = "github"
  }
  
  branch        = "main"
  build_command = "npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV = "production"
    PORT     = "3000"
  }
  
  instances = 1
  memory    = 512
  cpu       = 250
}

# Create a database
resource "sevalla_database" "example" {
  name     = "example-db"
  type     = "postgresql"
  version  = "14"
  size     = "small"
  password = var.db_password
}

variable "sevalla_token" {
  description = "Sevalla API token"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

output "app_id" {
  value = sevalla_application.example.id
}

output "database_host" {
  value = sevalla_database.example.host
}