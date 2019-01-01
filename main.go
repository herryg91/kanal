package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/herryg91/kanal/socket"
	"github.com/julienschmidt/httprouter"
)

var socketServer *socket.Server

var port int

func flagInit() {
	flag.IntVar(&port, "port", 9000, "")
}
func main() {
	flagInit()
	flag.Parse()
	socketServer = socket.New(50000)
	router := httprouter.New()

	router.GET("/connect", wsConnect)

	fmt.Println("Listen and Serve: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))

	// grace.Serve(":9000", router)
}

func wsConnect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("id") == "" {
		return
	}
	socketServer.Connect(r.FormValue("id"), w, r)

}
