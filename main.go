package main

import (
	"github.com/Sirupsen/logrus"
	"os"
	"runtime"

	resourced_agent "github.com/resourced/resourced/agent"
	resourced_util "github.com/resourced/resourced/util"
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

	allowedNetworks, cidrErr := resourced_util.ParseCIDRs(os.Getenv("RESOURCED_ALLOWED_NETWORKS"))
	if cidrErr != nil {
		panic(cidrErr)
	}

	agent, err := resourced_agent.NewAgent(allowedNetworks)
	if err != nil {
		panic(err)
	}

	agent.RunAllForever()

	httpAddr := os.Getenv("RESOURCED_ADDR")
	httpsCertFile := os.Getenv("RESOURCED_CERT_FILE")
	httpsKeyFile := os.Getenv("RESOURCED_KEY_FILE")

	if httpsCertFile != "" && httpsKeyFile != "" {
		err = agent.ListenAndServeTLS(httpAddr, httpsCertFile, httpsKeyFile)
	} else {
		err = agent.ListenAndServe(httpAddr)
	}

	if err != nil {
		panic(err)
	}
}
