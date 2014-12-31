package main

import (
	resourced_agent "github.com/resourced/resourced/agent"
	"os"
)

func main() {
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
