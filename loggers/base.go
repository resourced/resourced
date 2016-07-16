// Package loggers provides objects that gathers resource data from a host.
package loggers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
	"github.com/sethgrid/pester"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libmap"
	"github.com/resourced/resourced/logline"
)

var loggerConstructors = make(map[string]func() ILogger)

func init() {
	Register("Base", NewBase)
}

// Register makes any logger constructor available by name.
func Register(name string, constructor func() ILogger) {
	if constructor == nil {
		panic("logger: Register logger constructor is nil")
	}
	if _, dup := loggerConstructors[name]; dup {
		panic("logger: Register called twice for logger constructor " + name)
	}
	loggerConstructors[name] = constructor
}

// NewGoStruct instantiates ILogger
func NewGoStruct(name string) (ILogger, error) {
	constructor, ok := loggerConstructors[name]
	if !ok {
		return nil, errors.New("GoStruct is undefined.")
	}

	return constructor(), nil
}

// NewGoStructByConfig instantiates ILogger given Config struct
func NewGoStructByConfig(config resourced_config.Config) (ILogger, error) {
	lgr, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	lgr.SetSource(config.Source)
	lgr.SetBufferSize(config.BufferSize)
	lgr.SetTargets(config.Targets)

	// Populate ILogger fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(lgr).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return lgr, err
}

// ILogger is generic interface for all loggers.
type ILogger interface {
	SetSource(string)
	GetSource() string

	SetBufferSize(int64)
	GetBufferSize() int64

	SetTargets([]resourced_config.LogTargetConfig)
	GetTargets() []resourced_config.LogTargetConfig

	SetLoglines(string, []string)
	GetLoglines(string) []string
	GetLoglinesLength(string) int
	ResetLoglines(string)
	ProcessOutgoingLoglines([]string, []string) []string

	LogErrorAndResetLoglinesIfNeeded(string, error, string)

	SendLogToMaster(string, string, string, *host.Host, []string, string) error

	SendLogToAgent(string, int, []string, string) error

	WriteToFile(string, []string) error
}

type ILoggerChannel interface {
	ILogger
	RunBlockingChannel(string, <-chan interface{})
}

type ILoggerFile interface {
	ILogger
	RunBlockingFile(string)
}

func NewBase() ILogger {
	b := &Base{}
	b.Data = libmap.NewTSafeMapStrings(nil)
	b.BufferSize = 1000000

	return b
}

type Base struct {
	Source     string
	BufferSize int64
	Targets    []resourced_config.LogTargetConfig

	Data *libmap.TSafeMapStrings
}

// RunBlockingChannel pulls log line from channel continuously.
func (b *Base) RunBlockingChannel(name string, ch <-chan interface{}) {
	for line := range ch {
		b.Data.Append(name, line.(string))
	}
}

// RunBlockingFile tails the file continuously.
func (b *Base) RunBlockingFile(file string) {
	t, err := tail.TailFile(file, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
	})
	if err == nil {
		if !b.Data.Exists(file) {
			b.Data.Set(file, make([]string, 0))
		}

		for line := range t.Lines {
			b.Data.Append(file, line.Text)
		}
	}
}

// SetSource
func (b *Base) SetSource(source string) {
	b.Source = source
}

// GetSource returns the source field.
func (b *Base) GetSource() string {
	return b.Source
}

// SetBufferSize sets BufferSize
func (b *Base) SetBufferSize(bufferSize int64) {
	b.BufferSize = bufferSize
}

// GetBufferSize returns BufferSize
func (b *Base) GetBufferSize() int64 {
	return b.BufferSize
}

// SetTargets sets []LogTargetConfig
func (b *Base) SetTargets(targets []resourced_config.LogTargetConfig) {
	b.Targets = targets
}

// SetLoglines sets loglines.
func (b *Base) SetLoglines(file string, loglines []string) {
	b.Data.Set(file, loglines)
}

// GetLoglines returns loglines.
func (b *Base) GetLoglines(file string) []string {
	return b.Data.Get(file)
}

// GetLoglinesLength returns the count of loglines.
func (b *Base) GetLoglinesLength(file string) int {
	return len(b.Data.Get(file))
}

// ResetLoglines wipes it clean.
func (b *Base) ResetLoglines(file string) {
	b.Data.Reset(file)
}

