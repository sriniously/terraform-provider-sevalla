# Full-Stack Application Example

This example demonstrates how to deploy a complete full-stack application using the Sevalla Terraform provider. It includes:

- **Frontend**: React/Vue static site with custom domain
- **Backend API**: Node.js/Python application with load balancing
- **Database**: PostgreSQL for persistent data
- **Cache**: Redis for session storage and caching
- **Storage**: Object storage for user uploads and static assets
- **CI/CD**: Automated deployment pipelines

## Architecture

```
Internet
    ↓
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │
│   (Static Site) │────│   (Application) │
│                 │    │                 │
└─────────────────┘    └─────────────────┘
                              ↓
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │   (Database)    │
                    └─────────────────┘
                              ↓
                    ┌─────────────────┐
                    │      Redis      │
                    │     (Cache)     │
                    └─────────────────┘
                              ↓
                    ┌─────────────────┐
                    │ Object Storage  │
                    │ (Uploads/Assets)│
                    └─────────────────┘
```

## Prerequisites

1. **Sevalla Account**: Sign up at [sevalla.com](https://sevalla.com)
2. **API Token**: Generate an API token from your account settings
3. **GitHub Repositories**: Have your frontend and backend code in GitHub repositories
4. **Domain**: A domain name configured to point to Sevalla (optional)

## Quick Start

1. **Clone and Setup**:
   ```bash
   git clone https://github.com/sriniously/terraform-provider-sevalla.git
   cd terraform-provider-sevalla/examples/full-stack-app
   ```

2. **Configure Variables**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your values
   ```

3. **Deploy**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Configuration

### Required Variables

- `sevalla_token`: Your Sevalla API token
- `domain`: Your primary domain (e.g., "mycompany.com")
- `api_repo_url`: GitHub URL for your backend API
- `frontend_repo_url`: GitHub URL for your frontend
- `db_password`: Strong password for PostgreSQL
- `jwt_secret`: Secret for JWT authentication (32+ characters)

### Optional Variables

- `app_name`: Application name prefix (default: "myapp")
- `region`: Deployment region (default: "us-east-1")
- `api_branch`: Git branch for API (default: "main")
- `frontend_branch`: Git branch for frontend (default: "main")
- `api_instances`: Number of API instances (default: 2)

## Application Requirements

### Backend API Requirements

Your backend application should:

1. **Database Connection**: Use the `DATABASE_URL` environment variable
2. **Redis Connection**: Use the `REDIS_URL` environment variable
3. **Object Storage**: Use the object storage environment variables for file uploads
4. **Health Check**: Respond to health checks on the configured port
5. **Build Process**: Include `npm ci && npm run build` or equivalent

Example environment variables your API will receive:
```bash
DATABASE_URL=postgresql://username:password@host:port/database
REDIS_URL=redis://host:port
UPLOADS_BUCKET_ENDPOINT=https://uploads.s3.region.amazonaws.com
UPLOADS_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
UPLOADS_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Frontend Requirements

Your frontend application should:

1. **API Configuration**: Use build-time environment variables to configure API endpoints
2. **Build Process**: Include `npm ci && npm run build` or equivalent
3. **Build Output**: Generate static files in a `dist` or `build` directory
4. **Environment Variables**: Use `REACT_APP_API_URL` or equivalent for API endpoint

Example build command:
```bash
REACT_APP_API_URL=https://api.mycompany.com npm run build
```

## Scaling and Performance

### Horizontal Scaling
```hcl
# Scale API instances based on load
api_instances = 5
```

### Vertical Scaling
```hcl
# Increase memory and CPU for API
memory = 2048  # 2GB
cpu    = 1000  # 1 CPU core
```

### Database Scaling
```hcl
# Upgrade database size
size = "large"  # or "xlarge"
```

## Security Best Practices

1. **Environment Variables**: Store sensitive data as environment variables
2. **Secrets Management**: Use strong, unique passwords and secrets
3. **CORS Configuration**: Configure CORS_ORIGIN to match your domain
4. **HTTPS**: All resources use HTTPS by default
5. **Database Security**: Database passwords are marked as sensitive

## Monitoring and Debugging

### Check Application Status
```bash
# View current infrastructure
terraform show

# Check specific resource
terraform show 'sevalla_application.api'
```

### Access Logs
Check your Sevalla dashboard for:
- Application logs
- Database performance metrics
- Storage usage statistics
- Pipeline deployment history

## Customization Examples

### Add a Worker Service
```hcl
resource "sevalla_application" "worker" {
  name = "${var.app_name}-worker"
  
  repository {
    url    = var.worker_repo_url
    type   = "github"
    branch = "main"
  }
  
  environment = {
    DATABASE_URL = "postgresql://${sevalla_database.main_db.username}:${var.db_password}@${sevalla_database.main_db.host}:${sevalla_database.main_db.port}/${sevalla_database.main_db.name}"
    REDIS_URL    = "redis://${sevalla_database.cache.host}:${sevalla_database.cache.port}"
    QUEUE_NAME   = "background-jobs"
  }
  
  instances = 1
  memory    = 512
  cpu       = 250
}
```

### Add a Staging Environment
```hcl
# Use Terraform workspaces
terraform workspace new staging
terraform workspace select staging

# Or use different variable values
app_name = "myapp-staging"
domain   = "staging.mycompany.com"
api_instances = 1
```

### Add Multiple Regions
```hcl
# Deploy to multiple regions
locals {
  regions = ["us-east-1", "eu-west-1"]
}

resource "sevalla_application" "api" {
  for_each = toset(local.regions)
  
  name = "${var.app_name}-api-${each.key}"
  # ... rest of configuration
}
```

## Troubleshooting

### Common Issues

1. **Build Failures**: Check repository access and build commands
2. **Database Connection**: Verify password and connection string format
3. **Domain Configuration**: Ensure DNS is properly configured
4. **Resource Limits**: Check if you're hitting account limits

### Debug Commands
```bash
# Enable debug logging
export TF_LOG=DEBUG

# Validate configuration
terraform validate

# Check planned changes
terraform plan

# Import existing resource
terraform import sevalla_application.api app-12345
```

## Cost Optimization

- **Right-size Resources**: Start with smaller instances and scale up
- **Use Staging**: Deploy to staging first with minimal resources
- **Monitor Usage**: Check storage and database usage regularly
- **Auto-deployment**: Use pipelines to reduce manual deployment overhead

## Migration from Existing Infrastructure

If you have existing resources in Sevalla:

1. **List Resources**: Note IDs from Sevalla dashboard
2. **Import Resources**: Use `terraform import` for each resource
3. **Verify State**: Run `terraform plan` to ensure no changes
4. **Gradually Migrate**: Move one resource at a time

## Support

- **Provider Issues**: [GitHub Issues](https://github.com/sriniously/terraform-provider-sevalla/issues)
- **Sevalla Platform**: [support@sevalla.com](mailto:support@sevalla.com)
- **Documentation**: [docs.sevalla.com](https://docs.sevalla.com)