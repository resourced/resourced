// Package mysql gathers MySQL related data from a host.
package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"sync"
	"time"
)

var connections map[string]*sqlx.DB
var connectionsLock = &sync.RWMutex{}

type Base struct {
	HostAndPort string
	Retries     int
}

func (m *Base) initConnection() error {
	var err error

	connectionsLock.Lock()
	if connections == nil {
		connections = make(map[string]*sqlx.DB)
	}
	connectionsLock.Unlock()

	if m.Retries <= 0 {
		m.Retries = 10
	}

	for i := 0; i < m.Retries; i++ {
		// Do not create connection if one already exist.
		connectionsLock.RLock()
		conn, ok := connections[m.HostAndPort]
		connectionsLock.RUnlock()

		if !ok {
			newConn, err := sqlx.Open("mysql", fmt.Sprintf("root:@(%v)/?parseTime=true", m.HostAndPort))
			if err == nil && newConn != nil {
				connectionsLock.Lock()
				connections[m.HostAndPort] = newConn
				connectionsLock.Unlock()

				conn = newConn
			}
		}

		if err != nil {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		}

		err = conn.Ping()
		if err != nil {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		} else {
			break
		}
	}

	return err
}
