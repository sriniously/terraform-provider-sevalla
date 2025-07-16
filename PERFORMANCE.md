# Performance Optimization Guide

This document describes the performance optimizations implemented in the Sevalla Terraform provider and how to configure them for optimal performance.

## Overview

The Sevalla Terraform provider includes several performance optimizations to reduce API calls, improve response times, and handle rate limiting effectively:

- **Caching**: Reduces redundant API calls by caching responses
- **Rate Limiting**: Prevents API rate limit violations
- **Batch Processing**: Groups related operations for efficiency
- **Connection Pooling**: Manages HTTP connections efficiently
- **Request Retry**: Handles temporary failures gracefully

## Configuration

### Environment Variables

You can configure performance settings using environment variables:

#### Caching Configuration
```bash
# Enable/disable caching (default: true)
export SEVALLA_CACHE_ENABLED=true

# Cache TTL (default: 5m)
export SEVALLA_CACHE_TTL=5m
```

#### Rate Limiting Configuration
```bash
# Enable/disable rate limiting (default: true)
export SEVALLA_RATE_LIMIT_ENABLED=true

# Requests per second (default: 10)
export SEVALLA_RATE_LIMIT_PER_SECOND=10

# Burst capacity (default: 20)
export SEVALLA_RATE_LIMIT_BURST=20
```

#### Batch Processing Configuration
```bash
# Enable/disable batch processing (default: true)
export SEVALLA_BATCH_ENABLED=true

# Batch size (default: 10)
export SEVALLA_BATCH_SIZE=10

# Batch timeout (default: 100ms)
export SEVALLA_BATCH_TIMEOUT=100ms
```

#### Connection Pooling Configuration
```bash
# Maximum idle connections (default: 10)
export SEVALLA_MAX_IDLE_CONNS=10

# Maximum open connections (default: 20)
export SEVALLA_MAX_OPEN_CONNS=20

# Connection max lifetime (default: 30m)
export SEVALLA_CONN_MAX_LIFETIME=30m

# Connection max idle time (default: 10m)
export SEVALLA_CONN_MAX_IDLE_TIME=10m
```

#### Request Timeout Configuration
```bash
# Request timeout (default: 30s)
export SEVALLA_REQUEST_TIMEOUT=30s

# Retry attempts (default: 3)
export SEVALLA_RETRY_ATTEMPTS=3

# Retry delay (default: 1s)
export SEVALLA_RETRY_DELAY=1s
```

## Performance Features

### 1. Caching

The provider implements intelligent caching for read operations:

- **Resource Data**: Application, database, static site, object storage, and pipeline data
- **Cache Keys**: Based on resource type and ID
- **TTL**: Configurable time-to-live (default: 5 minutes)
- **Cache Invalidation**: Automatic invalidation on resource updates

#### Benefits:
- Reduces API calls for repeated reads
- Improves `terraform plan` performance
- Reduces latency for data sources

#### When to Disable:
- When working with rapidly changing resources
- During development/debugging
- When strict consistency is required

### 2. Rate Limiting

Built-in rate limiting prevents API rate limit violations:

- **Token Bucket Algorithm**: Smooth rate limiting
- **Configurable Limits**: Requests per second and burst capacity
- **Backoff Strategy**: Automatic retry with exponential backoff
- **Context Awareness**: Respects Terraform's context cancellation

#### Benefits:
- Prevents 429 (Too Many Requests) errors
- Ensures stable API performance
- Reduces need for manual delays

#### Configuration Tips:
- Set limits based on your Sevalla API plan
- Increase burst capacity for large deployments
- Monitor API usage patterns

### 3. Batch Processing

Groups related operations for efficiency:

- **Operation Grouping**: Similar operations are batched together
- **Configurable Batch Size**: Control how many operations are batched
- **Timeout-based Flushing**: Ensures operations don't wait indefinitely
- **Type-specific Batching**: Different resource types can be batched differently

#### Benefits:
- Reduces total API calls
- Improves performance for large configurations
- Better resource utilization

#### Best Practices:
- Keep batch sizes reasonable (5-20 operations)
- Use shorter timeouts for interactive operations
- Monitor batch effectiveness in logs

### 4. Connection Pooling

Efficient HTTP connection management:

- **Connection Reuse**: Reuses existing connections
- **Configurable Limits**: Control connection pool size
- **Idle Connection Management**: Automatic cleanup of unused connections
- **Connection Lifecycle**: Proper connection lifecycle management

