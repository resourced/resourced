// Package executors provides objects that gathers resource data from a host.
package executors

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/Sirupsen/logrus"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libmap"
	"github.com/resourced/resourced/queryparser"
)

var executorConstructors = make(map[string]func() IExecutor)

// Register makes any executor constructor available by name.
func Register(name string, constructor func() IExecutor) {
	if constructor == nil {
		panic("executor: Register executor constructor is nil")
	}
	if _, dup := executorConstructors[name]; dup {
		panic("executor: Register called twice for executor constructor " + name)
	}
	executorConstructors[name] = constructor
}

// NewGoStruct instantiates IExecutor
func NewGoStruct(name string) (IExecutor, error) {
	constructor, ok := executorConstructors[name]
	if !ok {
		return nil, errors.New("GoStruct is undefined.")
	}

	return constructor(), nil
}

// NewGoStructByConfig instantiates IExecutor given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IExecutor, error) {
	executor, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	executor.SetPath(config.Path)
	executor.SetInterval(config.Interval)

	// Populate IExecutor fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(executor).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return executor, nil
}

// IExecutor is generic interface for all executors.
type IExecutor interface {
	SetPath(string)
	GetPath() string
	SetInterval(string)
	SetResourcedMasterURL(string)
	SetResourcedMasterAccessToken(string)
	SetHostData(*host.Host)
	SetCounterDB(*libmap.TSafeMapCounter)
	Run() error
	ToJson() ([]byte, error)
	SetQueryParser(map[string][]byte)
	SetReadersDataInBytes(map[string][]byte)
	SetTags(map[string]string)
	IsConditionMet() bool
	LowThresholdExceeded() bool
	HighThresholdExceeded() bool
}

type Base struct {
	// Command: Shell command to execute.
	Command string

	// Path: ResourceD URL path. Example:
	// /uptime -> http://localhost:55555/x/uptime
	Path string

	Interval string

	// LowThreshold: minimum count of valid conditions
	LowThreshold int64

	// HighThreshold: maximum count of valid conditions
	HighThreshold int64

	// Conditions for when executor should run.
	Conditions string

	// Host data
	Host *host.Host

	ResourcedMasterURL         string
	ResourcedMasterAccessToken string

	ReadersDataBytes map[string][]byte

	qp *queryparser.QueryParser

	counterDB *libmap.TSafeMapCounter
}

func (b *Base) SetPath(path string) {
	b.Path = path
}

func (b *Base) GetPath() string {
	return b.Path
}

func (b *Base) SetInterval(interval string) {
	b.Interval = interval
}

func (b *Base) SetQueryParser(readersJsonBytes map[string][]byte) {
	b.qp = queryparser.New(readersJsonBytes, nil)
}

func (b *Base) SetCounterDB(db *libmap.TSafeMapCounter) {
	b.counterDB = db
}

// SetReadersDataInBytes pulls readers data and store them on ReadersData field.
func (b *Base) SetReadersDataInBytes(readersJsonBytes map[string][]byte) {
	b.ReadersDataBytes = readersJsonBytes

	b.SetQueryParser(readersJsonBytes)
}

// SetTags assigns all host tags to qp (QueryParser).
func (b *Base) SetTags(tags map[string]string) {
	b.qp.SetTags(tags)
}

func (b *Base) SetResourcedMasterURL(resourcedMasterURL string) {
	b.ResourcedMasterURL = resourcedMasterURL
}

func (b *Base) SetResourcedMasterAccessToken(resourcedMasterAccessToken string) {
	b.ResourcedMasterAccessToken = resourcedMasterAccessToken
}

func (b *Base) SetHostData(hostData *host.Host) {
	b.Host = hostData
}

func (b *Base) IsConditionMet() bool {
	if b.Conditions == "" {
		b.Conditions = "true"
	}

	result, err := b.qp.Parse(b.Conditions)
	if err != nil {
		return false
	}

	if result == true {
		// Condition is met, increment counter by 1
		b.counterDB.Incr(b.Path, 1)
	} else {
		// Condition is no longer met, reset counter to 0
		b.counterDB.Reset(b.Path)
	}
	return result
}

func (b *Base) LowThresholdExceeded() bool {
	return int64(b.counterDB.Get(b.Path)) > b.LowThreshold
}

func (b *Base) HighThresholdExceeded() bool {
	if b.HighThreshold == 0 {
		return false
	}
	return int64(b.counterDB.Get(b.Path)) > b.HighThreshold
}

// NewHttpRequest builds and returns http.Request struct.
func (b *Base) NewHttpRequest(dataJson []byte) (*http.Request, error) {
	var err error

	if b.ResourcedMasterURL == "" {
		return nil, errors.New("ResourcedMasterURL is undefined.")
	}

	if b.ResourcedMasterAccessToken == "" {
		return nil, errors.New("ResourcedMasterAccessToken is undefined.")
	}

	url := b.ResourcedMasterURL + "/api/executors"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(b.ResourcedMasterAccessToken, "")

	return req, err
}

// Send executor data to master
func (b *Base) SendToMaster(data map[string]interface{}) error {
	toSend := make(map[string]interface{})
	toSend["Data"] = data
	toSend["Host"] = b.Host

	dataJson, err := json.Marshal(toSend)
	if err != nil {
		return err
	}

	req, err := b.NewHttpRequest(dataJson)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err.Error(),
			"req.URL":    req.URL.String(),
			"req.Method": req.Method,
		}).Error("Failed to send executor data to ResourceD Master")

		return err
	}

	if resp.Body != nil {
		resp.Body.Close()
	}

	return nil
}
