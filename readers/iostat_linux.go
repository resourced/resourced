// +build linux
package readers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	Register("IOStat", NewIOStat)
}

func NewIOStat() IReader {
	ios := &IOStat{}
	ios.Data = make(map[string]map[string]float64)
	return ios
}

type IOStat struct {
	Data map[string]map[string]float64
}

// Run gathers load average information from gosigar.
func (ios *IOStat) Run() error {
	rawData, err := exec.Command("iostat", "-x").Output()
	if err != nil {
		return err
	}

	if rawData != nil {
		reader := bytes.NewReader(rawData)

		var titles []string
		pastTitles := false

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "Device:") {
				titles = strings.Fields(line)
				pastTitles = true
				continue
			}

			if pastTitles {
				rowData := strings.Fields(line)

				var device string

				for i, d := range rowData {
					d = strings.TrimSpace(d)

					title := strings.TrimSpace(titles[i])

					if i == 0 {
						device = d
						ios.Data[device] = make(map[string]float64)

					} else {
						dFloat, err := strconv.ParseFloat(d, 64)
						if err != nil {
							return err
						}

						ios.Data[device][title] = dFloat
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (ios *IOStat) ToJson() ([]byte, error) {
	return json.Marshal(ios.Data)
}
