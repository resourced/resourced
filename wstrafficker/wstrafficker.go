package wstrafficker

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func NewWSTrafficker(client *websocket.Conn) *WSTrafficker {
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

func (ws *WSTrafficker) Write(messageType int, payload []byte) error {
	ws.Client.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.Client.WriteMessage(messageType, payload)
}

func NewWSTraffickers() *WSTraffickers {
	wss := &WSTraffickers{}
	wss.ByAccessTokenHostname = make(map[string]*WSTrafficker)

	return wss
}

type WSTraffickers struct {
	ByAccessTokenHostname map[string]*WSTrafficker
}

func (wss *WSTraffickers) SaveConnection(accessToken string, conn *websocket.Conn) (*WSTrafficker, error) {
	_, payloadJson, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	err = json.Unmarshal(payloadJson, &payload)
	if err != nil {
		return nil, err
	}

	hostnameInterface, hostnameExists := payload["Hostname"]
	if !hostnameExists {
		return nil, errors.New("Unable to establish websocket connection without Hostname payload")
	}

	wsTrafficker := NewWSTrafficker(conn)

	hostname := hostnameInterface.(string)
	wss.ByAccessTokenHostname[accessToken+"-"+hostname] = NewWSTrafficker(conn)

	logrus.WithFields(logrus.Fields{
		"AccessToken": accessToken,
		"Hostname":    hostname,
	}).Info("Saved websocket connection in memory")

	return wsTrafficker, nil
}

func (wss *WSTraffickers) GetConnectionByAccessTokenHostname(accessToken, hostname string) *WSTrafficker {
	return wss.ByAccessTokenHostname[accessToken+"-"+hostname]
}
