package writers

import (
	"encoding/json"
	"github.com/jmoiron/jsonq"
)

// NewNewrelicInsights is NewrelicInsights constructor.
func NewNewrelicInsights() *NewrelicInsights {
	rm := &NewrelicInsights{}
	return rm
}

// NewrelicInsights is a writer that serialize readers data to New Relic Insights.
type NewrelicInsights struct {
	Http
	EventType string
}

func (nr *NewrelicInsights) reformatDataBeforeToJson(data map[string]interface{}) map[string]interface{} {
	newReadersData := make(map[string]interface{})

	hasOnlyOneReadersData := len(data) == 1

	jq := jsonq.NewQuery(data)

	for readerPath, _ := range data {
		dataPayload, err := jq.Object(readerPath, "Data")
		if err == nil {
			if hasOnlyOneReadersData {
				newReadersData = dataPayload
			} else {
				newReadersData[readerPath] = dataPayload
			}
		}

		// If Hostname key is missing...
		if _, ok := newReadersData["Hostname"]; !ok {
			hostname, err := jq.String(readerPath, "Host", "Name")
			if err == nil {
				newReadersData["Hostname"] = hostname
			}
		}
	}

	newReadersData["eventType"] = nr.EventType

	return newReadersData
}

// ToJson serialize Data field to JSON.
func (nr *NewrelicInsights) ToJson() ([]byte, error) {
	return json.Marshal(nr.reformatDataBeforeToJson(nr.Data))
}
