package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var connections map[string]*sqlx.DB

type Base struct {
	HostAndPort string
}

func (m *Base) initConnection() error {
	var err error

	if connections == nil {
		connections = make(map[string]*sqlx.DB)
	}

	if _, ok := connections[m.HostAndPort]; !ok {
		connections[m.HostAndPort], err = sqlx.Open("mysql", fmt.Sprintf("root:@(%v)/?parseTime=true", m.HostAndPort))
	}

	err = connections[m.HostAndPort].Ping()

	return err
}
