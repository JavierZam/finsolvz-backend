package utils

import (
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// IsExpired checks if the cache item has expired
func (item CacheItem) IsExpired() bool {
	return time.Now().After(item.Expiration)
}

// Cache is a simple in-memory cache with expiration
type Cache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]CacheItem),
	}

	// Start cleanup goroutine
	go c.cleanup()

	return c
}

// Set adds an item to the cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if item.IsExpired() {
		// Remove expired item
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]CacheItem)
}

// cleanup removes expired items every minute
func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			for key, item := range c.items {
				if item.IsExpired() {
					delete(c.items, key)
				}
			}
			c.mutex.Unlock()
		}
	}
}

// Global cache instance
var globalCache *Cache

func init() {
	globalCache = NewCache()
}

// GetCache returns the global cache instance
func GetCache() *Cache {
	return globalCache
}
