package cache

import (
	"github.com/mediocregopher/radix/v3"
)

type Cache interface {
	// Store key and value
	Store(key string, val []byte) error
	// Load value by key
	Load(key string) ([]byte, error)
	// Invalidate key
	Invalidate() error

	Close() error
}

func NewRedis(network, addr string, size int) (*redisCache, error) {
	pool, err := radix.NewPool(network, addr, size)
	if err != nil {
		return nil, err
	}
	return &redisCache{radix: pool}, nil
}

type redisCache struct {
	radix *radix.Pool
}

func (r *redisCache) Store(key string, val []byte) error {
	return r.radix.Do(radix.Cmd(nil, "SET", key, string(val)))
}

func (r *redisCache) Load(key string) ([]byte, error) {
	var val []byte
	err := r.radix.Do(radix.Cmd(&val, "GET", key))
	return val, err
}

func (r *redisCache) Invalidate() error {
	return r.radix.Do(radix.Cmd(nil, "FLUSHDB"))
}

func (r *redisCache) Close() error {
	return r.radix.Close()
}
