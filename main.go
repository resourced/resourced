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

	logFields := logrus.Fields{
		"Addr":                       a.GeneralConfig.Addr,
		"LogLevel":                   a.GeneralConfig.LogLevel,
		"AllowedNetworks":            a.GeneralConfig.AllowedNetworks,
		"ResourcedMaster.URL":        a.GeneralConfig.ResourcedMaster.URL,
		"ResourcedStacks.Root":       a.GeneralConfig.ResourcedStacks.Root,
		"ResourcedStacks.PythonPath": a.GeneralConfig.ResourcedStacks.PythonPath,
		"ResourcedStacks.PipPath":    a.GeneralConfig.ResourcedStacks.PipPath,
		"ResourcedStacks.Conditions": a.GeneralConfig.ResourcedStacks.Conditions,
		"ResourcedStacks.DryRun":     a.GeneralConfig.ResourcedStacks.DryRun,
		"ResourcedStacks.Git.HTTPS":  a.GeneralConfig.ResourcedStacks.Git.HTTPS,
		"ResourcedStacks.Git.Branch": a.GeneralConfig.ResourcedStacks.Git.Branch,
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
