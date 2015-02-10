package writers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// NewResourcedMaster is ResourcedMaster constructor.
func NewResourcedMaster() *ResourcedMaster {
	rm := &ResourcedMaster{}
	return rm
}

// ResourcedMaster is a writer that simply serialize all readers data to ResourcedMaster.
type ResourcedMaster struct {
	Http
}

// Run pushes every reader data to ResourcedMaster keyed by reader's key itself.
func (rm *ResourcedMaster) Run() error {
	readersData := rm.GetReadersData()

	responseChan := make(chan *http.Response)
	defer close(responseChan)

	errorChan := make(chan error)
	defer close(errorChan)

	errorStringSlice := make([]string, 0)

	var wg sync.WaitGroup

	for path, data := range readersData {
		wg.Add(1)

		dataJson, err := json.Marshal(data)
		if err == nil {
			req, err := rm.NewHttpRequest(dataJson)
			if err == nil {
				go func(path string, req *http.Request, responseChan chan *http.Response, errorChan chan error) {
					defer wg.Done()

					client := &http.Client{}
					resp, err := client.Do(req)

					if resp != nil {
						responseChan <- resp
					}
					if err != nil {
						errorChan <- err
					}
				}(path, req, responseChan, errorChan)
			}
		}
	}

	go func(errorStringSlice []string) {
		for resp := range responseChan {
			if resp.StatusCode != 200 {
				errorStringSlice = append(errorStringSlice, fmt.Sprintf("Failed to POST to: %v", resp.Request.URL))
			}
		}
	}(errorStringSlice)

	go func(errorStringSlice []string) {
		for err := range errorChan {
			errorStringSlice = append(errorStringSlice, err.Error())
		}
	}(errorStringSlice)

	wg.Wait()

	if len(errorStringSlice) > 0 {
		return errors.New("Errors: " + strings.Join(errorStringSlice, ", "))
	}

	return nil
}
