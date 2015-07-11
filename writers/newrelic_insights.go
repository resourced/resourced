package writers

import (
	"encoding/json"
	"errors"
	"github.com/jmoiron/jsonq"
	"net/http"
	"os"
)

func init() {
	Register("NewrelicInsights", NewNewrelicInsights)
}

// NewNewrelicInsights is NewrelicInsights constructor.
func NewNewrelicInsights() IWriter {
	return &NewrelicInsights{}
}

// NewrelicInsights is a writer that serialize readers data to New Relic Insights.
type NewrelicInsights struct {
	Http
	EventType string
}

func (nr *NewrelicInsights) reformatDataBeforeToJson(data interface{}) interface{} {
	hostname, err := os.Hostname()
	if err != nil {
		return data
	}

	dataAsMapStringInterface, isMapStringInterface := data.(map[string]interface{})
	dataAsSliceInterface, isSliceInterface := data.([]interface{})

	if isSliceInterface {
		newData := make([]interface{}, 0)

		for _, sliceData := range dataAsSliceInterface {
			sliceDataAsMapStringInterface, isSliceDataAMapStringInterface := sliceData.(map[string]interface{})
			if isSliceDataAMapStringInterface {
				sliceDataAsMapStringInterface["Hostname"] = hostname
				sliceDataAsMapStringInterface["eventType"] = nr.EventType

				newData = append(newData, sliceDataAsMapStringInterface)
			}
		}
		return newData
	}

	if isMapStringInterface {
		newData := make(map[string]interface{})

		hasOnlyOneReadersData := len(dataAsMapStringInterface) == 1

		jq := jsonq.NewQuery(dataAsMapStringInterface)

		for readerPath, _ := range dataAsMapStringInterface {
			dataPayload, err := jq.Object(readerPath, "Data")
			if err == nil {
				if hasOnlyOneReadersData {
					newData = dataPayload
				} else {
					newData[readerPath] = dataPayload
				}
			}
		}

		newData["eventType"] = nr.EventType
		newData["Hostname"] = hostname

		return newData
	}

	return data
}

// ToJson serialize Data field to JSON.
func (nr *NewrelicInsights) ToJson() ([]byte, error) {
	return json.Marshal(nr.reformatDataBeforeToJson(nr.Data))
}

// Run executes the writer.
func (nr *NewrelicInsights) Run() error {
	if nr.EventType == "" {
		return errors.New("EventType field is missing.")
	}

	if nr.Data == nil {
		return errors.New("Data field is nil.")
	}

	dataJson, err := nr.ToJson()
	if err != nil {
		return err
	}

	req, err := nr.NewHttpRequest(dataJson)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		resp.Body.Close()
	}

	return nil
}
