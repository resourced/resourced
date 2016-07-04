package writers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/resourced/resourced/libmap"
)

func init() {
	Register("Graphite", NewGraphite)
}

// NewGraphite is Graphite constructor.
func NewGraphite() IWriter {
	return &Graphite{}
}

// Graphite is a writer that simply serialize all readers data to Graphite.
type Graphite struct {
	Base
	Addr     string
	Protocol string
	Prefix   string
}

// Run executes the writer.
func (g *Graphite) Run() error {
	now := time.Now().UTC().Unix()

	if g.Data == nil {
		return errors.New("Data field is nil.")
	}

	flatten, err := libmap.Flatten(g.Data, ".")
	if err != nil {
		return err
	}

	if strings.ToLower(g.Protocol) == "tcp" {
		conn, err := net.Dial("tcp", g.Addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		w := bufio.NewWriter(conn)

		for key, value := range flatten {
			if g.Prefix != "" {
				key = g.Prefix + "." + key
			}

			logrus.WithFields(logrus.Fields{
				"Key":       key,
				"Value":     value,
				"Timestamp": now,
			}).Error("Sending metric to Graphite TCP endpoint")

			fmt.Fprintf(w, "%s %f %d\n", key, value, now)
			w.Flush()
		}

	} else if strings.ToLower(g.Protocol) == "udp" {
		conn, err := net.Dial("udp", g.Addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		buffer := make([]byte, 1024)

		for key, value := range flatten {
			if g.Prefix != "" {
				key = g.Prefix + "." + key
			}

			logrus.WithFields(logrus.Fields{
				"Key":       key,
				"Value":     value,
				"Timestamp": now,
			}).Error("Sending metric to Graphite UDP endpoint")

			fmt.Fprintf(conn, "%s %f %d\n", key, value, now)

			_, err = bufio.NewReader(conn).Read(buffer)
			if err != nil {
				return err
			}
			buffer = nil
		}
	}

	return nil
}
