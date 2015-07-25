package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"runtime"

	"github.com/resourced/resourced/agent"
	_ "github.com/resourced/resourced/readers/docker"
	_ "github.com/resourced/resourced/readers/mcrouter"
	_ "github.com/resourced/resourced/readers/mysql"
	_ "github.com/resourced/resourced/readers/procfs"
	_ "github.com/resourced/resourced/readers/redis"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// main runs the web server for resourced.
func main() {
	a, err := agent.New()
	if err != nil {
		logrus.Fatal(err)
	}

	logLevel, err := logrus.ParseLevel(a.GeneralConfig.LogLevel)
	if err == nil {
		logrus.SetLevel(logLevel)
	}

	a.RunAllForever()

	if a.GeneralConfig.HTTPS.CertFile != "" && a.GeneralConfig.HTTPS.KeyFile != "" {
		logrus.WithFields(logrus.Fields{
			"GeneralConfig.Addr":           a.GeneralConfig.Addr,
			"GeneralConfig.HTTPS.CertFile": a.GeneralConfig.HTTPS.CertFile,
			"GeneralConfig.HTTPS.KeyFile":  a.GeneralConfig.HTTPS.KeyFile,
		}).Info("Running HTTPS server")

		err = http.ListenAndServeTLS(a.GeneralConfig.Addr, a.GeneralConfig.HTTPS.CertFile, a.GeneralConfig.HTTPS.KeyFile, a.HttpRouter())

	} else {
		logrus.WithFields(logrus.Fields{
			"GeneralConfig.Addr": a.GeneralConfig.Addr,
		}).Info("Running HTTP server")

		err = http.ListenAndServe(a.GeneralConfig.Addr, a.HttpRouter())
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
