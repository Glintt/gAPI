package cache

import (
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"strings"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"

	"github.com/allegro/bigcache"
)

type Cache struct {
	ServiceDiscovery *bigcache.BigCache
	Apis             *bigcache.BigCache
	OAuth            *bigcache.BigCache
	gAPIApi          *bigcache.BigCache
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

	if config.GApiConfiguration.Cache.Enabled {
		gAPIApiConfig := bigcache.Config{
			Shards:             1024,
			LifeWindow:         2 * time.Second,
			CleanWindow:        2 * time.Second,
			MaxEntriesInWindow: 1000 * 10 * 60,
			Verbose:            true,
			HardMaxCacheSize:   8192,
			OnRemove:           RemoveCache,
		}
		GatewayCache.gAPIApi, _ = bigcache.NewBigCache(gAPIApiConfig)
	}
}
func InvalidateCache() {
	GatewayCache.ServiceDiscovery.Reset()
	GatewayCache.Apis.Reset()
	GatewayCache.OAuth.Reset()
	GatewayCache.gAPIApi.Reset()
}
func RemoveCache(key string, value []byte) {

}

func GApiCacheStore(key string, value []byte) {
	GatewayCache.gAPIApi.Set(key, value)
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

func GApiCacheGet(key string) ([]byte, error) {
	return GatewayCache.gAPIApi.Get(key)
}

func GApiCacheKey(c *routing.Context) string {
	return string(c.Method()) + "-" + c.URI().String() + c.Request.Header.String()
}

func StoreCacheGApi(c *routing.Context) error {

	if config.GApiConfiguration.Cache.Enabled && strings.ToLower(string(c.Method())) == "get" {
		GApiCacheStore(GApiCacheKey(c), c.Response.Body())
	}
	return nil
}

func ResponseCacheGApi(c *routing.Context) error {

	if !config.GApiConfiguration.Cache.Enabled {
		return nil
	}

	resp, err := GApiCacheGet(GApiCacheKey(c))

	if err != nil {
		return nil
	}

	http.Response(c, string(resp), 200, constants.SERVICE_NAME, config.APPLICATION_JSON)
	c.Abort()
	return nil
}
