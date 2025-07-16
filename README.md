# Terraform Provider for Sevalla (Community)

The Sevalla Terraform provider allows you to manage Sevalla cloud resources using Infrastructure as Code (IaC). This provider supports creating and managing applications, databases, static sites, object storage, and deployment pipelines.

**Note:** This is a community-maintained provider, not an official Sevalla provider. It is maintained by [@sriniously](https://github.com/sriniously).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)
- A Sevalla account with API access

## Installation

### Terraform Registry (Recommended)

The Sevalla provider is available on the [Terraform Registry](https://registry.terraform.io/providers/sriniously/sevalla/latest). To use it, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    sevalla = {
      source  = "sriniously/sevalla"
      version = "~> 0.1.0"
    }
  }
}
```

### Manual Installation (Development)

#### Option 1: Install from Source

1. Clone the repository:
```bash
git clone https://github.com/sriniously/terraform-provider-sevalla.git
cd terraform-provider-sevalla
```

2. Build and install the provider:
```bash
go install
```

#### Option 2: Development Override

Create a `.terraformrc` file in your home directory:

```hcl
provider_installation {
  dev_overrides {
    "sriniously/sevalla" = "/Users/YOUR_USERNAME/go/bin"
  }
  
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Getting Started

### Step 1: Obtain API Credentials

1. Log in to your Sevalla account
2. Navigate to Account Settings → API Tokens
3. Create a new API token with appropriate permissions
4. Save the token securely

### Step 2: Configure Authentication

Set your API token as an environment variable:

```bash
export SEVALLA_TOKEN="your-api-token-here"
```

Or configure it in your Terraform configuration (not recommended for production):

```hcl
provider "sevalla" {
  token = "your-api-token-here"
}
```

### Step 3: Create Your First Configuration

Create a `main.tf` file:

```hcl
terraform {
  required_providers {
    sevalla = {
      source = "sriniously/sevalla"
      version = "~> 0.1.0"
    }
  }
}

provider "sevalla" {
  # Token will be read from SEVALLA_TOKEN environment variable
}

# Create a simple Node.js application
resource "sevalla_application" "my_app" {
  name        = "hello-world-app"
  description = "My first Sevalla app via Terraform"
  
  repository {
    url    = "https://github.com/your-username/your-app"
    type   = "github"
    branch = "main"
  }
  
  build_command = "npm install && npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV = "production"
    PORT     = "3000"
  }
  
  instances = 1
  memory    = 512  # MB
  cpu       = 250  # millicores
}

# Output the application URL
output "app_url" {
  value = sevalla_application.my_app.domain
}
```

### Step 4: Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Preview the changes
terraform plan

# Apply the configuration
terraform apply
```

## Available Examples

We've created practical examples demonstrating real-world use cases:

### 1. [Basic Example](examples/basic/)
Simple application with database setup:
- **Application**: Node.js application with GitHub repository
- **Database**: PostgreSQL database with connection configuration
- **Configuration**: Environment variables and resource allocation

### 2. [Full-Stack Application](examples/full-stack-app/)
Complete web application with frontend, backend, database, cache, and storage:
- **Frontend**: Static site with custom domain
- **Backend**: API application with environment configuration
- **Database**: PostgreSQL for persistent data
- **Cache**: Redis for sessions and caching
- **Storage**: Object storage for uploads and assets
- **CI/CD**: Automated deployment pipelines

### 3. [Provider Configuration](examples/provider/)
Basic provider setup examples:
- **Authentication**: API token configuration
- **Basic Resources**: Simple resource creation examples

## Quick Start Examples

### Simple Web Application
```hcl
# PostgreSQL database
resource "sevalla_database" "app_db" {
  name     = "myapp-postgres"
  type     = "postgresql"
  version  = "14"
  size     = "medium"
  password = var.db_password
}

# Application with database connection
resource "sevalla_application" "web_app" {
  name = "my-web-application"
  
  repository {
    url  = "https://github.com/mycompany/webapp"
    type = "github"
  }
  
  branch        = "main"
  build_command = "npm ci && npm run build"
  start_command = "npm start"
  
  environment = {
    NODE_ENV     = "production"
    DATABASE_URL = "postgresql://${sevalla_database.app_db.username}:${var.db_password}@${sevalla_database.app_db.host}:${sevalla_database.app_db.port}/${sevalla_database.app_db.name}"
  }
  
  instances = 2
  memory    = 1024
  cpu       = 500
}
```

### Static Website
```hcl
resource "sevalla_static_site" "marketing" {
  name = "company-website"
  
  repository {
    url  = "https://github.com/mycompany/website"
    type = "github"
  }
  
  branch    = "main"
  build_dir = "dist"
  build_cmd = "npm ci && npm run build"
}
```

### 4. Multi-Environment Setup

```hcl
# Use Terraform workspaces for environments
locals {
  env = terraform.workspace
  
  instance_counts = {
    dev     = 1
    staging = 2
    prod    = 3
  }
  
  memory_sizes = {
    dev     = 512
    staging = 1024
    prod    = 2048
  }
}

resource "sevalla_application" "app" {
  name = "myapp-${local.env}"
  
  repository {
    url    = "https://github.com/mycompany/app"
    type   = "github"
    branch = local.env == "prod" ? "main" : local.env
  }
  
  environment = {
    ENVIRONMENT = local.env
    LOG_LEVEL   = local.env == "prod" ? "info" : "debug"
  }
  
  instances = local.instance_counts[local.env]
  memory    = local.memory_sizes[local.env]
  cpu       = 500
}
```

## Managing State

### Remote State Storage

For production use, store your Terraform state remotely:

```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    key    = "sevalla/terraform.tfstate"
    region = "us-east-1"
  }
}
```

### State Locking

Use DynamoDB for state locking when using S3 backend:

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "sevalla/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-state-lock"
  }
}
```

