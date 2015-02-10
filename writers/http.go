package writers

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
)

// NewHttp is Http constructor.
func NewHttp() *Http {
	h := &Http{}
	return h
}

// Http is a writer that simply serialize all readers data to Http.
type Http struct {
	Base
	Url      string
	Method   string
	Headers  string
	Username string
	Password string
}

// headersAsMap parses the headers data as string and returns them as map.
func (h *Http) headersAsMap() map[string]string {
	if h.Headers == "" {
		return nil
	}

	headersInMap := make(map[string]string)

	pairs := strings.Split(h.Headers, ",")

	for _, pairInString := range pairs {
		pair := strings.Split(pairInString, "=")
		if len(pair) >= 2 {
			headersInMap[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	return headersInMap
}

// NewHttpRequest builds and returns http.Request struct.
func (h *Http) NewHttpRequest(dataJson []byte) (*http.Request, error) {
	var err error

	if h.Url == "" {
		return nil, errors.New("Url is undefined.")
	}

	if h.Method == "" {
		return nil, errors.New("Method is undefined.")
	}

	req, err := http.NewRequest(h.Method, h.Url, bytes.NewBuffer(dataJson))
	if err != nil {
		return nil, err
	}

	for key, value := range h.headersAsMap() {
		req.Header.Set(key, value)
	}

	if h.Username != "" {
		req.SetBasicAuth(h.Username, h.Password)
	}

	return req, err
}

// Run executes the writer.
func (h *Http) Run() error {
	h.Data = h.GetReadersData()
	dataJson, err := h.ToJson()

	if err != nil {
		return err
	}

	req, err := h.NewHttpRequest(dataJson)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
