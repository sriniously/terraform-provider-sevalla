package provider

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// CacheEntry represents a cached API response.
type CacheEntry struct {
	Data      interface{}
	Timestamp time.Time
	TTL       time.Duration
}

// IsExpired checks if the cache entry is expired.
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.Timestamp) > c.TTL
}

// ProviderCache provides caching for API responses to reduce API calls.
type ProviderCache struct {
	cache map[string]*CacheEntry
	mutex sync.RWMutex
}

// NewProviderCache creates a new provider cache.
func NewProviderCache() *ProviderCache {
	return &ProviderCache{
		cache: make(map[string]*CacheEntry),
	}
}

// Get retrieves an item from the cache.
func (pc *ProviderCache) Get(key string) (interface{}, bool) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()
	
	entry, exists := pc.cache[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}
	
	return entry.Data, true
}

// Set stores an item in the cache.
func (pc *ProviderCache) Set(key string, data interface{}, ttl time.Duration) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	
	pc.cache[key] = &CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

// Clear removes all items from the cache.
func (pc *ProviderCache) Clear() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	
	pc.cache = make(map[string]*CacheEntry)
}

// ClearExpired removes all expired entries from the cache.
func (pc *ProviderCache) ClearExpired() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	
	for key, entry := range pc.cache {
		if entry.IsExpired() {
			delete(pc.cache, key)
		}
	}
}

// BatchOperation represents a batch operation for API calls.
type BatchOperation struct {
	ID         string
	Operation  string
	Parameters interface{}
	Result     interface{}
	Error      error
	Done       chan bool
}

// BatchProcessor handles batch operations to reduce API calls.
type BatchProcessor struct {
	operations chan *BatchOperation
	results    map[string]*BatchOperation
	mutex      sync.RWMutex
	batchSize  int
	batchTime  time.Duration
}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor(batchSize int, batchTime time.Duration) *BatchProcessor {
	bp := &BatchProcessor{
		operations: make(chan *BatchOperation, batchSize*2),
		results:    make(map[string]*BatchOperation),
		batchSize:  batchSize,
		batchTime:  batchTime,
	}
	
	// Start the batch processor
	go bp.processBatches()
	
	return bp
}

// Submit submits an operation to the batch processor.
func (bp *BatchProcessor) Submit(op *BatchOperation) {
	bp.mutex.Lock()
	bp.results[op.ID] = op
	bp.mutex.Unlock()
	
	bp.operations <- op
}

// Wait waits for an operation to complete.
func (bp *BatchProcessor) Wait(id string) (*BatchOperation, error) {
	bp.mutex.RLock()
	op, exists := bp.results[id]
	bp.mutex.RUnlock()
	
	if !exists {
		return nil, nil
	}
	
	<-op.Done
	return op, op.Error
}

// processBatches processes operations in batches.
func (bp *BatchProcessor) processBatches() {
	ticker := time.NewTicker(bp.batchTime)
	defer ticker.Stop()
	
	batch := make([]*BatchOperation, 0, bp.batchSize)
	
	for {
		select {
		case op := <-bp.operations:
			batch = append(batch, op)
			
			if len(batch) >= bp.batchSize {
				bp.executeBatch(batch)
				batch = make([]*BatchOperation, 0, bp.batchSize)
			}
			
		case <-ticker.C:
			if len(batch) > 0 {
				bp.executeBatch(batch)
				batch = make([]*BatchOperation, 0, bp.batchSize)
			}
		}
	}
}

// executeBatch executes a batch of operations.
func (bp *BatchProcessor) executeBatch(batch []*BatchOperation) {
	// Group operations by type for more efficient processing
	operationGroups := make(map[string][]*BatchOperation)
	
	for _, op := range batch {
		operationGroups[op.Operation] = append(operationGroups[op.Operation], op)
	}
	
	// Execute each group
	for operationType, ops := range operationGroups {
		switch operationType {
		case "get_application":
			bp.executeGetApplicationBatch(ops)
		case "get_database":
			bp.executeGetDatabaseBatch(ops)
		case "get_static_site":
			bp.executeGetStaticSiteBatch(ops)
		case "get_object_storage":
			bp.executeGetObjectStorageBatch(ops)
		case "get_pipeline":
			bp.executeGetPipelineBatch(ops)
		default:
			// Execute individually if no batch support
			for _, op := range ops {
				bp.executeIndividualOperation(op)
			}
		}
	}
}

