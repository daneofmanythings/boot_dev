package internal

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		mu:      sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}
	c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if value, ok := c.entries[key]; !ok {
		return []byte{}, false
	} else {
		return value.val, true
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			go func() {
				for k, v := range c.entries {
					// TODO: verify this is working
					if time.Now().Nanosecond()-v.createdAt.Nanosecond() > int(interval.Nanoseconds()) {
						fmt.Printf("deleting %v from cache", k)
						c.mu.Lock()
						delete(c.entries, k)
						c.mu.Unlock()
					}
				}
			}()
		}
	}()
}
