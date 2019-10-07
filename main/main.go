package main

import (
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/cache"
	"github.com/romangurevitch/redis-cache-go/contact"
	"log"
)

func main() {
	redisCache, err := cache.NewRedis("tcp", config.RedisHost+":"+config.RedisPort, config.RedisPoolSize)
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}
	defer redisCache.Close()

	cachedServer, err := contact.NewContactServer(config.ApiBaseUrl, redisCache)
	if err != nil {
		log.Fatalf("could not create contact server: %v", err)
	}
	cachedServer.Start()
}
