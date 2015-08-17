package writers

import (
	"net/http"

	"github.com/gogap/logrus"
)

func init() {
	Register("ResourcedMasterHost", NewResourcedMasterHost)
	Register("ResourcedMasterExecutors", NewResourcedMasterExecutors)
}

// NewResourcedMasterHost is ResourcedMasterHost constructor.
func NewResourcedMasterHost() IWriter {
	return &ResourcedMasterHost{}
}

// ResourcedMasterHost is a writer that serialize readers data to ResourcedMasterHost.
type ResourcedMasterHost struct {
	Http
}

// NewResourcedMasterExecutors is ResourcedMasterExecutors constructor.
func NewResourcedMasterExecutors() IWriter {
	return &ResourcedMasterExecutors{}
}

// ResourcedMasterExecutors is a writer that serialize readers data to ResourcedMasterExecutors.
type ResourcedMasterExecutors struct {
	Http
}

// Run executes the writer.
func (rm *ResourcedMasterExecutors) Run() error {
	// Build executors config data.
	executorData := make(map[string]interface{})

	for _, executor := range rm.Configs.Executors {
		executorData[executor.Path] = executor
	}
	rm.SetData(executorData)

	dataJson, err := rm.ToJson()
	if err != nil {
		return err
	}

	req, err := rm.NewHttpRequest(dataJson)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err.Error(),
			"req.URL":    req.URL.String(),
			"req.Method": req.Method,
		}).Error("Failed to send HTTP request")

		return err
	}

	if resp.Body != nil {
		resp.Body.Close()
	}

	return nil
}
