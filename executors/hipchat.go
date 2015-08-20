package executors

import (
	"encoding/json"
	"fmt"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

func init() {
	Register("HipChat", NewHipChat)
}

func NewHipChat() IExecutor {
	hc := &HipChat{}
	hc.Data = make(map[string]interface{})

	return hc
}

type HipChat struct {
	Base
	Data      map[string]interface{}
	AuthToken string
	RoomName  string
	Message   string
}

// Run shells out external program and store the output on c.Data.
func (hc *HipChat) Run() error {
	hc.Data["Conditions"] = hc.Conditions

	if hc.IsConditionMet() && hc.LowThresholdExceeded() && !hc.HighThresholdExceeded() {
		c := hipchat.NewClient(hc.AuthToken)

		rooms, _, err := c.Room.List()
		if err != nil {
			return err
		}

		hc.Data["Message"] = fmt.Sprintf("Conditions: %v, Message: %v", hc.Conditions, hc.Message)

		notificationReq := &hipchat.NotificationRequest{Message: hc.Data["Message"].(string)}

		for _, room := range rooms.Items {
			if room.Name == hc.RoomName {
				_, err := c.Room.Notification(room.Name, notificationReq)
				if err != nil {
					return err
				}
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (hc *HipChat) ToJson() ([]byte, error) {
	return json.Marshal(hc.Data)
}
