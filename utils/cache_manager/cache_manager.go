package cache_manager

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

type CacheManager struct {
	cache *bigcache.BigCache
}

const (
	cacheDuration = 10 * time.Minute
)

var cacheManager *CacheManager

func newCacheManager() (*CacheManager, error) {
	cacheConfig := bigcache.Config{
		Shards:             1024,
		LifeWindow:         cacheDuration,
		CleanWindow:        cacheDuration,
		MaxEntriesInWindow: 10 * 60,
		MaxEntrySize:       500,
	}
	cache, err := bigcache.New(context.Background(), cacheConfig)
	if err != nil {
		return nil, err
	}

	return &CacheManager{cache: cache}, nil
}

func InitCacheManager() error {
	var err error
	cacheManager, err = newCacheManager()
	if err != nil {
		return err
	}

	return nil
}

func GetCacheManager() *CacheManager {
	return cacheManager
}

func (c *CacheManager) Set(key string, value []byte) error {
	return c.cache.Set(key, value)
}

func (c *CacheManager) Get(key string) ([]byte, error) {
	return c.cache.Get(key)
}
