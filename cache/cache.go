package cache

import (
    "context"
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
    Contains(ctx context.Context, key string) bool

    // Delete remove the cached key
    Delete(ctx context.Context, key string) error

    // Fetch retrieve the cached key value
    Fetch(ctx context.Context, key string) (string, error)

    // FetchMulti retrieve multiple cached keys value
    FetchMulti(ctx context.Context, keys []string) map[string]string

    // Flush remove all cached keys
    Flush(ctx context.Context) error

    // Save cache a value by key
    Save(ctx context.Context, key string, value string, lifeTime time.Duration) error
}

