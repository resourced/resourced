package wsclient

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
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
			"err": err.Error(),
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
