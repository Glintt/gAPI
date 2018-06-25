package cache

import (
	"time"

	"github.com/allegro/bigcache"
)

type Cache struct {
	ServiceDiscovery *bigcache.BigCache
	Apis             *bigcache.BigCache
	OAuth            *bigcache.BigCache
}

var GatewayCache Cache

func InitCachingService() {
	sdCacheConfig := bigcache.Config{
		Shards:             1024,
		LifeWindow:         1 * time.Hour,
		CleanWindow:        1 * time.Hour,
		MaxEntriesInWindow: 50,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           RemoveCache,
	}
	GatewayCache.ServiceDiscovery, _ = bigcache.NewBigCache(sdCacheConfig)

	apisCacheConfig := bigcache.Config{
		Shards:             1024,
		LifeWindow:         2 * time.Second,
		CleanWindow:        2 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           RemoveCache,
	}
	GatewayCache.Apis, _ = bigcache.NewBigCache(apisCacheConfig)

	oauthCacheConfig := bigcache.Config{
		Shards:             1024,
		LifeWindow:         30 * time.Minute,
		CleanWindow:        30 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           RemoveCache,
	}
	GatewayCache.OAuth, _ = bigcache.NewBigCache(oauthCacheConfig)
}
func InvalidateCache() {
	GatewayCache.ServiceDiscovery.Reset()
	GatewayCache.Apis.Reset()
	GatewayCache.OAuth.Reset()
}
func RemoveCache(key string, value []byte) {

}
func ServiceDiscoveryCacheStore(key string, value []byte) {
	GatewayCache.ServiceDiscovery.Set(key, value)
}

func ApisCacheStore(key string, value []byte) {
	GatewayCache.Apis.Set(key, value)
}

func OAuthCacheStore(key string, value []byte) {
	GatewayCache.OAuth.Set(key, value)
}

func ServiceDiscoveryCacheGet(key string) ([]byte, error) {
	return GatewayCache.ServiceDiscovery.Get(key)
}

func ApisCacheGet(key string) ([]byte, error) {
	return GatewayCache.Apis.Get(key)
}

func OAuthCacheGet(key string) ([]byte, error) {
	return GatewayCache.OAuth.Get(key)
}
