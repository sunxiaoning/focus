package cachetype

import (
	"sync"
	"time"
)

type CachedData interface {
}
type expiredCache struct {
	Timestamp time.Time
	Cache     CachedData
}

var caches = &sync.Map{}

func SetExpiredCache(cacheName string, cacheKey string, cacheData CachedData, duration time.Duration) {
	cache, _ := caches.LoadOrStore(cacheName, &sync.Map{})
	cache.(*sync.Map).Store(cacheKey, &expiredCache{
		Timestamp: time.Now().Add(duration),
		Cache:     cacheData,
	})
}

func GetCache(cacheName string, cacheKey string) CachedData {
	cache, ok := caches.Load(cacheName)
	if !ok || cache == nil {
		return nil
	}
	cacheData, ok := cache.(*sync.Map).Load(cacheKey)
	if !ok || cache == nil {
		return nil
	}
	expiredCache := cacheData.(*expiredCache)
	if expiredCache.Timestamp.Before(time.Now()) {
		cache.(*sync.Map).Delete(cacheKey)
		return nil
	}
	return expiredCache.Cache
}
