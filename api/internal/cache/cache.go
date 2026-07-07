package cache

import (
	"sync"
	"time"
)

type Item struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]*Item
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*Item),
		ttl:   ttl,
	}
	go c.evictLoop()
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &Item{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) evictLoop() {
	for {
		time.Sleep(c.ttl)
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.ExpiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
