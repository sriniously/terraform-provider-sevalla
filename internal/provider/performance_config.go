package provider

import (
	"os"
	"strconv"
	"time"
)

// PerformanceConfig holds configuration for performance optimizations.
type PerformanceConfig struct {
	// Caching configuration
	CacheEnabled bool
	CacheTTL     time.Duration

	// Rate limiting configuration
	RateLimitEnabled   bool
	RateLimitPerSecond int
	RateLimitBurst     int

	// Batch processing configuration
	BatchEnabled bool
	BatchSize    int
	BatchTimeout time.Duration

	// Connection pooling configuration
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Request timeout configuration
	RequestTimeout time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
}

// DefaultPerformanceConfig returns default performance configuration.
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		// Caching defaults
		CacheEnabled: true,
		CacheTTL:     5 * time.Minute,

		// Rate limiting defaults
		RateLimitEnabled:   true,
		RateLimitPerSecond: 10,
		RateLimitBurst:     20,

		// Batch processing defaults
		BatchEnabled: true,
		BatchSize:    10,
		BatchTimeout: 100 * time.Millisecond,

		// Connection pooling defaults
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,

		// Request timeout defaults
		RequestTimeout: 30 * time.Second,
		RetryAttempts:  3,
		RetryDelay:     1 * time.Second,
	}
}

// LoadPerformanceConfigFromEnv loads performance configuration from environment variables.
func LoadPerformanceConfigFromEnv() *PerformanceConfig {
	config := DefaultPerformanceConfig()

	loadCacheConfig(config)
	loadRateLimitConfig(config)
	loadBatchConfig(config)
	loadConnectionConfig(config)
	loadRequestConfig(config)

	return config
}

// loadCacheConfig loads cache configuration from environment variables.
func loadCacheConfig(config *PerformanceConfig) {
	if val := os.Getenv("SEVALLA_CACHE_ENABLED"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.CacheEnabled = enabled
		}
	}

	if val := os.Getenv("SEVALLA_CACHE_TTL"); val != "" {
		if ttl, err := time.ParseDuration(val); err == nil {
			config.CacheTTL = ttl
		}
	}
}

// loadRateLimitConfig loads rate limiting configuration from environment variables.
func loadRateLimitConfig(config *PerformanceConfig) {
	if val := os.Getenv("SEVALLA_RATE_LIMIT_ENABLED"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.RateLimitEnabled = enabled
		}
	}

	if val := os.Getenv("SEVALLA_RATE_LIMIT_PER_SECOND"); val != "" {
		if limit, err := strconv.Atoi(val); err == nil {
			config.RateLimitPerSecond = limit
		}
	}

	if val := os.Getenv("SEVALLA_RATE_LIMIT_BURST"); val != "" {
		if burst, err := strconv.Atoi(val); err == nil {
			config.RateLimitBurst = burst
		}
	}
}

// loadBatchConfig loads batch processing configuration from environment variables.
func loadBatchConfig(config *PerformanceConfig) {
	if val := os.Getenv("SEVALLA_BATCH_ENABLED"); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			config.BatchEnabled = enabled
		}
	}

	if val := os.Getenv("SEVALLA_BATCH_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.BatchSize = size
		}
	}

	if val := os.Getenv("SEVALLA_BATCH_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.BatchTimeout = timeout
		}
	}
}

// loadConnectionConfig loads connection pooling configuration from environment variables.
func loadConnectionConfig(config *PerformanceConfig) {
	if val := os.Getenv("SEVALLA_MAX_IDLE_CONNS"); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			config.MaxIdleConns = conns
		}
	}

	if val := os.Getenv("SEVALLA_MAX_OPEN_CONNS"); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			config.MaxOpenConns = conns
		}
	}

	if val := os.Getenv("SEVALLA_CONN_MAX_LIFETIME"); val != "" {
		if lifetime, err := time.ParseDuration(val); err == nil {
			config.ConnMaxLifetime = lifetime
		}
	}

	if val := os.Getenv("SEVALLA_CONN_MAX_IDLE_TIME"); val != "" {
		if idleTime, err := time.ParseDuration(val); err == nil {
			config.ConnMaxIdleTime = idleTime
		}
	}
}

// loadRequestConfig loads request timeout configuration from environment variables.
func loadRequestConfig(config *PerformanceConfig) {
	if val := os.Getenv("SEVALLA_REQUEST_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.RequestTimeout = timeout
		}
	}

	if val := os.Getenv("SEVALLA_RETRY_ATTEMPTS"); val != "" {
		if attempts, err := strconv.Atoi(val); err == nil {
			config.RetryAttempts = attempts
		}
	}

	if val := os.Getenv("SEVALLA_RETRY_DELAY"); val != "" {
		if delay, err := time.ParseDuration(val); err == nil {
			config.RetryDelay = delay
		}
	}
}

// Validate validates the performance configuration.
func (pc *PerformanceConfig) Validate() error {
	if pc.RateLimitPerSecond <= 0 {
		pc.RateLimitPerSecond = 10
	}

	if pc.RateLimitBurst <= 0 {
		pc.RateLimitBurst = pc.RateLimitPerSecond * 2
	}

	if pc.BatchSize <= 0 {
		pc.BatchSize = 10
	}

	if pc.BatchTimeout <= 0 {
		pc.BatchTimeout = 100 * time.Millisecond
	}

	if pc.MaxIdleConns <= 0 {
		pc.MaxIdleConns = 10
	}

	if pc.MaxOpenConns <= 0 {
		pc.MaxOpenConns = 20
	}

	if pc.ConnMaxLifetime <= 0 {
		pc.ConnMaxLifetime = 30 * time.Minute
	}

	if pc.ConnMaxIdleTime <= 0 {
		pc.ConnMaxIdleTime = 10 * time.Minute
	}

	if pc.RequestTimeout <= 0 {
		pc.RequestTimeout = 30 * time.Second
	}

	if pc.RetryAttempts < 0 {
		pc.RetryAttempts = 3
	}

	if pc.RetryDelay <= 0 {
		pc.RetryDelay = 1 * time.Second
	}

	if pc.CacheTTL <= 0 {
		pc.CacheTTL = 5 * time.Minute
	}

	return nil
}
