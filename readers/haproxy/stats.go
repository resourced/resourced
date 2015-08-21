// Package haproxy gathers haproxy related data from a host.
package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("HAProxyStats", NewHAProxyStats)
}

func NewHAProxyStats() readers.IReader {
	hap := &HAProxyStats{}
	hap.Data = make([]map[string]interface{}, 0)

	return hap
}

type HAProxyStats struct {
	Data []map[string]interface{}
	Url  string
}

// Run executes the writer.
func (hap *HAProxyStats) Run() error {
	resp, err := http.Get(hap.Url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"Url":   hap.Url,
		}).Error("Failed to send HTTP request")

		return err
	}

	defer resp.Body.Close()

	csv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"Url":   hap.Url,
		}).Error("Failed to get CSV data")

		return err
	}

	inJson, err := libstring.CSVtoJSON(string(csv))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"Url":   hap.Url,
		}).Error("Failed to parse CSV to JSON")

		return err
	}

	err = json.Unmarshal(inJson, &hap.Data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"Url":   hap.Url,
		}).Error("Failed to decode JSON")

		return err
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (hap *HAProxyStats) ToJson() ([]byte, error) {
	return json.Marshal(hap.Data)
}
