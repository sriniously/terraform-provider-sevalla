terraform {
  required_providers {
    sevalla = {
      source  = "sriniously/sevalla"
      version = "~> 1.0"
    }
  }
}

provider "sevalla" {
  # Configuration options
  api_token = var.sevalla_api_token
  base_url  = "https://api.sevalla.com" # Optional: defaults to https://api.sevalla.com
}

# Example: Create an application
resource "sevalla_application" "example" {
  name          = "my-app"
  github_repo   = "https://github.com/user/repo"
  branch        = "main"
  build_command = "npm run build"
  environment_variables = {
    NODE_ENV = "production"
  }
}

# Example: Create a database
resource "sevalla_database" "example" {
  name    = "my-database"
  type    = "postgresql"
  version = "15"
  plan    = "starter"
  region  = "us-east-1"
}