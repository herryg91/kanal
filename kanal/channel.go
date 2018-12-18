package kanal

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func New(name string) *Channel {
	result := &Channel{
		name:    name,
		sockets: map[string]Socket{},
		mutex:   &sync.RWMutex{},
	}

	return result
}

func (c *Channel) JoinChannel(id string, conn *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn.SetCloseHandler(func(code int, text string) error {
		c.LeftChannel(id)
		return nil
	})

	c.sockets[id] = newSocket(id, conn)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(fmt.Printf("error: %v", err))
				}
				break
			}
			message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
			c.BroadcastMessage(string(message))
		}
	}()
}

func (c *Channel) LeftChannel(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.sockets, id)
}

func (c *Channel) CountEstablishedConnection() int {
	return len(c.sockets)
}

func (s *Channel) GetConnections() map[string]Socket {
	return s.sockets
}

func (s *Channel) GetConnection(id string) Socket {
	return s.sockets[id]
}

func (c *Channel) BroadcastMessage(message string) {
	for _, socket := range c.sockets {
		socket.sendMessage <- []byte(message)
	}
}
