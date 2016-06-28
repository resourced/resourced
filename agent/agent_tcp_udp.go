package agent

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/narqo/go-dogstatsd-parser"

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

func (a *Agent) NewUDPServer(config resourced_config.ITCPServer, name string) (*net.UDPConn, error) {
	if config.GetAddr() != "" {
		logFields := logrus.Fields{
			"LogReceiver.Addr": config.GetAddr(),
			"LogLevel":         a.GeneralConfig.LogLevel,
		}

		udpAddr, err := net.ResolveUDPAddr("udp", config.GetAddr())
		if err != nil {
			logrus.WithFields(logFields).Error("Failed to run " + name + " server")
			return nil, err
		}

		logrus.WithFields(logFields).Info("Running " + name + " server")
		return net.ListenUDP("udp4", udpAddr)
	}

	return nil, nil
}

func (a *Agent) HandleGraphite(dataInBytes []byte) {
	dataInChunks := strings.Split(string(dataInBytes), " ")

	if len(dataInChunks) >= 2 {
		key := dataInChunks[0]
		value, err := strconv.ParseFloat(dataInChunks[1], 64)

		if err == nil {
			// Loop through blacklist and set key-value if everything is good
			doSetValue := true

			for _, blacklistRegex := range a.GeneralConfig.MetricReceiver.BlacklistCompiled {
				if blacklistRegex.MatchString(key) {
					doSetValue = false
					break
				}
			}

			if doSetValue {
				a.GraphiteDB.Set(key, value)
			}
		}
	}
}

func (a *Agent) HandleStatsD(dataInBytes []byte) {
	metric, err := dogstatsd.Parse(string(dataInBytes))
	if err != nil {
		return
	}

	a.GraphiteDB.Set(metric.Name, metric.Value)
}

func (a *Agent) HandleLog(dataInBytes []byte) {
	a.TCPLogDB.Append("Loglines", string(dataInBytes))
}
