package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	grace "gopkg.in/tokopedia/grace.v1"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	router := httprouter.New()

	router.GET("/connect", wsConnect)
	// router.GET("/push", wsPush)
	grace.Serve(":9000", router)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsConnect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
