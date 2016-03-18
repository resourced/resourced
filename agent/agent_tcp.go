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

func (a *Agent) NewTCPServer(config resourced_config.ITCPServer, name string) (net.Listener, error) {
	if config.GetAddr() != "" {
		logFields := logrus.Fields{
			"LogReceiver.Addr": config.GetAddr(),
			"LogLevel":         a.GeneralConfig.LogLevel,
		}

		if config.GetCertFile() != "" && config.GetKeyFile() != "" {
			logFields["LogReceiver.CertFile"] = config.GetCertFile()
			logFields["LogReceiver.KeyFile"] = config.GetKeyFile()

			cert, err := tls.LoadX509KeyPair(config.GetCertFile(), config.GetKeyFile())
			if err != nil {
				logrus.WithFields(logFields).Fatal(err)
				return nil, err
			}

			logrus.WithFields(logFields).Info("Running " + name + "+SSL server")

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

			return tls.Listen("tcp", config.GetAddr(), tlsConfig)

		} else {
			logrus.WithFields(logFields).Info("Running " + name + " server")

			return net.Listen("tcp", config.GetAddr())
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
		a.TCPLogDB.Append("Loglines", string(dataInBytes))
	}

	conn.Write([]byte(""))
	conn.Close()
}
