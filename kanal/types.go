package kanal

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Channel struct {
	name    string
	sockets map[string]Socket
	mutex   *sync.RWMutex
}

type Socket struct {
	Name        string
	Conn        *websocket.Conn
	LastUpdate  time.Time
	sendMessage chan []byte

	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int
}

// func (s Socket) Read() {
// 	for {
// 		_, message, err := s.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Println(fmt.Printf("error: %v", err))
// 			}
// 			break
// 		}
// 		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
// 		s.PushWriteAll(string(message))
// 	}
// }

// func newSocket(name string, conn *websocket.Conn) *Socket {
// 	result := Socket{
// 		Name:       name,
// 		Conn:       conn,
// 		LastUpdate: time.Now(),
// 	}
// 	return result
// }

// conn.SetCloseHandler(func(code int, text string) error {
// 	s.Left(id)
// 	return nil
// })

// go func() {
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Println(fmt.Printf("error: %v", err))
// 			}
// 			break
// 		}
// 		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
// 		s.PushWriteAll(string(message))
// 	}
// }()

// func (s *Channel) GetConnList() map[string]*websocket.Conn {
// 	return s.connList
// }

// func (s *Channel) GetConn(id string) *websocket.Conn {
// 	return s.connList[id]
// }

// func (s *Channel) Join(id string, conn *websocket.Conn) {
// 	s.connList[id] = conn
// 	conn.SetCloseHandler(func(code int, text string) error {
// 		s.Left(id)
// 		return nil
// 	})

// 	go func() {
// 		for {
// 			_, message, err := conn.ReadMessage()
// 			if err != nil {
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 					log.Println(fmt.Printf("error: %v", err))
// 				}
// 				break
// 			}
// 			message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
// 			s.PushWriteAll(string(message))
// 		}
// 	}()
// }

// func (s *Channel) PushWriteAll(message string) {
// 	for _, userConn := range s.connList {
// 		userConn.WriteMessage(websocket.TextMessage, []byte(message))
// 	}
// }

// func (s *Channel) PushWrite(uid, message string) {
// 	userConn := s.GetConn(uid)
// 	if userConn == nil {
// 		return
// 	}

// 	writer, err := userConn.NextWriter(websocket.TextMessage)
// 	if err != nil {
// 		return
// 	}
// 	writer.Write([]byte(message))
// 	writer.Close()
// }

// func (s *Channel) Left(id string) {
// 	delete(s.connList, id)
// }
