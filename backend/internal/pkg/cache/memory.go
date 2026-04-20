package cache

import (
	"sync"
	"time"
)

type item struct {
	value     interface{}
	expiresAt time.Time
}

type MemoryCache struct {
	data map[string]*item
	mu   sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		data: make(map[string]*item),
	}
	go c.cleanup()
	return c
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.data[key]
	if !ok {
		return nil, false
	}
	if !v.expiresAt.IsZero() && time.Now().After(v.expiresAt) {
		return nil, false
	}
	return v.value, true
}

func (m *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	m.data[key] = &item{value: value, expiresAt: exp}
}

func (m *MemoryCache) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for k, v := range m.data {
			if !v.expiresAt.IsZero() && now.After(v.expiresAt) {
				delete(m.data, k)
			}
		}
		m.mu.Unlock()
	}
}
