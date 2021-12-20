package cache

import (
    "sync"
    "time"
)

var CacheMap sync.Map

var Cached Cache

func GetCache(c string) Cache {
    v, ok := CacheMap.Load(c);
    if ok {
        return v.(Cache)
    }
    return nil
}

type Cache interface {
    // Contains check if a cached key exists
    Contains(key string) bool

    // Delete remove the cached key
    Delete(key string) error

    // Fetch retrieve the cached key value
    Fetch(key string) (string, error)

    // FetchMulti retrieve multiple cached keys value
    FetchMulti(keys []string) map[string]string

    // Flush remove all cached keys
    Flush() error

    // Save cache a value by key
    Save(key string, value string, lifeTime time.Duration) error
}

