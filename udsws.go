package udsws

import (
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func ListenAndServe(sockPath string, handler http.Handler) error {
	if handler == nil {
		handler = http.DefaultServeMux
	}
	l, err := net.Listen("unix", sockPath)
	if err != nil {
		return err
	}
	return http.Serve(l, handler)
}

func NewHandleFunc(path string, upgrader websocket.Upgrader, reader func(conn *websocket.Conn)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("Unable to create handle function %s", err.Error())
		}
		reader(ws)
	})
}

func NewClient(sockPath, wsPath string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: "unix", Path: wsPath}
	dialer := &websocket.Dialer{
		NetDial:           func(_ ,_ string) (net.Conn, error) {
			return net.Dial("unix", sockPath)
		},
	}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}