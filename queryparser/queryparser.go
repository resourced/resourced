package queryparser

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/robertkrimen/otto"
)

func New(data map[string][]byte, tags map[string]string) *QueryParser {
	hostname, _ := os.Hostname()

	qp := &QueryParser{}
	qp.hostname = hostname
	qp.tags = tags
	qp.data = data
	return qp
}

type QueryParser struct {
	hostname string
	tags     map[string]string
	data     map[string][]byte
}

func (qp *QueryParser) SetTags(tags map[string]string) {
	qp.tags = tags
}

func (qp *QueryParser) Parse(query string) (bool, error) {
	query, err := qp.replaceDataPathWithValue(query)
	if err != nil {
		return false, err
	}

	value, err := otto.New().Run(query)
	if err != nil {
		return false, err
	}

	return value.ToBoolean()
}

func (qp *QueryParser) dataValue(datapath, jsonSelector string) (interface{}, error) {
	dataJsonBytes := qp.data[datapath]
	var dataJson map[string]interface{}

	err := json.Unmarshal(dataJsonBytes, &dataJson)
	if err != nil {
		return nil, err
	}

	jq := jsonq.NewQuery(dataJson)

	jsonSelectorChunks := strings.Split(jsonSelector, ".")
	jsonSelectorChunks = append([]string{"Data"}, jsonSelectorChunks...) // Always query from "Data" sub-structure.

	return jq.Interface(jsonSelectorChunks...)
}

func (qp *QueryParser) replaceDataPathWithValue(query string) (string, error) {
	queryChunks := strings.Fields(query)

	for i, chunk := range queryChunks {
		if strings.Contains(chunk, "/r/") || strings.Contains(chunk, "/w/") || strings.Contains(chunk, "/x/") {
			openParensCount := strings.Count(chunk, "(")
			chunk = strings.Replace(chunk, "(", "", -1)

			dataPathAndJsonSelectorChunks := strings.Split(chunk, ".")

			dataPath := dataPathAndJsonSelectorChunks[0]
			jsonSelectorChunks := dataPathAndJsonSelectorChunks[1:]
			jsonSelector := strings.Join(jsonSelectorChunks, ".")

			valueInterface, err := qp.dataValue(dataPath, jsonSelector)
			if err != nil {
				return "", err
			}

			valueBytes, err := json.Marshal(valueInterface)
			if err != nil {
				return "", err
			}

			if openParensCount == 0 {
				queryChunks[i] = string(valueBytes)

			} else {
				queryChunks[i] = strings.Repeat("(", openParensCount) + string(valueBytes)
			}

		}
	}

	return strings.Join(queryChunks, " "), nil
}
