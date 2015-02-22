package redis

import (
	"github.com/garyburd/redigo/redis"
)

var connections map[string]redis.Conn

type Base struct {
	HostAndPort string
}

func (r *Base) initConnection() error {
	var err error

	if connections == nil {
		connections = make(map[string]redis.Conn)
	}

	if r.HostAndPort == "" {
		r.HostAndPort = ":6379"
	}

	if _, ok := connections[r.HostAndPort]; !ok {
		connections[r.HostAndPort], err = redis.Dial("tcp", r.HostAndPort)
	}

	return err
}
