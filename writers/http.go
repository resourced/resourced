package writers

import (
	"bytes"
	"github.com/ddliu/go-httpclient"
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
	Headers map[string]string
}

func (h *Http) Run() error {
	h.Data = h.GetReadersData()
	inJson, err := h.ToJson()

	if err != nil {
		return err
	}

	client := httpclient.NewHttpClient().Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "ResourceD/1.0",
		"Accept-Language":        "en-us",
	})

	_, err = client.Do(h.Method, h.Url, h.Headers, bytes.NewReader(inJson))
	return err
}
