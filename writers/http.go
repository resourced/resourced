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

func (h *Http) Run() error {
	h.Data = h.GetReadersData()
	inJson, err := h.ToJson()

	if err != nil {
		return err
	}

	if h.Url == "" {
		return errors.New("Url is undefined.")
	}

	if h.Method == "" {
		return errors.New("Method is undefined.")
	}

	req, err := http.NewRequest(h.Method, h.Url, bytes.NewBuffer(inJson))
	if err != nil {
		return err
	}

	for key, value := range h.headersAsMap() {
		req.Header.Set(key, value)
	}

	if h.Username != "" {
		req.SetBasicAuth(h.Username, h.Password)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
