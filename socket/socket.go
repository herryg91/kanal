package socket

import (
	"bytes"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

func newSocket(name string, conn *websocket.Conn) *Socket {
	result := &Socket{
		Name: name,
		Conn: conn,

		redisClient:   nil,
		channelEngine: nil,

		LastUpdate:     time.Now(),
		sendMessage:    make(chan []byte),
		writeWait:      10 * time.Second,
		pongWait:       60 * time.Second,
		pingPeriod:     ((60 * time.Second) * 9) / 10,
		maxMessageSize: 512,
	}

	result.redisClient = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go result.subscribeChannel([]string{"test"})
	go result.readPump()
	go result.writePump()

	result.Conn.SetCloseHandler(func(code int, text string) error {
		result.channelEngine.Close()
		return nil
	})
	return result
}

func (s *Socket) subscribeChannel(channelName []string) {
	s.channelEngine = s.redisClient.Subscribe(channelName...)
	for msg := range s.channelEngine.Channel() {
		// channelName := msg.Channel
		payload := msg.Payload
		s.sendMessage <- []byte(payload)
	}
	log.Println("subscribe channel stopped")
}

func (s *Socket) changeChannel(channelName []string) {
	s.channelEngine.Close()
	go s.subscribeChannel(channelName)
}

func (s *Socket) readPump() {
	defer func() {
		s.Conn.Close()
	}()

	s.Conn.SetReadLimit(512)
	s.Conn.SetReadDeadline(time.Now().Add(s.pongWait))
	s.Conn.SetPongHandler(func(string) error { s.Conn.SetReadDeadline(time.Now().Add(s.pongWait)); return nil })
	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		// s.sendMessage <- message
		s.redisClient.Publish("test", message)
	}
}

func (s *Socket) writePump() {
	ticker := time.NewTicker(s.pingPeriod)
	defer func() {
		ticker.Stop()
		s.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-s.sendMessage:
			s.Conn.SetWriteDeadline(time.Now().Add(s.writeWait))
			if !ok {
				s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(s.sendMessage)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-s.sendMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			s.Conn.SetWriteDeadline(time.Now().Add(s.writeWait))
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
