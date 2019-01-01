package chRedis

import (
	"github.com/go-redis/redis"
)

type RedisChEngine struct {
	redisClient *redis.Client
	pubsub      *redis.PubSub
	channels    []string
}