// executeGetApplicationBatch executes a batch of get application operations.
func (bp *BatchProcessor) executeGetApplicationBatch(ops []*BatchOperation) {
	// In a real implementation, this would make a batch API call
	// For now, we'll execute individually but could be optimized
	for _, op := range ops {
		bp.executeIndividualOperation(op)
	}
}

// executeGetDatabaseBatch executes a batch of get database operations.
func (bp *BatchProcessor) executeGetDatabaseBatch(ops []*BatchOperation) {
	// In a real implementation, this would make a batch API call
	for _, op := range ops {
		bp.executeIndividualOperation(op)
	}
}

// executeGetStaticSiteBatch executes a batch of get static site operations.
func (bp *BatchProcessor) executeGetStaticSiteBatch(ops []*BatchOperation) {
	// In a real implementation, this would make a batch API call
	for _, op := range ops {
		bp.executeIndividualOperation(op)
	}
}

// executeGetObjectStorageBatch executes a batch of get object storage operations.
func (bp *BatchProcessor) executeGetObjectStorageBatch(ops []*BatchOperation) {
	// In a real implementation, this would make a batch API call
	for _, op := range ops {
		bp.executeIndividualOperation(op)
	}
}

// executeGetPipelineBatch executes a batch of get pipeline operations.
func (bp *BatchProcessor) executeGetPipelineBatch(ops []*BatchOperation) {
	// In a real implementation, this would make a batch API call
	for _, op := range ops {
		bp.executeIndividualOperation(op)
	}
}

// executeIndividualOperation executes a single operation.
func (bp *BatchProcessor) executeIndividualOperation(op *BatchOperation) {
	// This would contain the actual API call logic
	// For now, we'll just mark it as done
	close(op.Done)
}

// RateLimiter implements rate limiting for API calls.
type RateLimiter struct {
	tokens    chan struct{}
	ticker    *time.Ticker
	rateLimit int
	interval  time.Duration
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rateLimit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:    make(chan struct{}, rateLimit),
		ticker:    time.NewTicker(interval),
		rateLimit: rateLimit,
		interval:  interval,
	}
	
	// Fill the token bucket initially
	for i := 0; i < rateLimit; i++ {
		rl.tokens <- struct{}{}
	}
	
	// Start the token refill process
	go rl.refillTokens()
	
	return rl
}

// Wait waits for a token to be available.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// refillTokens refills the token bucket at the specified interval.
func (rl *RateLimiter) refillTokens() {
	for range rl.ticker.C {
		select {
		case rl.tokens <- struct{}{}:
		default:
			// Token bucket is full, skip
		}
	}
}

// Stop stops the rate limiter.
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
}

// PerformanceOptimizedClient wraps the Sevalla API client with performance optimizations.
type PerformanceOptimizedClient struct {
	client        *sevallaapi.Client
	cache         *ProviderCache
	batchProcessor *BatchProcessor
	rateLimiter   *RateLimiter
}

// NewPerformanceOptimizedClient creates a new performance optimized client.
func NewPerformanceOptimizedClient(client *sevallaapi.Client) *PerformanceOptimizedClient {
	return &PerformanceOptimizedClient{
		client:        client,
		cache:         NewProviderCache(),
		batchProcessor: NewBatchProcessor(10, 100*time.Millisecond),
		rateLimiter:   NewRateLimiter(10, 1*time.Second),
	}
}

// GetApplicationCached gets an application with caching.
func (poc *PerformanceOptimizedClient) GetApplicationCached(ctx context.Context, id string) (*sevallaapi.Application, error) {
	cacheKey := "application:" + id
	
	// Check cache first
	if cached, found := poc.cache.Get(cacheKey); found {
		tflog.Debug(ctx, "Application retrieved from cache", map[string]interface{}{"id": id})
		if app, ok := cached.(*sevallaapi.Application); ok {
			return app, nil
		}
	}
	
	// Wait for rate limiter
	if err := poc.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	
	// Make API call
	tflog.Debug(ctx, "Making API call for application", map[string]interface{}{"id": id})
	app, err := sevallaapi.NewApplicationService(poc.client).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	poc.cache.Set(cacheKey, app, 5*time.Minute)
	
	return app, nil
}

