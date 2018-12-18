package kanal

import (
	"time"

	"github.com/gorilla/websocket"
)

func newSocket(name string, conn *websocket.Conn) Socket {
	result := Socket{
		Name:           name,
		Conn:           conn,
		LastUpdate:     time.Now(),
		sendMessage:    make(chan []byte),
		writeWait:      10 * time.Second,
		pongWait:       60 * time.Second,
		pingPeriod:     ((60 * time.Second) * 9) / 10,
		maxMessageSize: 512,
	}

	go func() {
		ticker := time.NewTicker(result.pingPeriod)
		defer func() {
			ticker.Stop()
			conn.Close()
		}()
		for {
			select {
			case message, ok := <-result.sendMessage:
				conn.SetWriteDeadline(time.Now().Add(result.writeWait))
				if !ok {
					// The hub closed the channel.
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(message)

				// Add queued chat messages to the current websocket message.
				n := len(result.sendMessage)
				for i := 0; i < n; i++ {
					w.Write([]byte{'\n'})
					w.Write(<-result.sendMessage)
				}

				if err := w.Close(); err != nil {
					return
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(result.writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
	return result
}
