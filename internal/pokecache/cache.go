package pokecache

import (
	"sync"
	"time"
)

// cacheEntry holds the cached data and its creation time.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is a threadsafe structure for storing cached entries.
type Cache struct {
	mu       sync.Mutex
	store    map[string]cacheEntry
	interval time.Duration
}

// NewCache creates and returns a new Cache instance.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		store:    make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

// Set adds or updates an entry in the cache.
func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

// Get retrieves an entry by key, returning the value and whether it was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.store[key]
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.store {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.store, key)
			}
		}
		c.mu.Unlock()
	}
}