#### Benefits:
- Reduces connection overhead
- Improves overall throughput
- Better resource utilization

#### Configuration Guidelines:
- Set pool size based on concurrency needs
- Monitor connection usage patterns
- Adjust timeouts based on network conditions

### 5. Request Retry

Handles temporary failures gracefully:

- **Exponential Backoff**: Increasing delays between retries
- **Configurable Attempts**: Control retry behavior
- **Error Classification**: Only retries appropriate errors
- **Circuit Breaker**: Prevents cascading failures

#### Benefits:
- Improves reliability
- Handles transient network issues
- Reduces manual intervention

#### Configuration:
- Set retry attempts based on reliability needs
- Adjust delay based on expected recovery time
- Monitor retry patterns

## Performance Tuning

### For Small Deployments (< 10 resources)
```bash
export SEVALLA_CACHE_TTL=2m
export SEVALLA_RATE_LIMIT_PER_SECOND=5
export SEVALLA_BATCH_SIZE=5
export SEVALLA_MAX_OPEN_CONNS=10
```

### For Medium Deployments (10-50 resources)
```bash
export SEVALLA_CACHE_TTL=5m
export SEVALLA_RATE_LIMIT_PER_SECOND=10
export SEVALLA_BATCH_SIZE=10
export SEVALLA_MAX_OPEN_CONNS=20
```

### For Large Deployments (> 50 resources)
```bash
export SEVALLA_CACHE_TTL=10m
export SEVALLA_RATE_LIMIT_PER_SECOND=20
export SEVALLA_BATCH_SIZE=20
export SEVALLA_MAX_OPEN_CONNS=50
```

### For CI/CD Environments
```bash
export SEVALLA_CACHE_ENABLED=false
export SEVALLA_RATE_LIMIT_PER_SECOND=15
export SEVALLA_BATCH_SIZE=15
export SEVALLA_RETRY_ATTEMPTS=5
export SEVALLA_RETRY_DELAY=2s
```

## Monitoring and Debugging

### Enable Debug Logging
```bash
export TF_LOG=DEBUG
terraform apply
```

### Performance Metrics

The provider logs performance metrics including:
- Cache hit/miss ratios
- API call counts
- Batch processing efficiency
- Rate limiting events
- Connection pool usage

### Common Performance Issues

1. **High Cache Miss Rate**
   - Check cache TTL settings
   - Verify resource access patterns
   - Consider increasing cache TTL

2. **Rate Limit Violations**
   - Reduce rate limit settings
   - Increase retry delays
   - Use batch processing

3. **Slow Response Times**
   - Check connection pool settings
   - Verify network connectivity
   - Monitor API endpoint health

4. **Memory Usage**
   - Monitor cache size
   - Adjust cache TTL
   - Clear cache periodically

## Best Practices

### 1. Resource Organization
- Group related resources in modules
- Use data sources for cross-references
- Minimize resource dependencies

### 2. State Management
- Use remote state for team collaboration
- Enable state locking
- Regular state cleanup

### 3. Plan Optimization
- Use `terraform plan -target` for focused changes
- Leverage `terraform refresh` for state updates
- Cache plan files for repeated operations

### 4. Testing
- Use separate environments for testing
- Mock API calls in unit tests
- Load test with realistic data

### 5. Monitoring
- Track API usage patterns
- Monitor provider performance
- Set up alerts for failures

## Troubleshooting

### Performance Issues

1. **Slow terraform plan**
   - Enable caching
   - Reduce batch timeout
   - Check network connectivity

2. **API rate limit errors**
   - Reduce rate limit settings
   - Increase retry delays
   - Use parallelism controls

3. **Memory usage**
   - Clear cache periodically
   - Reduce cache TTL
   - Monitor resource usage

### Debug Commands

```bash
# Check current performance settings
terraform show -json | jq '.configuration.provider_config.sevalla'

# Monitor API calls
TF_LOG=DEBUG terraform plan 2>&1 | grep -i "api"

# Check cache statistics
TF_LOG=DEBUG terraform plan 2>&1 | grep -i "cache"

# Monitor rate limiting
TF_LOG=DEBUG terraform plan 2>&1 | grep -i "rate"
```

## Contributing

When contributing performance improvements:

1. Measure performance impact
2. Add appropriate tests
3. Update documentation
4. Consider backward compatibility
5. Monitor resource usage

## Support

For performance-related issues:
- Check the troubleshooting guide
- Enable debug logging
- Report issues with performance metrics
- Include configuration details