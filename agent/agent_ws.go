package agent

import (
	"encoding/json"
	"fmt"
	net_url "net/url"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced/wsclient"
	"github.com/resourced/resourced/wstrafficker"
)

const (
	websocketPathByAccessTokenPrefix = "/api/ws/access-tokens"
)

// setWSTrafficker construct WSTrafficker instance.
// WSTrafficker carry its own websocket client.
func (a *Agent) setWSTrafficker() error {
	var masterUrl *net_url.URL
	var accessToken string
	var err error

	for _, writer := range a.Configs.Writers {
		if writer.GoStruct == "ResourcedMaster" {
			urlInterface := writer.GoStructFields["Url"]
			if urlInterface == nil {
				return nil
			}

			accessTokenInterface := writer.GoStructFields["Username"]
			if accessTokenInterface == nil {
				return nil
			}
			accessToken = accessTokenInterface.(string)

			masterUrl, err = net_url.Parse(urlInterface.(string))
			if err != nil {
				return err
			}

			break
		}
	}

	if masterUrl != nil && accessToken != "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		originScheme := "ws"
		if a.IsTLS() {
			originScheme = "wss"
		}

		originAddr := os.Getenv("RESOURCED_ADDR")
		if originAddr == "" {
			originAddr = "localhost:55555"
		}

		targetScheme := "ws"
		if masterUrl.Scheme == "https" {
			targetScheme = "wss"
		}

		originURL := fmt.Sprintf("%v://%v", originScheme, originAddr)
		targetURL := fmt.Sprintf("%v://%v%v/%v", targetScheme, masterUrl.Host, websocketPathByAccessTokenPrefix, accessToken)

		wsSettings := make(map[string]interface{})
		wsSettings["Timeout"] = 1 * time.Second

		wsClient, _, err := wsclient.NewClient(originURL, targetURL, wsSettings)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":     err.Error(),
				"OriginURL": originURL,
				"TargetURL": targetURL,
				"Timeout":   wsSettings["Timeout"],
			}).Error("Failed to establish websocket connection")
			return nil
		}

		a.WSTrafficker = wstrafficker.NewWSTrafficker(wsClient)

		payload := make(map[string]string)
		payload["Hostname"] = hostname

		payloadJson, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		err = a.WSTrafficker.Write(1, payloadJson)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"OriginURL": originURL,
			"TargetURL": targetURL,
			"Timeout":   wsSettings["Timeout"],
		}).Info("Established websocket connection")
	}

	return nil
}
