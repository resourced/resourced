package wstrafficker

import (
	"github.com/gorilla/websocket"
)

func New(client *websocket.Conn) *WSTrafficker {
	ws := &WSTrafficker{}
	ws.Chans.Send = make(chan []byte)
	ws.Chans.Receive = make(chan []byte)

	ws.Client = client

	return ws
}

type WSTrafficker struct {
	Chans struct {
		Send    chan []byte
		Receive chan []byte
	}
	Client *websocket.Conn
}
