package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Sab94/udsws"
	"github.com/gorilla/websocket"
)
var socketPath = "/tmp/udsws.sock"

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, err := udsws.NewClient(socketPath, "/ws")
	if err != nil {
		log.Fatalf("NewClient failed : %s", err.Error())
	}

	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Fatal("write:", err)
				return
			}
		case <-interrupt:
			log.Println("Got interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatal("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
