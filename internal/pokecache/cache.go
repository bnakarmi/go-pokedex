package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]cacheEntry
	mux        *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntry: make(map[string]cacheEntry),
        mux: &sync.Mutex{},
	}

	go cache.reapLoop(interval)

    return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.cacheEntry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	val, ok := cache.cacheEntry[key]
	return val.val, ok
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cache.reap(time.Now().UTC(), interval)
	}
}

func (cache *Cache) reap(now time.Time, last time.Duration) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	for k, v := range cache.cacheEntry {
		if v.createdAt.Before(now.Add(-last)) {
			delete(cache.cacheEntry, k)
		}
	}
}
