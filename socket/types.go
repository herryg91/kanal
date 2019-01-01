package socket

import (
	"sync"
	"time"

	"github.com/go-redis/redis"

	"github.com/gorilla/websocket"
)

type Server struct {
	Sockets map[string]*Socket
	Limit   int
	Locker  *sync.Mutex

	upgrader websocket.Upgrader
}

type Socket struct {
	Name string
	Conn *websocket.Conn

	redisClient   *redis.Client
	channelEngine *redis.PubSub

	LastUpdate  time.Time
	sendMessage chan []byte

	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int
}
