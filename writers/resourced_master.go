package writers

import (
	"encoding/json"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/sethgrid/pester"

	"github.com/resourced/resourced/libmap"
)

func init() {
	Register("ResourcedMasterHost", NewResourcedMasterHost)
}

// NewResourcedMasterHost is ResourcedMasterHost constructor.
func NewResourcedMasterHost() IWriter {
	return &ResourcedMasterHost{}
}

// ResourcedMasterHost is a writer that serialize readers data to ResourcedMasterHost.
type ResourcedMasterHost struct {
	Http
}

func (rmh *ResourcedMasterHost) preProcessData() (map[string]interface{}, error) {
	newData := make(map[string]interface{})

	for readerPath, dataAndHostInterface := range rmh.Data.(map[string]interface{}) {
		dataAndHost := dataAndHostInterface.(map[string]interface{})

		flatten, err := libmap.Flatten(dataAndHost["Data"].(map[string]interface{}), ".")
		if err != nil {
			return nil, err
		}

		dataAndHost["Data"] = flatten
		newData[readerPath] = dataAndHost
	}

	return newData, nil
}

// Run executes the writer.
func (rmh *ResourcedMasterHost) Run() error {
	if rmh.Data == nil {
		return errors.New("Data field is nil.")
	}

	flatten, err := rmh.preProcessData()
	if err != nil {
		return err
	}

	dataJson, err := json.Marshal(flatten)
	if err != nil {
		return err
	}

	req, err := rmh.NewHttpRequest(dataJson)
	if err != nil {
		return err
	}

	client := pester.New()
	client.MaxRetries = int(rmh.MaxRetries)
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = false

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
