package executors

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dselans/pagerduty"
)

func init() {
	Register("PagerDuty", NewPagerDuty)
}

func NewPagerDuty() IExecutor {
	pd := &PagerDuty{}
	pd.Data = make(map[string]interface{})

	return pd
}

type PagerDuty struct {
	Base
	Data        map[string]interface{}
	ServiceKey  string
	Description string
	IncidentKey string
}

// Run shells out external program and store the output on c.Data.
func (pd *PagerDuty) Run() error {
	pd.Data["Conditions"] = pd.Conditions

	if pd.IsConditionMet() && pd.LowThresholdExceeded() && !pd.HighThresholdExceeded() {
		event := pagerduty.NewTriggerEvent(pd.ServiceKey, pd.Description)

		if pd.IncidentKey != "" {
			event.IncidentKey = pd.IncidentKey
		}

		response, statusCode, err := pagerduty.Submit(event)

		if response != nil {
			pd.Data["IncidentKey"] = response.IncidentKey
			pd.Data["Status"] = response.Status
			pd.Data["StatusCode"] = statusCode
			pd.Data["Message"] = response.Message
			pd.Data["Errors"] = response.Errors

			go func() {
				logline := pd.formatBeforeSendingToMaster(pd.Data)
				err := pd.SendToMaster(logline)
				if err != nil {
					logrus.Error(err)
				}
			}()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (pd *PagerDuty) ToJson() ([]byte, error) {
	return json.Marshal(pd.Data)
}

func (pd *PagerDuty) formatBeforeSendingToMaster(data map[string]interface{}) AgentLoglinePayload {
	logline := fmt.Sprintf("Conditions: %v. IncidentKey: %v. ", pd.Conditions, pd.IncidentKey)

	if status, ok := data["Status"]; ok {
		logline = logline + fmt.Sprintf("Status: %v. ", status)
	}
	if statusCode, ok := data["StatusCode"]; ok {
		logline = logline + fmt.Sprintf("StatusCode: %v. ", statusCode)
	}
	if message, ok := data["Message"]; ok {
		logline = logline + fmt.Sprintf("Message: %v. ", message)
	}
	if errors, ok := data["Errors"]; ok && len(errors.([]string)) > 0 {
		logline = logline + fmt.Sprintf("Errors: %v. ", errors)
	}

	return AgentLoglinePayload{Created: time.Now().UTC().Unix(), Content: logline}
}
