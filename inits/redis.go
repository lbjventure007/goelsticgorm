package inits

import (
	redis "github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func Script(str string) *redis.Script {
	return redis.NewScript(str)
}
