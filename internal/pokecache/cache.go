package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]cacheEntry
	Mux      sync.RWMutex
	Interval time.Duration
}

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	c.Entries[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mux.RLock()
	defer c.Mux.RUnlock()
	entry, ok := c.Entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.Val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	for range ticker.C {
		c.Mux.Lock()
		for k, v := range c.Entries {
			if time.Since(v.CreatedAt) > c.Interval {
				delete(c.Entries, k)
			}
		}
		c.Mux.Unlock()
	}

}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		Entries:  map[string]cacheEntry{},
		Interval: interval,
	}

	go c.reapLoop()

	return c
}
