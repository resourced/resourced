package main

import (
	"github.com/Sirupsen/logrus"
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

	ag, err := agent.New()
	if err != nil {
		panic(err)
	}

	ag.RunAllForever()

	httpAddr := os.Getenv("RESOURCED_ADDR")
	if httpAddr == "" {
		httpAddr = ":55555"
	}

	httpsCertFile := os.Getenv("RESOURCED_CERT_FILE")
	httpsKeyFile := os.Getenv("RESOURCED_KEY_FILE")

	if httpsCertFile != "" && httpsKeyFile != "" {
		err = ag.ListenAndServeTLS(httpAddr, httpsCertFile, httpsKeyFile)
	} else {
		err = ag.ListenAndServe(httpAddr)
	}

	if err != nil {
		panic(err)
	}
}
