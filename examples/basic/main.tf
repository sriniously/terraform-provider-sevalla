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

# Create a database
resource "sevalla_database" "example" {
  name     = "example-db"
  type     = "postgresql"
  version  = "17"
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

output "database_host" {
  value = sevalla_database.example.host
}