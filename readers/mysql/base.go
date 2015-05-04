// Package mysql gathers MySQL related data from a host.
package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"time"
)

var connections map[string]*sqlx.DB

type Base struct {
	HostAndPort string
	Retries     int
}

func (m *Base) initConnection() error {
	var err error

	if connections == nil {
		connections = make(map[string]*sqlx.DB)
	}

	if m.Retries <= 0 {
		m.Retries = 10
	}

	createConnection := func() error {
		// Do not create connection if one already exist.
		if existingConnection, ok := connections[m.HostAndPort]; ok && existingConnection != nil {
			return nil
		}

		conn, connectionError := sqlx.Open("mysql", fmt.Sprintf("root:@(%v)/?parseTime=true", m.HostAndPort))
		if connectionError == nil && conn != nil {
			connections[m.HostAndPort] = conn
		}

		return connectionError
	}

	for i := 0; i < m.Retries; i++ {
		err = createConnection()

		if err != nil {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		}

		err = connections[m.HostAndPort].Ping()
		if err != nil {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		} else {
			break
		}

	}

	return err
}
