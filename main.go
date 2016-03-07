package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced/agent"
	resourced_config "github.com/resourced/resourced/config"
	_ "github.com/resourced/resourced/readers/docker"
	_ "github.com/resourced/resourced/readers/haproxy"
	_ "github.com/resourced/resourced/readers/mcrouter"
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

	// Create default configDir if necessary
	if _, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			err := resourced_config.NewDefaultConfigs(configDir)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Generated default configurations inside %v", configDir)
		}
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
	if a.GeneralConfig.Graphite.Addr != "" {
		var graphiteListener net.Listener

		logFields := logrus.Fields{
			"Graphite.Addr": a.GeneralConfig.Graphite.Addr,
			"LogLevel":      a.GeneralConfig.LogLevel,
		}

		if a.GeneralConfig.Graphite.CertFile != "" && a.GeneralConfig.Graphite.KeyFile != "" {
			logFields["Graphite.CertFile"] = a.GeneralConfig.Graphite.CertFile
			logFields["Graphite.KeyFile"] = a.GeneralConfig.Graphite.KeyFile

			cert, err := tls.LoadX509KeyPair(a.GeneralConfig.Graphite.CertFile, a.GeneralConfig.Graphite.KeyFile)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.WithFields(logFields).Info("Running Graphite TCP+SSL server")

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

			graphiteListener, err = tls.Listen("tcp", a.GeneralConfig.Graphite.Addr, tlsConfig)
			if err != nil {
				logrus.Fatal(err)
			}
			defer graphiteListener.Close()

		} else {
			logrus.WithFields(logFields).Info("Running Graphite TCP server")

			graphiteListener, err = net.Listen("tcp", a.GeneralConfig.Graphite.Addr)
			if err != nil {
				logrus.Fatal(err)
			}
			defer graphiteListener.Close()
		}

		if graphiteListener != nil {
			go func(graphiteListener net.Listener) {
				for {
					conn, err := graphiteListener.Accept()
					if err == nil {
						go a.HandleGraphite(conn)
					}
				}
			}(graphiteListener)
		}
	}

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
