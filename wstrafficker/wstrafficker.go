// Package wstrafficker provides smart ws:// client to connect to ResourceD Master
package wstrafficker

import (
	"math"
	"net"
	"net/http"
	"net/url"
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
)

func NewClient(originURL, targetURL string, settings map[string]interface{}) (c *websocket.Conn, response *http.Response, err error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, nil, err
	}

	dialer := net.Dialer{}

	timeoutInterface, ok := settings["Timeout"]
	if ok {
		dialer.Timeout = timeoutInterface.(time.Duration)
	}

	keepAliveInterface, ok := settings["KeepAlive"]
	if ok {
		dialer.KeepAlive = keepAliveInterface.(time.Duration)
	}

	rawConn, err := dialer.Dial("tcp", u.Host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Errorf("TCP connection error when dialing %v", u.Host)

		return nil, nil, err
	}

	wsHeaders := http.Header{
		"Origin": {originURL},
		// your milage may differ
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	return websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
}

func NewWSTrafficker(originURL, targetURL string, settings map[string]interface{}) (*WSTrafficker, error) {
	client, _, err := NewClient(originURL, targetURL, settings)
	if err != nil {
		return nil, err
	}

	ws := &WSTrafficker{}
	ws.Chans.Send = make(chan []byte)
	ws.Chans.Receive = make(chan []byte)

	ws.OriginURL = originURL
	ws.TargetURL = targetURL
	ws.ClientSettings = settings
	ws.Client = client

	return ws, nil
}

type WSTrafficker struct {
	Chans struct {
		Send    chan []byte
		Receive chan []byte
	}
	OriginURL        string
	TargetURL        string
	ClientSettings   map[string]interface{}
	Client           *websocket.Conn
	ReconnectBackoff int64
}

func (ws *WSTrafficker) PingAndReconnect() {
	go func() {
		for {
			err := ws.Client.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(pingPeriod))
			if err == nil {
				retryInterval := math.Pow(2, float64(ws.ReconnectBackoff))
				time.Sleep(time.Duration(retryInterval) * time.Second)
				continue
			}

			// Disconnect old client
			ws.Client.Close()

			// Create a new client
			client, _, err := NewClient(ws.OriginURL, ws.TargetURL, ws.ClientSettings)
			if err != nil {
				ws.ReconnectBackoff += 1

				logrus.WithFields(logrus.Fields{
					"Error":            err.Error(),
					"ReconnectBackoff": ws.ReconnectBackoff,
					"OriginURL":        ws.OriginURL,
					"TargetURL":        ws.TargetURL,
					"Timeout":          ws.ClientSettings["Timeout"],
				}).Error("Failed to create new websocket client during reconnect")

			} else {
				ws.Client = client
				ws.ReconnectBackoff = 0

				logrus.WithFields(logrus.Fields{
					"OriginURL": ws.OriginURL,
					"TargetURL": ws.TargetURL,
					"Timeout":   ws.ClientSettings["Timeout"],
				}).Info("Created new websocket client during reconnect")
			}

			retryInterval := math.Pow(2, float64(ws.ReconnectBackoff))
			time.Sleep(time.Duration(retryInterval) * time.Second)
		}
	}()
}

func (ws *WSTrafficker) Write(messageType int, payload []byte) error {
	ws.Client.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.Client.WriteMessage(messageType, payload)
}
