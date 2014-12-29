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

	err = agent.ListenAndServe(os.Getenv("RESOURCED_ADDR"))

	if err != nil {
		panic(err)
	}
}
