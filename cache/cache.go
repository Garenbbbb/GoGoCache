package cache

import (
	"gogocache/cache_type"
	"sync"
)

type Cache struct {
	mu       sync.Mutex
	cache    *cache_type.LRU
	capacity int64
}

func newCache(capacity int64) *Cache {
	return &Cache{capacity: capacity}
}

func (cache *Cache) Add(key string, value ByteView) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	//lazy initialization
	if cache.cache == nil {
		cache.cache = cache_type.NewLRU(cache.capacity, nil)
	}
	cache.cache.Add(key, value)
}

func (cache *Cache) Get(key string) (ByteView, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.cache == nil {
		return ByteView{}, false
	}

	if v, flag := cache.cache.Get(key); flag {
		return v.(ByteView), flag
	}

	return ByteView{}, false
}
