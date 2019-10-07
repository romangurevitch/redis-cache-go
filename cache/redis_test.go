package cache

import (
	"github.com/alicebob/miniredis"
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/test"
	"reflect"
	"testing"
)

func TestStoreLoad(t *testing.T) {
	miniRedis := getMiniRedis(t)
	defer miniRedis.Close()
	redisCache := getRedisCache(t, miniRedis)
	defer redisCache.Close()

	expected := []byte("expected value")
	err := redisCache.Store("key", expected)
	test.CheckError(t, err)

	actual, err := redisCache.Load("key")
	test.CheckError(t, err)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Load() = %v, want %v", actual, expected)
	}
}

func TestMultiStoreLoad(t *testing.T) {
	miniRedis := getMiniRedis(t)
	defer miniRedis.Close()
	redisCache := getRedisCache(t, miniRedis)
	defer redisCache.Close()

	expected := []byte("expected value")
	const key = "key"
	err := redisCache.Store(key, expected)
	test.CheckError(t, err)

	err = redisCache.Store(key, expected)
	test.CheckError(t, err)

	actual, err := redisCache.Load(key)
	test.CheckError(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Load() = %v, want %v", actual, expected)
	}
}

func TestEmptyLoad(t *testing.T) {
	miniRedis := getMiniRedis(t)
	defer miniRedis.Close()
	redisCache := getRedisCache(t, miniRedis)
	defer redisCache.Close()

	const key = "key"
	actual, err := redisCache.Load(key)
	test.CheckError(t, err)
	if actual != nil {
		t.Errorf("Load() = %v, want %v", actual, nil)
	}
}

func TestInvalidate(t *testing.T) {
	miniRedis := getMiniRedis(t)
	defer miniRedis.Close()
	redisCache := getRedisCache(t, miniRedis)
	defer redisCache.Close()

	const key = "key"
	err := redisCache.Store(key, []byte("expected value"))
	test.CheckError(t, err)

	err = redisCache.Invalidate()
	test.CheckError(t, err)

	actual, err := redisCache.Load(key)
	test.CheckError(t, err)

	if actual != nil {
		t.Errorf("Load() = %v, want %v", actual, nil)
	}
}

func getMiniRedis(t *testing.T) *miniredis.Miniredis {
	miniRedis, err := miniredis.Run()
	test.CheckError(t, err)
	return miniRedis
}

func getRedisCache(t *testing.T, redis *miniredis.Miniredis) *redisCache {
	redisCache, err := NewRedis("tcp", redis.Addr(), config.RedisPoolSize)
	test.CheckError(t, err)
	return redisCache
}
