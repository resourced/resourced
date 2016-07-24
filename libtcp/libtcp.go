package libtcp

import (
	"net"
	"time"

	"github.com/sethgrid/pester"
)

func NewConnectionWithRetries(addr string, maxRetries int) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	attempts := 0

	for {
		if err == nil {
			break
		}

		if err != nil && attempts > maxRetries {
			return conn, err
		}

		if err != nil {
			attempts = attempts + 1
			time.Sleep(pester.ExponentialJitterBackoff(attempts))
			conn, err = net.Dial("tcp", addr)
			continue
		}
	}

	return conn, err
}
