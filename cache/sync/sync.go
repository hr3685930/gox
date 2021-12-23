package sync

import (
	"context"
	"errors"
	"sync"
	"time"
)

type (
	syncMapItem struct {
		data     string
		duration int64
	}

	syncMap struct {
		storage *sync.Map
	}
)

type err string

// Error returns the string errors value.
func (e err) Error() string {
	return string(e)
}

const (
	// ErrCacheExpired returns an errors when the cache key was expired.
	ErrCacheExpired = err("cache expired")

	// ErrFlush returns an errors when flush fails.
	ErrFlush = err("unable to flush")

	// ErrSave returns an errors when save fails.
	ErrSave = err("unable to save")

	// ErrDelete returns an errors when deletion fails.
	ErrDelete = err("unable to delete")

	// ErrDecode returns an errors when decode fails.
	ErrDecode = err("unable to decode")
)

// New creates an instance of SyncMap cache driver
func New() *syncMap {
	return &syncMap{&sync.Map{}}
}

func (sm *syncMap) read(ctx context.Context, key string) (*syncMapItem, error) {
	v, ok := sm.storage.Load(key)
	if !ok {
		return nil, errors.New("key not found")
	}

	item := v.(*syncMapItem)

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = sm.Delete(ctx, key)
		return nil, ErrCacheExpired
	}

	return item, nil
}

// Contains checks if cached key exists in SyncMap storage
func (sm *syncMap) Contains(ctx context.Context, key string) bool {
	_, err := sm.Fetch(ctx, key)
	return err == nil
}

// Delete the cached key from SyncMap storage
func (sm *syncMap) Delete(ctx context.Context, key string) error {
	sm.storage.Delete(key)
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap) Fetch(ctx context.Context, key string) (string, error) {
	item, err := sm.read(ctx, key)
	if err != nil {
		return "", err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap) FetchMulti(ctx context.Context, keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := sm.Fetch(ctx, key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap) Flush(ctx context.Context) error {
	sm.storage = &sync.Map{}
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap) Save(ctx context.Context, key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	sm.storage.Store(key, &syncMapItem{value, duration})
	return nil
}
