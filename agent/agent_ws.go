package agent

import (
	"fmt"
	net_url "net/url"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced/wsclient"
	"github.com/resourced/resourced/wstrafficker"
)

// setWSTrafficker construct WSTrafficker instance.
// WSTrafficker carry its own websocket client.
func (a *Agent) setWSTrafficker() error {
	wsPath := "/api/ws"
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

			wsUrl = fmt.Sprintf("%v://%v%v", url.Scheme, url.Host, wsPath)
			break
		}
	}

	if wsUrl != "" {
		wsSettings := make(map[string]interface{})
		wsSettings["Timeout"] = 1 * time.Second

		httpAddr := os.Getenv("RESOURCED_ADDR")
		if httpAddr == "" {
			httpAddr = "localhost:55555"
		}

		wsClient, _, err := wsclient.NewClient("http://"+httpAddr, wsUrl, wsSettings)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":     err.Error(),
				"OriginURL": "http://" + httpAddr,
				"TargetURL": wsUrl,
				"Timeout":   wsSettings["Timeout"],
			}).Error("Failed to establish websocket connection")
			return nil
		}

		a.WSTrafficker = wstrafficker.New(wsClient)
	}

	return nil
}