## Best Practices

### 1. Use Variables for Configuration

```hcl
variable "app_name" {
  description = "Application name"
  type        = string
}

variable "instance_count" {
  description = "Number of application instances"
  type        = number
  default     = 1
}

variable "environment_variables" {
  description = "Application environment variables"
  type        = map(string)
  default     = {}
}

resource "sevalla_application" "app" {
  name        = var.app_name
  instances   = var.instance_count
  environment = var.environment_variables
  # ... other configuration
}
```

### 2. Use Outputs for Important Information

```hcl
output "application_url" {
  description = "The URL of the deployed application"
  value       = sevalla_application.app.domain
}

output "database_connection" {
  description = "Database connection information"
  value = {
    host     = sevalla_database.db.host
    port     = sevalla_database.db.port
    username = sevalla_database.db.username
  }
  sensitive = true
}
```

### 3. Organize with Modules

Create reusable modules for common patterns:

```hcl
# modules/web-app/main.tf
variable "name" {
  type = string
}

variable "repo_url" {
  type = string
}

variable "with_database" {
  type    = bool
  default = false
}

resource "sevalla_application" "app" {
  name = var.name
  # ... configuration
}

resource "sevalla_database" "db" {
  count = var.with_database ? 1 : 0
  name  = "${var.name}-db"
  # ... configuration
}

# Use the module
module "my_app" {
  source = "./modules/web-app"
  
  name          = "my-application"
  repo_url      = "https://github.com/user/repo"
  with_database = true
}
```

### 4. Handle Sensitive Data Properly

```hcl
# terraform.tfvars (add to .gitignore)
db_password = "actual-password-here"

# variables.tf
variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

# Or use environment variables
# export TF_VAR_db_password="actual-password-here"
```

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   ```
   Error: Unable to find token
   ```
   Solution: Ensure `SEVALLA_TOKEN` is set or token is provided in the provider configuration.

2. **Resource Not Found**
   ```
   Error: HTTP 404: Resource not found
   ```
   Solution: Check if the resource exists in Sevalla and you have appropriate permissions.

