package agent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/sethgrid/pester"
)

type IAutoPrune interface {
	GetBufferSize() int64
}

// LogPayloadForMaster packages the log data before sending to master.
func (a *Agent) LogPayloadForMaster(loglines []string, filename string) map[string]interface{} {
	toSend := make(map[string]interface{})

	data := make(map[string]interface{})
	data["Loglines"] = loglines
	data["Filename"] = filename

	toSend["Data"] = data

	host, err := a.hostData()
	if err != nil {
		return toSend
	}
	toSend["Host"] = host

	return toSend
}

// SendLogToMaster sends log lines to master.
func (a *Agent) SendLogToMaster(loglines []string, filename string) ([]string, error) {
	data := a.LogPayloadForMaster(loglines, filename)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := a.GeneralConfig.ResourcedMaster.URL + "/api/logs"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.GeneralConfig.ResourcedMaster.AccessToken, "")

	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = false

	resp, err := client.Do(req)

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err.Error(),
			"req.URL":    req.URL.String(),
			"req.Method": req.Method,
		}).Error("Failed to send logs data to ResourceD Master")

		return nil, err
	}

	return loglines, err
}
