package agent

import (
	"fmt"
	net_url "net/url"
)

func (a *Agent) initializeWSConnection() error {
	var wsUrl string

	for _, writer := range a.Configs.Writers {
		if writer.GoStruct == "ResourcedMaster" {
			urlInterface := writer.GoStructFields["Url"]
			if urlInterface == nil {
				return nil
			}

			url, err := net_url.Parse(urlInterface.(string))
			if err != nil {
				return err
			}

			wsUrl = fmt.Sprintf("%v://%v/api/ws", url.Scheme, url.Host)
			break
		}
	}

	if wsUrl != "" {
		return nil
	}

	return nil
}
