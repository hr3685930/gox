package redis

import (
    "net"
    "time"

    rd "gopkg.in/redis.v4"
)

type Redis struct {
    rd.BaseCmdable
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
func (r *Redis) Contains(key string) bool {
    status, _ := r.Exists(key).Result()
    return status
}

// Delete the cached key from Redis storage
func (r *Redis) Delete(key string) error {
    return r.Del(key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *Redis) Fetch(key string) (string, error) {
    return r.Get(key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *Redis) FetchMulti(keys []string) map[string]string {
    result := make(map[string]string)

    items, err := r.MGet(keys...).Result()
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
func (r *Redis) Flush() error {
    return r.FlushAll().Err()
}

// Save a value in Redis storage by key
func (r *Redis) Save(key string, value string, lifeTime time.Duration) error {
    return r.Set(key, value, lifeTime).Err()
}
