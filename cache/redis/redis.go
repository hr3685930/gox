package redis

import (
	"context"
	"net"
	"time"

	rd "github.com/go-redis/redis/v8"
)

type Redis struct {
	*rd.Client
}

// New creates an instance of Redis cache driver
func New(addr, port string, db int, pass string) (*Redis, error) {
	conn := rd.NewClient(&rd.Options{
		Addr:     addr + ":" + port,
		DB:       db,
		Password: pass,
	})

	if _, err := net.Dial("tcp", addr+":"+port); err != nil {
		return nil, err
	}
	return &Redis{conn}, nil
}

// Contains checks if cached key exists in Redis storage
func (r *Redis) Contains(ctx context.Context, key string) bool {
	status, _ := r.Exists(ctx, key).Result()
	if status > 0 {
		return true
	}
	return false
}

// Delete the cached key from Redis storage
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Del(ctx, key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *Redis) Fetch(ctx context.Context, key string) (string, error) {
	return r.Get(ctx, key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *Redis) FetchMulti(ctx context.Context, keys []string) map[string]string {
	result := make(map[string]string)

	items, err := r.MGet(ctx, keys...).Result()
	if err != nil {
		return result
	}

	for i := 0; i < len(keys); i++ {
		if items[i] != nil {
			result[keys[i]] = items[i].(string)
		}
	}

	return result
}

// Flush removes all cached keys of the Redis storage
func (r *Redis) Flush(ctx context.Context) error {
	return r.FlushAll(ctx).Err()
}

// Save a value in Redis storage by key
func (r *Redis) Save(ctx context.Context, key string, value string, lifeTime time.Duration) error {
	return r.Set(ctx, key, value, lifeTime).Err()
}

func (r *Redis) AddTracingHook() {
	r.AddHook(NewHook())
}
