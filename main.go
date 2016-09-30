package main

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	metrics_graphite "github.com/cyberdelia/go-metrics-graphite"

	"github.com/resourced/resourced/agent"
	"github.com/resourced/resourced/libtime"
	_ "github.com/resourced/resourced/readers/docker"
	_ "github.com/resourced/resourced/readers/haproxy"
	_ "github.com/resourced/resourced/readers/mcrouter"
	_ "github.com/resourced/resourced/readers/memcache"
	_ "github.com/resourced/resourced/readers/mysql"
	_ "github.com/resourced/resourced/readers/procfs"
	_ "github.com/resourced/resourced/readers/redis"
	_ "github.com/resourced/resourced/readers/varnish"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// main runs the web server for resourced.
func main() {
	configDir := os.Getenv("RESOURCED_CONFIG_DIR")
	if configDir == "" {
		err := errors.New("RESOURCED_CONFIG_DIR is required. Setting it to .")
		logrus.Error(err)

		configDir = "."
	}

	a, err := agent.New(configDir)
	if err != nil {
		logrus.Fatal(err)
	}

	logLevel, err := logrus.ParseLevel(a.GeneralConfig.LogLevel)
	if err == nil {
		logrus.SetLevel(logLevel)
	}

	a.RunAllForever()

	// Metrics TCP Settings
	metricTCPListener, err := a.NewTCPServer(a.GeneralConfig.MetricReceiver, "Metrics Receiver TCP")
	if err != nil {
		logrus.Fatal(err)
	}
	if metricTCPListener != nil {
		defer metricTCPListener.Close()

		go func(metricTCPListener net.Listener) {
			for {
				conn, err := metricTCPListener.Accept()
				if err != nil {
					libtime.SleepString("1s")
					continue
				}

				dataInBytes, err := ioutil.ReadAll(conn)
				if err != nil {
					libtime.SleepString("1s")
					continue
				}

				if strings.Contains(string(dataInBytes), " ") {
					go a.HandleGraphite(dataInBytes)
				} else {
					go a.HandleStatsD(dataInBytes)
				}

				conn.Write([]byte(""))
				conn.Close()
			}
		}(metricTCPListener)
	}

	// Metrics UDP Settings
	metricUDPListener, err := a.NewUDPServer(a.GeneralConfig.MetricReceiver, "Metrics Receiver UDP")
	if err != nil {
		logrus.Fatal(err)
	}
	if metricUDPListener != nil {
		defer metricUDPListener.Close()

		go func(metricUDPListener *net.UDPConn) {
			bufferReader := make([]byte, 1024)

			for {
				n, _, err := metricUDPListener.ReadFromUDP(bufferReader)
				if err != nil {
					libtime.SleepString("1s")
					continue
				}

				if strings.Contains(string(bufferReader[0:n]), " ") {
					go a.HandleGraphite(bufferReader[0:n])
				} else {
					go a.HandleStatsD(bufferReader[0:n])
				}
			}
		}(metricUDPListener)
	}

	// LogReceiver TCP Settings
	logReceiverTCPListener, err := a.NewTCPServer(a.GeneralConfig.LogReceiver, "Log Receiver TCP")
	if err != nil {
		logrus.Fatal(err)
	}
	if logReceiverTCPListener != nil {
		defer logReceiverTCPListener.Close()

		go func(logReceiverTCPListener net.Listener) {
			for {
				conn, err := logReceiverTCPListener.Accept()
				if err != nil {
					libtime.SleepString("1s")
					continue
				}

				dataInBytes, err := ioutil.ReadAll(conn)
				if err != nil {
					libtime.SleepString("1s")
					continue
				}

				go a.HandleLog(dataInBytes)

				conn.Write([]byte(""))
				conn.Close()
			}
		}(logReceiverTCPListener)
	}

	// Create TCP connection to publish metrics to localhost
	addr, err := net.ResolveTCPAddr("tcp", a.GeneralConfig.MetricReceiver.GetAddr())
	if err != nil {
		logrus.Fatal(err)
	}

	// Publish self metrics to localhost
	if a.GeneralConfig.MetricReceiver.StatsInterval != "" {
		statsInterval, err := time.ParseDuration(a.GeneralConfig.MetricReceiver.StatsInterval)
		if err != nil {
			logrus.Fatal(err)
		}

		go metrics_graphite.Graphite(a.NewMetricsRegistryForSelf(), statsInterval, "ResourcedAgent", addr)
		go a.FlushStatsDMetricsToResultDB(statsInterval)
	}

	// HTTP Settings
	logFields := a.DefaultLogrusFieldsForHTTP()

	if a.GeneralConfig.HTTPS.CertFile != "" && a.GeneralConfig.HTTPS.KeyFile != "" {
		logFields["HTTPS.CertFile"] = a.GeneralConfig.HTTPS.CertFile
		logFields["HTTPS.KeyFile"] = a.GeneralConfig.HTTPS.KeyFile

		logrus.WithFields(logFields).Info("Running HTTPS server")

		err = http.ListenAndServeTLS(a.GeneralConfig.Addr, a.GeneralConfig.HTTPS.CertFile, a.GeneralConfig.HTTPS.KeyFile, a.HttpRouter())

	} else {
		logrus.WithFields(logFields).Info("Running HTTP server")

		err = http.ListenAndServe(a.GeneralConfig.Addr, a.HttpRouter())
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
