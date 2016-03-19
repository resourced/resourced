package agent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libmap"
	"github.com/resourced/resourced/libtime"
)

type IAutoPrune interface {
	GetAutoPruneLength() int64
}

// LogPayload packages the log data before sending to master.
func (a *Agent) LogPayload(logdb *libmap.TSafeMapStrings, filename string) map[string]interface{} {
	toSend := make(map[string]interface{})

	data := make(map[string]interface{})
	data["Loglines"] = logdb.Get("Loglines")
	data["Filename"] = filename

	toSend["Data"] = data

	host, err := a.hostData()
	if err != nil {
		return toSend
	}
	toSend["Host"] = host

	return toSend
}

// SendLog sends log lines to master.
func (a *Agent) SendLog(logdb *libmap.TSafeMapStrings, filename string) ([]string, error) {
	data := a.LogPayload(logdb, filename)

	loglines := data["Data"].(map[string]interface{})["Loglines"].([]string)
	if len(loglines) <= 0 {
		return nil, nil
	}

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

	client := &http.Client{}
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

	logdb.Reset("Loglines")

	return loglines, err
}

func (a *Agent) PruneLogs(autoPrunner IAutoPrune, logdb *libmap.TSafeMapStrings) error {
	loglines := logdb.Get("Loglines")
	if int64(len(loglines)) > autoPrunner.GetAutoPruneLength() {
		logdb.Reset("Loglines")
	}
	return nil
}

// SendTCPLogForever sends log lines to master in an infinite loop.
func (a *Agent) SendTCPLogForever(config resourced_config.LogReceiverConfig) {
	go func(a *Agent, config resourced_config.LogReceiverConfig) {
		for {
			a.SendLog(a.TCPLogDB, "")
			a.PruneLogs(config, a.TCPLogDB)
			libtime.SleepString(config.WriteToMasterInterval)
		}
	}(a, config)
}
