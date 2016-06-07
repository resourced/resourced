package main

import (
	"errors"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	metrics_graphite "github.com/cyberdelia/go-metrics-graphite"

	"github.com/resourced/resourced/agent"
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

	// Graphite Settings
	graphiteListener, err := a.NewTCPServer(a.GeneralConfig.Graphite, "Graphite TCP")
	if err != nil {
		logrus.Fatal(err)
	}
	if graphiteListener != nil {
		defer graphiteListener.Close()

		go func(graphiteListener net.Listener) {
			for {
				conn, err := graphiteListener.Accept()
				if err == nil {
					go a.HandleGraphite(conn)
				}
			}
		}(graphiteListener)
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
				if err == nil {
					go a.HandleLog(conn)
				}
			}
		}(logReceiverListener)
	}

	// Publish metrics to self graphite endpoint.
	addr, err := net.ResolveTCPAddr("tcp", a.GeneralConfig.Graphite.GetAddr())
	if err != nil {
		logrus.Fatal(err)
	}
	statsInterval, err := time.ParseDuration(a.GeneralConfig.Graphite.StatsInterval)
	if err != nil {
		logrus.Fatal(err)
	}
	go metrics_graphite.Graphite(a.NewMetricsRegistry(), statsInterval, "resourced_agent", addr)

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