3. **Rate Limiting**
   ```
   Error: HTTP 429: Too Many Requests
   ```
   Solution: Add delays between resource creation or use `-parallelism=1` flag.

### Debugging

Enable debug logging:
```bash
export TF_LOG=DEBUG
terraform apply
```

### Import Existing Resources

Import resources that were created outside of Terraform:

```bash
# Import an application
terraform import sevalla_application.app app-12345

# Import a database
terraform import sevalla_database.db db-67890

# Import a static site
terraform import sevalla_static_site.site site-abcde
```

## Migration Guide

### From Manual Configuration to Terraform

1. List existing resources in Sevalla dashboard
2. Create Terraform configurations matching existing resources
3. Import each resource using `terraform import`
4. Run `terraform plan` to verify no changes
5. Make adjustments as needed

### From Other Providers

If migrating from another cloud provider:

1. Export data from the current provider
2. Create equivalent Sevalla resources
3. Update application configurations
4. Test thoroughly before switching DNS

## Complete Provider Reference

### Supported Resources

The provider currently supports the following resources:

1. **sevalla_application** - Manages applications with repository integration, build configuration, and environment variables
2. **sevalla_database** - Manages databases (PostgreSQL, MySQL, MariaDB, Redis)
3. **sevalla_static_site** - Manages static websites with build configuration
4. **sevalla_object_storage** - Manages object storage buckets
5. **sevalla_pipeline** - Manages CI/CD deployment pipelines

### Supported Data Sources

The provider includes data sources for fetching existing resources:

1. **sevalla_application** - Fetches existing application details
2. **sevalla_database** - Fetches existing database details
3. **sevalla_static_site** - Fetches existing static site details
4. **sevalla_object_storage** - Fetches existing object storage details

### Provider Configuration

```hcl
provider "sevalla" {
  # Required - Sevalla API token
  # Can also be set via SEVALLA_TOKEN environment variable
  token = "your-api-token"
  
  # Optional - API base URL (defaults to https://api.sevalla.com)
  # Useful for testing or private Sevalla installations
  base_url = "https://api.sevalla.com"
}
```

### Environment Variables

The provider supports the following environment variables:

- `SEVALLA_TOKEN` - Your Sevalla API token (recommended for security)
- `TF_LOG` - Set to `DEBUG` for detailed logging

## CI/CD Integration

### GitHub Actions

```yaml
name: Deploy to Sevalla

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v2
      with:
        terraform_version: 1.5.0
    
    - name: Terraform Init
      run: terraform init
    
    - name: Terraform Plan
      run: terraform plan
      env:
        SEVALLA_TOKEN: ${{ secrets.SEVALLA_TOKEN }}
    
    - name: Terraform Apply
      if: github.ref == 'refs/heads/main'
      run: terraform apply -auto-approve
      env:
        SEVALLA_TOKEN: ${{ secrets.SEVALLA_TOKEN }}
```

### GitLab CI

```yaml
stages:
  - validate
  - plan
  - apply

variables:
  TF_ROOT: ${CI_PROJECT_DIR}
  TF_IN_AUTOMATION: "true"

cache:
  key: "${TF_ROOT}"
  paths:
    - ${TF_ROOT}/.terraform

before_script:
  - terraform --version
  - terraform init

validate:
  stage: validate
  script:
    - terraform validate

plan:
  stage: plan
  script:
    - terraform plan -out=plan.tfplan
  artifacts:
    paths:
      - plan.tfplan

apply:
  stage: apply
  script:
    - terraform apply -auto-approve plan.tfplan
  dependencies:
    - plan
  only:
    - main
```

## Advanced Configuration Examples

### Blue-Green Deployments

