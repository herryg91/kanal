package socket

import (
	"sync"
	"time"

	chRedis "github.com/herryg91/kanal/channel/redis"

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

	channelEngine *chRedis.RedisChEngine

	LastUpdate  time.Time
	sendMessage chan []byte

	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int
}
