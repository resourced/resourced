package agent

import (
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func (a *Agent) HandleGraphite(conn net.Conn) {
	dataInBytes, err := ioutil.ReadAll(conn)
	if err == nil {
		dataInChunks := strings.Split(string(dataInBytes), " ")

		if len(dataInChunks) >= 2 {
			key := dataInChunks[0]
			value, err := strconv.ParseFloat(dataInChunks[1], 64)

			if err == nil {
				a.GraphiteDB.Set(key, value)
			}
		}
	}

	conn.Write([]byte(""))
	conn.Close()
}