```hcl
locals {
  # Toggle between blue and green
  active_color = "blue"  # or "green"
  
  colors = {
    blue  = "blue"
    green = "green"
  }
}

# Blue environment
resource "sevalla_application" "app_blue" {
  name = "myapp-blue"
  
  repository {
    url  = "https://github.com/company/app"
    type = "github"
  }
  
  branch    = "release/blue"
  instances = local.active_color == "blue" ? 3 : 1
  
  environment = {
    APP_COLOR = "blue"
  }
}

# Green environment
resource "sevalla_application" "app_green" {
  name = "myapp-green"
  
  repository {
    url  = "https://github.com/company/app"
    type = "github"
  }
  
  branch    = "release/green"
  instances = local.active_color == "green" ? 3 : 1
  
  environment = {
    APP_COLOR = "green"
  }
}

# DNS pointing to active environment
output "active_app_url" {
  value = local.active_color == "blue" ? 
    sevalla_application.app_blue.domain : 
    sevalla_application.app_green.domain
}
```

### Auto-scaling Configuration

```hcl
variable "traffic_level" {
  description = "Expected traffic level"
  type        = string
  default     = "normal"
  
  validation {
    condition     = contains(["low", "normal", "high", "peak"], var.traffic_level)
    error_message = "Traffic level must be low, normal, high, or peak."
  }
}

locals {
  scaling_config = {
    low = {
      instances = 1
      memory    = 512
      cpu       = 250
    }
    normal = {
      instances = 2
      memory    = 1024
      cpu       = 500
    }
    high = {
      instances = 4
      memory    = 2048
      cpu       = 1000
    }
    peak = {
      instances = 8
      memory    = 4096
      cpu       = 2000
    }
  }
}

resource "sevalla_application" "scalable_app" {
  name = "auto-scaling-app"
  
  instances = local.scaling_config[var.traffic_level].instances
  memory    = local.scaling_config[var.traffic_level].memory
  cpu       = local.scaling_config[var.traffic_level].cpu
  
  # ... other configuration
}
```

### Multi-Region Deployment

```hcl
variable "regions" {
  description = "Regions to deploy to"
  type        = list(string)
  default     = ["us-east-1", "eu-west-1", "ap-southeast-1"]
}

# Deploy application to each region
resource "sevalla_application" "regional_app" {
  for_each = toset(var.regions)
  
  name = "myapp-${each.key}"
  
  repository {
    url  = "https://github.com/company/app"
    type = "github"
  }
  
  environment = {
    REGION = each.key
    # Region-specific configuration
  }
  
  # ... other configuration
}

# Object storage in each region
resource "sevalla_object_storage" "regional_storage" {
  for_each = toset(var.regions)
  
  name   = "storage-${each.key}"
  region = each.key
}
```

## Testing Your Configuration

### Validate Configuration

```bash
# Check syntax and configuration
terraform validate

# Format your configuration files
terraform fmt -recursive

# Show the execution plan
terraform plan
```

### Test in Staging First

```hcl
# Use workspaces for different environments
terraform workspace new staging
terraform workspace select staging
terraform apply

# Test your staging deployment
# ...

# Switch to production
terraform workspace select production
terraform apply
```

### Running Tests

If you're contributing to the provider development:

```bash
# Run unit tests
make test

# Run acceptance tests (requires SEVALLA_TOKEN)
make testacc
```

## Security Best Practices

### 1. Never Commit Secrets

```hcl
# ❌ BAD - Don't do this
provider "sevalla" {
  token = "sk_live_abcd1234"  # Never commit tokens!
}

# ✅ GOOD - Use environment variables
provider "sevalla" {
  # Token read from SEVALLA_TOKEN env var
}
```

### 2. Use Terraform Cloud for Secrets

```hcl
terraform {
  cloud {
    organization = "my-company"
    
    workspaces {
      name = "sevalla-prod"
    }
  }
}
```

### 3. Least Privilege Access

Create separate API tokens with minimal required permissions:
- Read-only tokens for `terraform plan`
- Write tokens only for CI/CD systems

### 4. Encrypt State Files

```hcl
terraform {
  backend "s3" {
    bucket     = "my-terraform-state"
    key        = "sevalla/terraform.tfstate"
    region     = "us-east-1"
    encrypt    = true  # Enable encryption
    kms_key_id = "arn:aws:kms:us-east-1:123456789012:key/12345678"
  }
}
```

