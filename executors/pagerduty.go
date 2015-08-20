package executors

import (
	"encoding/json"

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

		pd.Data["IncidentKey"] = response.IncidentKey
		pd.Data["Status"] = response.Status
		pd.Data["StatusCode"] = statusCode
		pd.Data["Message"] = response.Message
		pd.Data["Errors"] = response.Errors

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
