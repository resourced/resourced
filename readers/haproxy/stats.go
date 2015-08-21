// Package haproxy gathers haproxy related data from a host.
package haproxy

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
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

// CSVtoJSON parses CSV to JSON
// This function is very specific to parsing HAProxy output.
func CSVtoJSON(csvInput string) ([]byte, error) {
	csvReader := csv.NewReader(strings.NewReader(csvInput))
	lineCount := 0

	var headers []string
	var result bytes.Buffer
	var item bytes.Buffer

	result.WriteString("[")

	for {
		// read just one record, but we could ReadAll() as well
		record, err := csvReader.Read()

		if err == io.EOF {
			result.Truncate(int(len(result.String()) - 1))
			result.WriteString("]")
			break

		} else if err != nil {
			fmt.Println("Error:", err)
			return []byte(""), err
		}

		if lineCount == 0 {
			headers = record[:]

			// Unfortunate hack, HAProxy CSV has 1 key that starts with #
			if headers[0] == "# pxname" {
				headers[0] = "pxname"
			}

			lineCount += 1
		} else {
			item.WriteString("{")

			rowsSlice := make([]string, 0)

			for i := 0; i < len(headers); i++ {
				if headers[i] == "" || record[i] == "" {
					continue
				} else {
					_, err := strconv.ParseFloat(record[i], 64)
					if err != nil {
						rowsSlice = append(rowsSlice, fmt.Sprintf(`"%v": "%v"`, headers[i], record[i]))
					} else {
						rowsSlice = append(rowsSlice, fmt.Sprintf(`"%v": %v`, headers[i], record[i]))
					}
				}
			}

			item.WriteString(strings.Join(rowsSlice, ","))
			item.WriteString("}")

			result.WriteString(item.String() + ",")
			item.Reset()
			lineCount += 1
		}
	}

	return result.Bytes(), nil
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

	inJson, err := CSVtoJSON(string(csv))
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