## Performance Optimization

### 1. Use Data Sources to Reduce API Calls

```hcl
# Instead of creating multiple similar resources
data "sevalla_application" "existing" {
  id = "app-12345"
}

# Reference the data source
resource "sevalla_pipeline" "pipeline" {
  app_id = data.sevalla_application.existing.id
  # ...
}
```

### 2. Parallelize Resource Creation

```bash
# Default parallelism is 10
terraform apply

# Increase for faster execution (be mindful of rate limits)
terraform apply -parallelism=20

# Decrease if hitting rate limits
terraform apply -parallelism=1
```

### 3. Use Targeted Applies for Large Configurations

```bash
# Apply only specific resources
terraform apply -target=sevalla_application.web_app
terraform apply -target=module.database
```

## Monitoring and Observability

### Export Metrics

```hcl
# Output important metrics
output "infrastructure_summary" {
  value = {
    applications = {
      count = length(sevalla_application.apps)
      total_memory = sum([for app in sevalla_application.apps : app.memory * app.instances])
      total_instances = sum([for app in sevalla_application.apps : app.instances])
    }
    databases = {
      count = length(sevalla_database.dbs)
      types = distinct([for db in sevalla_database.dbs : db.type])
    }
  }
}
```

### Integration with Monitoring Tools

```hcl
# Send deployment notifications
resource "null_resource" "notify_deployment" {
  triggers = {
    app_version = sevalla_application.app.id
  }
  
  provisioner "local-exec" {
    command = <<-EOT
      curl -X POST https://api.monitoring-tool.com/deployments \
        -H "Authorization: Bearer $MONITORING_TOKEN" \
        -d '{
          "service": "${sevalla_application.app.name}",
          "version": "${sevalla_application.app.id}",
          "environment": "production"
        }'
    EOT
  }
}
```

## FAQ

**Q: How do I handle provider upgrades?**
A: Test in a non-production environment first:
```bash
terraform init -upgrade
terraform plan
```

**Q: Can I use this provider with Terraform Cloud?**
A: Yes, upload the provider binary to a private registry or use the public registry once available.

**Q: How do I manage multiple Sevalla accounts?**
A: Use provider aliases:
```hcl
provider "sevalla" {
  alias = "staging"
  token = var.staging_token
}

provider "sevalla" {
  alias = "production"
  token = var.production_token
}

resource "sevalla_application" "staging_app" {
  provider = sevalla.staging
  # ...
}
```

**Q: What happens if a resource is modified outside Terraform?**
A: Run `terraform refresh` to update the state, or `terraform plan` to see the differences.

## Support

- **Documentation**: This README and the `/examples` directory
- **Issues**: [GitHub Issues](https://github.com/sriniously/terraform-provider-sevalla/issues)
- **Community**: [Sevalla Community Forum](https://community.sevalla.com)
- **API Reference**: [api-docs.sevalla.com](https://api-docs.sevalla.com)

## Development Status

This provider is actively maintained and supports all core Sevalla resources. Current implementation status:

### ✅ Fully Implemented
- **Resources**: All 5 resources (application, database, static_site, object_storage, pipeline)
- **Data Sources**: 4 data sources (application, database, static_site, object_storage)
- **API Client**: Complete REST API client with error handling
- **Documentation**: Comprehensive documentation for all resources

### 🧪 Testing Status
- **Application Resource**: ✅ Full test coverage
- **Database Resource**: ✅ Full test coverage  
- **Static Site Resource**: ⚠️ Implementation complete, tests pending
- **Object Storage Resource**: ⚠️ Implementation complete, tests pending
- **Pipeline Resource**: ⚠️ Implementation complete, tests pending

### 📋 TODO List
- [ ] Add tests for remaining resources
- [ ] Add pipeline data source
- [ ] Add more comprehensive integration tests
- [ ] Add performance optimizations

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Run tests (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.