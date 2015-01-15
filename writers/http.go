package writers

import (
	"bytes"
	"errors"
	"github.com/ddliu/go-httpclient"
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
	Url     string
	Method  string
	Headers string
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

	client := httpclient.NewHttpClient().Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "ResourceD/1.0",
		"Accept-Language":        "en-us",
	})

	_, err = client.Do(h.Method, h.Url, h.headersAsMap(), bytes.NewReader(inJson))
	return err
}
