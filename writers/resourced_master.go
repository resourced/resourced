package writers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
)

func init() {
	Register("ResourcedMasterHost", NewResourcedMasterHost)
	Register("ResourcedMasterStacks", NewResourcedMasterStacks)
}

// NewResourcedMasterHost is ResourcedMasterHost constructor.
func NewResourcedMasterHost() IWriter {
	return &ResourcedMasterHost{}
}

// ResourcedMasterHost is a writer that serialize readers data to ResourcedMasterHost.
type ResourcedMasterHost struct {
	Http
}

// NewResourcedMasterStacks is ResourcedMasterStacks constructor.
func NewResourcedMasterStacks() IWriter {
	return &ResourcedMasterStacks{}
}

// ResourcedMasterStacks is a writer that serialize readers data to ResourcedMasterStacks.
type ResourcedMasterStacks struct {
	Http
	Root string
}

// stacksData gathers complete list of ResourceD Stacks metadata.
func (rm *ResourcedMasterStacks) stacksData() map[string]interface{} {
	data := make(map[string]interface{})
	return data
}

// Run executes the writer.
func (rm *ResourcedMasterStacks) Run() error {
	if rm.Root == "" {
		return errors.New("ResourceD Stacks root should not be empty")
	}

	callback := func() {
		if rm.Data == nil {
			return
		}

		data := rm.stacksData()

		dataJson, err := json.Marshal(data)
		if err != nil {
			return
		}

		req, err := rm.NewHttpRequest(dataJson)
		if err != nil {
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":      err.Error(),
				"req.URL":    req.URL.String(),
				"req.Method": req.Method,
			}).Error("Failed to send HTTP request")

			return
		}

		if resp.Body != nil {
			resp.Body.Close()
		}
	}
	return rm.WatchDir(rm.Root, callback)
}
