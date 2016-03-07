package agent

import (
	"io/ioutil"
	"net"
)

func (a *Agent) HandleGraphite(conn net.Conn) {
	bytes, err := ioutil.ReadAll(conn)
	if err == nil {
		println("Incoming Data:")
		println(string(bytes))
	}

	conn.Write([]byte(""))
	conn.Close()
}
