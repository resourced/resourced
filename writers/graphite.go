package writers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/nytlabs/gojsonexplode"
)

func init() {
	Register("Graphite", NewGraphite)
}

// NewGraphite is NewGraphite constructor.
func NewGraphite() IWriter {
	return &Graphite{}
}

// Graphite is a writer that serialize readers data to New Relic Insights.
type Graphite struct {
	Base
	HostPort string
	conn     *net.TCPConn
}

// ToJson returns flattened data in JSON
func (g *Graphite) ToJson() ([]byte, error) {
	if g.Data == nil {
		return nil, errors.New("Data field is nil.")
	}

	dataInJson, err := json.Marshal(g.Data)
	if err != nil {
		return nil, err
	}

	return gojsonexplode.Explodejson(dataInJson, ".")
}

func (g *Graphite) NewTCPConn() error {
	if g.conn == nil {
		tcpAddr, err := net.ResolveTCPAddr("tcp", g.HostPort)
		if err != nil {
			return err
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return err
		}

		g.conn = conn
	}

	return nil
}

// Run sends data to remote graphite server
func (g *Graphite) Run() error {
	if g.HostPort == "" {
		return fmt.Errorf("Unable to connect to Graphite server: %s", g.HostPort)
	}

	err := g.NewTCPConn()
	if err != nil {
		return err
	}
	if g.conn == nil {
		return fmt.Errorf("Unable to connect to Graphite server: %s", g.HostPort)
	}

	flattenData := make(map[string]interface{})

	flattenDataJson, err := g.ToJson()
	if err != nil {
		return err
	}

	err = json.Unmarshal(flattenDataJson, &flattenData)
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()

	for key, value := range flattenData {
		key = strings.Replace(key, "/", "", 1)

		graphiteStringFormattedData := fmt.Sprintf("servers.%v.%v %v %v", hostname, key, value, time.Now().Unix())

		_, err = g.conn.Write([]byte(graphiteStringFormattedData))
		if err != nil {
			return err
		}
	}

	return nil
}
