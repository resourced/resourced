package agent

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"

	resourced_config "github.com/resourced/resourced/config"
)

func (a *Agent) NewTCPServer(config resourced_config.TCPConfig, name string) (net.Listener, error) {
	if config.Addr != "" {
		logFields := logrus.Fields{
			"LogReceiver.Addr": config.Addr,
			"LogLevel":         a.GeneralConfig.LogLevel,
		}

		if config.CertFile != "" && config.KeyFile != "" {
			logFields["LogReceiver.CertFile"] = config.CertFile
			logFields["LogReceiver.KeyFile"] = config.KeyFile

			cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
			if err != nil {
				logrus.WithFields(logFields).Fatal(err)
				return nil, err
			}

			logrus.WithFields(logFields).Info("Running " + name + "+SSL server")

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

			return tls.Listen("tcp", config.Addr, tlsConfig)

		} else {
			logrus.WithFields(logFields).Info("Running " + name + " server")

			return net.Listen("tcp", config.Addr)
		}
	}

	return nil, nil
}

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

func (a *Agent) HandleLog(conn net.Conn) {
	dataInBytes, err := ioutil.ReadAll(conn)
	if err == nil {
		a.LogDB.Append("Loglines", string(dataInBytes))
	}

	conn.Write([]byte(""))
	conn.Close()
}
