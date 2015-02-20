package main

import (
	resourced_agent "github.com/resourced/resourced/agent"
	"os"
	"runtime"
)

// main runs the web server for resourced.
func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	agent, err := resourced_agent.NewAgent()
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
