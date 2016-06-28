package main

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
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
		err := errors.New("RESOURCED_CONFIG_DIR is required")
		logrus.Fatal(err)
	}

	a, err := agent.New()
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

				go a.HandleGraphite(dataInBytes)
				go a.HandleStatsD(dataInBytes)

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

				go a.HandleGraphite(bufferReader[0:n])
				go a.HandleStatsD(bufferReader[0:n])
			}
		}(metricUDPListener)
	}

	// LogReceiver TCP Settings
	logReceiverListener, err := a.NewTCPServer(a.GeneralConfig.LogReceiver, "Log Receiver TCP")
	if err != nil {
		logrus.Fatal(err)
	}
	if logReceiverListener != nil {
		defer logReceiverListener.Close()

		go func(logReceiverListener net.Listener) {
			for {
				conn, err := logReceiverListener.Accept()
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
		}(logReceiverListener)
	}

	// Publish metrics to self graphite endpoint.
	addr, err := net.ResolveTCPAddr("tcp", a.GeneralConfig.MetricReceiver.GetAddr())
	if err != nil {
		logrus.Fatal(err)
	}
	statsInterval, err := time.ParseDuration(a.GeneralConfig.MetricReceiver.StatsInterval)
	if err != nil {
		logrus.Fatal(err)
	}
	go metrics_graphite.Graphite(a.NewMetricsRegistry(), statsInterval, "ResourcedAgent", addr)

	// HTTP Settings
	logFields := logrus.Fields{
		"Addr":                a.GeneralConfig.Addr,
		"LogLevel":            a.GeneralConfig.LogLevel,
		"ResourcedMaster.URL": a.GeneralConfig.ResourcedMaster.URL,
	}

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
