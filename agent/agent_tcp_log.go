package agent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libtime"
)

// SendTCPLog sends log lines to master.
func (a *Agent) SendTCPLog(config resourced_config.LogReceiverConfig) error {
	loglines := a.TCPLogDB.Get("Loglines")
	if len(loglines) <= 0 {
		return nil
	}

	toSend := make(map[string]interface{})
	toSend["Loglines"] = loglines
	toSend["Filename"] = ""

	host, err := a.hostData()
	if err != nil {
		return err
	}
	toSend["Host"] = host

	dataJson, err := json.Marshal(toSend)
	if err != nil {
		return err
	}

	url := a.GeneralConfig.ResourcedMaster.URL + "/api/logs"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		return err
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

		return err
	}

	a.TCPLogDB.Reset("Loglines")

	return nil
}

func (a *Agent) PruneTCPLogs(config resourced_config.LogReceiverConfig) error {
	loglines := a.TCPLogDB.Get("Loglines")
	if int64(len(loglines)) > config.AutoPruneLength {
		a.TCPLogDB.Reset("Loglines")
	}
	return nil
}

// SendTCPLogForever sends log lines to master in an infinite loop.
func (a *Agent) SendTCPLogForever(config resourced_config.LogReceiverConfig) {
	go func(a *Agent, config resourced_config.LogReceiverConfig) {
		for {
			a.SendTCPLog(config)
			a.PruneTCPLogs(config)
			libtime.SleepString(config.WriteToMasterInterval)
		}
	}(a, config)
}