// GetDatabaseCached gets a database with caching.
func (poc *PerformanceOptimizedClient) GetDatabaseCached(ctx context.Context, id string) (*sevallaapi.Database, error) {
	cacheKey := "database:" + id
	
	// Check cache first
	if cached, found := poc.cache.Get(cacheKey); found {
		tflog.Debug(ctx, "Database retrieved from cache", map[string]interface{}{"id": id})
		if db, ok := cached.(*sevallaapi.Database); ok {
			return db, nil
		}
	}
	
	// Wait for rate limiter
	if err := poc.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	
	// Make API call
	tflog.Debug(ctx, "Making API call for database", map[string]interface{}{"id": id})
	db, err := sevallaapi.NewDatabaseService(poc.client).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	poc.cache.Set(cacheKey, db, 5*time.Minute)
	
	return db, nil
}

// GetStaticSiteCached gets a static site with caching.
func (poc *PerformanceOptimizedClient) GetStaticSiteCached(ctx context.Context, id string) (*sevallaapi.StaticSite, error) {
	cacheKey := "static_site:" + id
	
	// Check cache first
	if cached, found := poc.cache.Get(cacheKey); found {
		tflog.Debug(ctx, "Static site retrieved from cache", map[string]interface{}{"id": id})
		if site, ok := cached.(*sevallaapi.StaticSite); ok {
			return site, nil
		}
	}
	
	// Wait for rate limiter
	if err := poc.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	
	// Make API call
	tflog.Debug(ctx, "Making API call for static site", map[string]interface{}{"id": id})
	site, err := sevallaapi.NewStaticSiteService(poc.client).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	poc.cache.Set(cacheKey, site, 5*time.Minute)
	
	return site, nil
}

// GetObjectStorageCached gets object storage with caching.
func (poc *PerformanceOptimizedClient) GetObjectStorageCached(ctx context.Context, id string) (*sevallaapi.ObjectStorage, error) {
	cacheKey := "object_storage:" + id
	
	// Check cache first
	if cached, found := poc.cache.Get(cacheKey); found {
		tflog.Debug(ctx, "Object storage retrieved from cache", map[string]interface{}{"id": id})
		if storage, ok := cached.(*sevallaapi.ObjectStorage); ok {
			return storage, nil
		}
	}
	
	// Wait for rate limiter
	if err := poc.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	
	// Make API call
	tflog.Debug(ctx, "Making API call for object storage", map[string]interface{}{"id": id})
	storage, err := sevallaapi.NewObjectStorageService(poc.client).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	poc.cache.Set(cacheKey, storage, 5*time.Minute)
	
	return storage, nil
}

// GetPipelineCached gets a pipeline with caching.
func (poc *PerformanceOptimizedClient) GetPipelineCached(ctx context.Context, id string) (*sevallaapi.Pipeline, error) {
	cacheKey := "pipeline:" + id
	
	// Check cache first
	if cached, found := poc.cache.Get(cacheKey); found {
		tflog.Debug(ctx, "Pipeline retrieved from cache", map[string]interface{}{"id": id})
		if pipeline, ok := cached.(*sevallaapi.Pipeline); ok {
			return pipeline, nil
		}
	}
	
	// Wait for rate limiter
	if err := poc.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	
	// Make API call
	tflog.Debug(ctx, "Making API call for pipeline", map[string]interface{}{"id": id})
	pipeline, err := sevallaapi.NewPipelineService(poc.client).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	poc.cache.Set(cacheKey, pipeline, 5*time.Minute)
	
	return pipeline, nil
}

// InvalidateCache invalidates cache entries for a specific resource type.
func (poc *PerformanceOptimizedClient) InvalidateCache(resourceType, id string) {
	cacheKey := resourceType + ":" + id
	poc.cache.mutex.Lock()
	defer poc.cache.mutex.Unlock()
	
	delete(poc.cache.cache, cacheKey)
}

// ClearCache clears all cache entries.
func (poc *PerformanceOptimizedClient) ClearCache() {
	poc.cache.Clear()
}

// Stop stops all performance optimization components.
func (poc *PerformanceOptimizedClient) Stop() {
	poc.rateLimiter.Stop()
	poc.cache.Clear()
}