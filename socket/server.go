package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func New(maxSocket int) *Server {
	result := &Server{
		Sockets: map[string]*Socket{},
		Limit:   maxSocket,
		Locker:  &sync.Mutex{},

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
	if result.Limit <= 0 {
		result.Limit = 50000
	}

	result.upgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost"
	}
	return result
}

func (s *Server) Connect(id string, w http.ResponseWriter, r *http.Request) *Socket {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	s.Locker.Lock()
	if _, ok := s.Sockets[id]; !ok && len(s.Sockets) < s.Limit {
		s.Sockets[id] = newSocket(id, conn)
	}
	s.Locker.Unlock()
	return s.Sockets[id]
}

func (s *Server) SendMessage(id string, message []byte) {
	if _, ok := s.Sockets[id]; ok {
		s.Sockets[id].sendMessage <- message
	}
}