// GetTargets returns slice of LogTargetConfig.
func (b *Base) GetTargets() []resourced_config.LogTargetConfig {
	return b.Targets
}

func (b *Base) LogErrorAndResetLoglinesIfNeeded(source string, err error, message string) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error(message)

		// Check if we have to prune in-memory log lines.
		if int64(b.GetLoglinesLength(source)) > b.GetBufferSize() {
			b.ResetLoglines(source)
		}
	}
}

// filterLoglines denies logline that matches denyList regex.
func (b *Base) filterLoglines(loglines []string, denyList []string) []string {
	newDenyList := make([]string, 0)
	for _, deny := range denyList {
		if deny != "" {
			newDenyList = append(newDenyList, deny)
		}
	}

	if len(newDenyList) == 0 {
		return loglines
	}

	newLoglines := make([]string, 0)

	for _, logline := range loglines {
		for _, deny := range denyList {
			match, err := regexp.MatchString(deny, logline)
			if err != nil || !match {
				newLoglines = append(newLoglines, logline)
			}
		}
	}

	return newLoglines
}

// ProcessOutgoingLoglines before forwarding to targets.
func (b *Base) ProcessOutgoingLoglines(loglines []string, denyList []string) []string {
	return b.filterLoglines(loglines, denyList)
}

// logPayloadForMaster packages the log data before sending to master.
func (b *Base) logPayloadForMaster(hostData *host.Host, loglines []string, source string) map[string]interface{} {
	toSend := make(map[string]interface{})

	data := make(map[string]interface{})
	data["Loglines"] = loglines
	data["Filename"] = source
	toSend["Data"] = data
	toSend["Host"] = hostData

	return toSend
}

// SendLogToMaster sends log lines to master.
func (b *Base) SendLogToMaster(accessToken, masterURLHost, masterURLPath string, hostData *host.Host, loglines []string, source string) error {
	// Check if loglines contain ResourceD base64 wire protocol.
	// If so, convert to plain text.
	for i, lg := range loglines {
		if strings.HasPrefix(lg, "type:base64") {
			loglines[i] = logline.ParseSingle(lg).EncodePlain()
		}
	}

	if masterURLPath == "" {
		masterURLPath = "/api/logs"
	}

	data := b.logPayloadForMaster(hostData, loglines, source)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := masterURLHost + masterURLPath

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to create request struct for sending data to ResourceD Master")

		return err
	}

	req.SetBasicAuth(accessToken, "")

	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = false

	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":      err.Error(),
			"req.URL":    req.URL.String(),
			"req.Method": req.Method,
		}).Error("Failed to send logs data to ResourceD Master")

		return err
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	return err
}

// SendLogToAgent sends log lines to another agent.
func (b *Base) SendLogToAgent(anotherAgentAddr string, maxRetries int, loglines []string, source string) error {
	if len(loglines) == 0 {
		return nil
	}

	conn, err := net.Dial("tcp", anotherAgentAddr)
	attempts := 0

	for {
		if err == nil {
			break
		}

		if err != nil && attempts > maxRetries {
			return err
		}

		if err != nil {
			attempts = attempts + 1
			time.Sleep(pester.ExponentialJitterBackoff(attempts))
			conn, err = net.Dial("tcp", anotherAgentAddr)
			continue
		}
	}

	if conn != nil {
		defer conn.Close()

		w := bufio.NewWriter(conn)

		for i, lg := range loglines {
			// Check if each logline is NOT encoded in ResourceD log wire protocol
			if !strings.HasPrefix(lg, "type:base64") && !strings.HasPrefix(lg, "type:plain") {
				loglines[i] = logline.LiveLogline{Created: time.Now().UTC().Unix(), Content: lg}.EncodeBase64()
				lg = loglines[i]
			}

			fmt.Fprint(w, lg)
			w.Flush()
		}
	}

	return err
}

// WriteToFile writes log lines to local file.
func (b *Base) WriteToFile(targetFile string, loglines []string) error {
	// Check if loglines contain ResourceD base64 wire protocol.
	// If so, convert to plain text.
	for i, lg := range loglines {
		if strings.HasPrefix(lg, "type:base64") {
			loglines[i] = logline.ParseSingle(lg).PlainContent()
		}
	}

	fileHandle, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	for _, logline := range loglines {
		fileHandle.WriteString(logline + "\n")
	}

	return nil
}
