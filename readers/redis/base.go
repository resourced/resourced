// Package redis gathers Redis related data from a host.
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

	// Does the connection exist?
	if _, ok := connections[r.HostAndPort]; !ok {
		// Are we able to establish a connection to it?
		if conn, redisErr := redis.Dial("tcp", r.HostAndPort); redisErr == nil {
			connections[r.HostAndPort] = conn
		} else {
			err = redisErr
		}
	}

	return err
}
