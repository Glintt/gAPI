package cacheable

import (
	"encoding/json"
	"github.com/Glintt/gAPI/api/utils"
	"github.com/allegro/bigcache"
	"time"
)

var genericCache = CacheableStorage{}

type CacheableStorage struct {
	cacheStorage *bigcache.BigCache
}

func GetCacheableStorageInstance() CacheableStorage {
	genericCache.InitializeCacheStorage()
	return genericCache
}

func (cs *CacheableStorage) Get(key string, out interface{}) error{
	var obj interface{}
	err := getFromCache(cs.cacheStorage, key, &obj)
	if err == nil {
		utils.LogMessage("CACHEABLE --- Get: " + key, utils.DebugLogType)
		jsonData, _ := json.Marshal(obj)
		json.Unmarshal(jsonData, &out)
		return err
	}
	return err
}

func (cs *CacheableStorage) Set(key string, value interface{}) {
	utils.LogMessage("CACHEABLE --- Set: " + key, utils.DebugLogType)
	addToCache(cs.cacheStorage, key, value)
}

func (cs *CacheableStorage) Delete(key string) error {
	utils.LogMessage("CACHEABLE --- Delete: " + key, utils.DebugLogType)
	return cs.cacheStorage.Delete(key)
}

func (cs *CacheableStorage) Reset() error {
	utils.LogMessage("CACHEABLE --- Reset", utils.DebugLogType)
	return cs.cacheStorage.Reset()
}

func (cs *CacheableStorage) Cacheable(key string, f func() (interface{}, error), out interface{}) error{	
	if (cs.cacheStorage == nil) {
		cs.InitializeCacheStorage()
	}

	var obj interface{}
	// Check on cache and return if found
	err := cs.Get(key, &out)
	if err == nil {
		return err
	}

	obj, err = f()
	if err == nil {
		cs.Set(key, obj)
	}
	
	if err == nil {
		jsonData, _ := json.Marshal(obj)
		json.Unmarshal(jsonData, &out)
	}

	return err
}

func (cs *CacheableStorage) InitializeCacheStorage() {
	if cs.cacheStorage == nil {
		cs.cacheStorage,_ = bigcache.NewBigCache(bigcache.Config{
			Shards:             1024,
			LifeWindow:         24 * time.Hour,
			CleanWindow:        24 * time.Hour,
			MaxEntriesInWindow: 1000 * 10 * 60,
			Verbose:            true,
			HardMaxCacheSize:   8192,
			//OnRemove:           RemoveCache,
		})
	}
}
func getFromCache(cacheStorage *bigcache.BigCache, key string, out interface{}) error {
	valueByte, err := cacheStorage.Get(key)
	if err != nil { return err }
	return json.Unmarshal(valueByte, out)
}

func addToCache(cacheStorage *bigcache.BigCache, key string, value interface{}) {
	valueByte, _ := json.Marshal(value)
	cacheStorage.Set(key, valueByte)
}
