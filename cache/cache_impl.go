package cache

import (
	"context"

	"time"

	rCache "github.com/go-redis/cache/v8"
)

type CacheImpl interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error)
	GetValue(ctx context.Context, key string) (string, error)
	DeleteValue(ctx context.Context, key string) error
}

type CacheRepositoryImpl struct {
	cache *Cache
}

func New(db *Cache) CacheImpl {
	return &CacheRepositoryImpl{
		cache: db,
	}
}

func (r *CacheRepositoryImpl) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	err = r.cache.Set(&rCache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})

	if err != nil {
		return
	}

	return
}

func (r *CacheRepositoryImpl) GetValue(ctx context.Context, key string) (value string, err error) {
	if err = r.cache.Get(ctx, key, &value); err != nil {
		return
	}

	return
}

func (r *CacheRepositoryImpl) DeleteValue(ctx context.Context, key string) (err error) {
	if err = r.cache.Delete(ctx, key); err != nil {
		return
	}

	return
}
