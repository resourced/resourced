package writers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"reflect"
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

func (g *Graphite) preProcessKey(key string) string {
	// Strip leading forward slash
	if strings.HasPrefix(key, "/") {
		key = key[1:len(key)]
	}

	// Prepend prefix if defined
	if g.Prefix != "" {
		key = g.Prefix + "." + key
	}

	return key
}

func (g *Graphite) preProcessData() (map[string]interface{}, error) {
	flatten := make(map[string]interface{})

	flattenWithData, err := libmap.Flatten(g.Data, ".")
	if err != nil {
		return flatten, err
	}

	for key, value := range flattenWithData {
		keyWithoutData := strings.Replace(key, ".Data.", ".", 1)
		flatten[keyWithoutData] = value
	}

	return flatten, nil
}

// Run executes the writer.
func (g *Graphite) Run() error {
	now := time.Now().UTC().Unix()

	if g.Data == nil {
		return errors.New("Data field is nil.")
	}

	flatten, err := g.preProcessData()
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
			key = g.preProcessKey(key)

			switch reflect.TypeOf(value).Kind() {
			case reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64:
				logrus.WithFields(logrus.Fields{
					"Key":       key,
					"Value":     value,
					"Timestamp": now,
				}).Info("Sending metric to Graphite TCP endpoint")

				fmt.Fprintf(w, "%s %v %d\n", key, value, now)
				w.Flush()
			}
		}

	} else if strings.ToLower(g.Protocol) == "udp" {
		addr, err := net.ResolveUDPAddr("udp", g.Addr)
		if err != nil {
			return err
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		for key, value := range flatten {
			key = g.preProcessKey(key)

			switch reflect.TypeOf(value).Kind() {
			case reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64:
				logrus.WithFields(logrus.Fields{
					"Key":       key,
					"Value":     value,
					"Timestamp": now,
				}).Info("Sending metric to Graphite UDP endpoint")

				conn.Write([]byte(fmt.Sprintf("%s %v %d\n", key, value, now)))
			}
		}
	}

	return nil
}
