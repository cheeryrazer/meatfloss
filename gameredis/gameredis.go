package gameredis

import (
	"fmt"

	"meatfloss/config"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

var (
	// redisClient ...
	redisClient *redis.Client
)

// Initialize redis.
func Initialize() {
	addr := fmt.Sprint(config.Get().RedisServer.Host, ":", config.Get().RedisServer.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",                          // no password set
		DB:       config.Get().RedisServer.Db, // use default DB
		PoolSize: 64,                          // max connections
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		glog.Error("redisClient.Ping() failed, error: ", err)
	} else {
		glog.Info("redisClient.Ping() ok!")
	}

	return
}

// GetUniqueID ...
func GetUniqueID() int64 {
	result, err := redisClient.Incr("meatFlossUniqueID").Result()
	if err != nil {
		return 0
	}
	return result
}

// GetGoodsUniqueID ...
func GetGoodsUniqueID() int64 {
	result, err := redisClient.Incr("meatFlossUniqueID").Result()
	if err != nil {
		return 0
	}
	return result + 10000000
}

// structures
