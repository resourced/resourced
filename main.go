package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"runtime"

	"github.com/resourced/resourced/agent"
	_ "github.com/resourced/resourced/readers/docker"
	_ "github.com/resourced/resourced/readers/mcrouter"
	_ "github.com/resourced/resourced/readers/mysql"
	_ "github.com/resourced/resourced/readers/procfs"
	_ "github.com/resourced/resourced/readers/redis"
)

func init() {
	logLevelString := os.Getenv("RESOURCED_LOG_LEVEL")
	if logLevelString == "" {
		logLevelString = "info"
	}
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err == nil {
		logrus.SetLevel(logLevel)
	}
}

// main runs the web server for resourced.
func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	a, err := agent.New()
	if err != nil {
		logrus.Fatal(err)
	}

	a.RunAllForever()

	addr := os.Getenv("RESOURCED_ADDR")
	if addr == "" {
		addr = ":55555"
	}

	certFile := os.Getenv("RESOURCED_CERT_FILE")
	keyFile := os.Getenv("RESOURCED_KEY_FILE")

	if certFile != "" && keyFile != "" {
		logrus.WithFields(logrus.Fields{
			"addr": addr,
		}).Info("Running HTTPS server")

		err = http.ListenAndServeTLS(addr, certFile, keyFile, a.HttpRouter())

	} else {
		logrus.WithFields(logrus.Fields{
			"addr": addr,
		}).Info("Running HTTP server")

		err = http.ListenAndServe(addr, a.HttpRouter())
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
