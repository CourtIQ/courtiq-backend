package access

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem represents a cached access check result
type CacheItem struct {
	Result    *AccessResult
	ExpiresAt time.Time
}

// isExpired checks if a cache item has expired
func (c *CacheItem) isExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// AccessCache provides an in-memory cache for access check results
type AccessCache struct {
	mu       sync.RWMutex
	cache    map[string]*CacheItem
	maxItems int
	ttl      time.Duration
}

// NewAccessCache creates a new cache with specified size and TTL
func NewAccessCache(maxItems int, ttl time.Duration) *AccessCache {
	return &AccessCache{
		cache:    make(map[string]*CacheItem),
		maxItems: maxItems,
		ttl:      ttl,
	}
}

// generateKey creates a cache key from owner and viewer IDs plus config
func generateKey(ownerID, viewerID string, config CheckConfig) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s",
		ownerID,
		viewerID,
		config.RequiredLevel,
		config.EntityID,
		config.EntityType)
}

// Get retrieves a cached access result if available
func (c *AccessCache) Get(ownerID, viewerID string, config CheckConfig) *AccessResult {
	key := generateKey(ownerID, viewerID, config)
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.cache[key]
	if !exists {
		return nil
	}
	
	// If expired, return nil
	if item.isExpired() {
		return nil
	}
	
	return item.Result
}

// Set stores an access result in the cache
func (c *AccessCache) Set(ownerID, viewerID string, config CheckConfig, result *AccessResult) {
	key := generateKey(ownerID, viewerID, config)
	
	// Mark result as cached and set expiration
	result.CachedResult = true
	result.ExpiresAt = time.Now().Add(c.ttl)
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if we need to clean up (simple approach)
	if len(c.cache) >= c.maxItems {
		c.cleanupExpired()
		
		// If still at max, just return without caching
		if len(c.cache) >= c.maxItems {
			return
		}
	}
	
	c.cache[key] = &CacheItem{
		Result:    result,
		ExpiresAt: result.ExpiresAt,
	}
}

// Clear removes specific cache entries related to the given user IDs
func (c *AccessCache) Clear(userIDs ...string) {
	if len(userIDs) == 0 {
		return
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Simple approach: iterate through all keys
	// In a production system, you'd want a more efficient way to index cache entries by user
	for key := range c.cache {
		for _, userID := range userIDs {
			if key[:len(userID)] == userID || key[len(userID)+1:len(userID)*2+1] == userID {
				delete(c.cache, key)
				break
			}
		}
	}
}

// cleanupExpired removes expired entries from the cache
func (c *AccessCache) cleanupExpired() {
	now := time.Now()
	for key, item := range c.cache {
		if now.After(item.ExpiresAt) {
			delete(c.cache, key)
		}
	}
}