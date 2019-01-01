package chRedis

import (
	"log"
	"sync"

	"github.com/go-redis/redis"
)

func New(redisClient *redis.Client) *RedisChEngine {
	result := &RedisChEngine{
		redisClient:  redisClient,
		pubsub:       nil,
		channels:     map[string]bool{},
		locker:       &sync.Mutex{},
		ReceiveMsgFn: defaultReceiveMsgFn,
	}
	return result
}

func defaultReceiveMsgFn(channel string, pattern string, payload string) {

}

func (r *RedisChEngine) SetMessageReceive(fn ReceiveMessageFunc) {
	r.ReceiveMsgFn = fn
}

func (r *RedisChEngine) SendMessage(channel string, message string) {
	r.redisClient.Publish(channel, message).Result()
}

func (r *RedisChEngine) Join(name string) {
	if _, ok := r.channels[name]; !ok {
		r.channels[name] = true
		go r.refreshPubsub()
	}
}

func (r *RedisChEngine) Left(name string) {
	if _, ok := r.channels[name]; !ok {
		delete(r.channels, name)
		go r.refreshPubsub()
	}
}

func (r *RedisChEngine) Close() {
	if r.pubsub != nil {
		err := r.pubsub.Close()
		if err != nil {
			log.Println("[error]", err)
		}
		r.pubsub = nil
	}

	if r.redisClient != nil {
		err := r.redisClient.Close()
		if err != nil {
			log.Println("[error]", err)
		}
	}
	r.channels = map[string]bool{}
}

func (r *RedisChEngine) refreshPubsub() {
	if r.pubsub != nil {
		err := r.pubsub.Close()
		if err != nil {
			log.Println("[error]", err)
		}
	}

	channelNames := MapStringBoolToArrOfStr(r.channels)
	r.pubsub = r.redisClient.Subscribe(channelNames...)
	for msg := range r.pubsub.Channel() {
		r.ReceiveMsgFn(msg.Channel, msg.Pattern, msg.Payload)
	}
}
