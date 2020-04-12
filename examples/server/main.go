package main

import (
	"log"
	"os"

	"github.com/Sab94/udsws"
	"github.com/gorilla/websocket"
)

var socketPath = "/tmp/udsws.sock"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	os.Remove(socketPath)
	udsws.NewHandleFunc("/ws", upgrader, reader)
	err := udsws.ListenAndServe(socketPath, nil)
	if err != nil {
		log.Fatalf("WS server failed : %s", err.Error())
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Got message : ", string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}