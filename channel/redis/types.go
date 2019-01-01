package chRedis

import (
	"sync"

	"github.com/go-redis/redis"
)

type RedisChEngine struct {
	redisClient *redis.Client
	pubsub      *redis.PubSub
	channels    map[string]bool

	ReceiveMsgFn ReceiveMessageFunc
	locker       *sync.Mutex
}

type ReceiveMessageFunc func(channel string, pattern string, payload string)
