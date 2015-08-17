// Package agent runs readers, writers, and HTTP server.
package agent

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/executors"
	"github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libnet"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	"github.com/resourced/resourced/readers"
	"github.com/resourced/resourced/storage"
	"github.com/resourced/resourced/writers"
	"github.com/resourced/resourced/wstrafficker"
	"github.com/satori/go.uuid"
)

// New is the constructor for Agent struct.
func New() (*Agent, error) {
	agent := &Agent{}

	agent.ID = uuid.NewV4().String()

	err := agent.setConfigs()
	if err != nil {
		return nil, err
	}

	err = agent.setAllowedNetworks()
	if err != nil {
		return nil, err
	}

	err = agent.setTags()
	if err != nil {
		return nil, err
	}

	err = agent.setWSTrafficker()
	if err != nil {
		return nil, err
	}

	err = agent.setStorages()
	if err != nil {
		return nil, err
	}

	return agent, err
}

// Agent struct carries most of the functionality of ResourceD.
// It collects information through readers and serve them up as HTTP+JSON.
type Agent struct {
	ID               string
	Tags             map[string]string
	Configs          *resourced_config.Configs
	GeneralConfig    resourced_config.GeneralConfig
	MetadataStorages *storage.MetadataStorages
	DbPath           string
	Db               *storage.Storage
	AllowedNetworks  []*net.IPNet
	WSTrafficker     *wstrafficker.WSTrafficker
}

func (a *Agent) IsTLS() bool {
	if a.GeneralConfig.HTTPS.CertFile != "" && a.GeneralConfig.HTTPS.KeyFile != "" {
		return true
	}
	return false
}

func (a *Agent) setAllowedNetworks() error {
	allowedNetworks, err := libnet.ParseCIDRs(a.GeneralConfig.AllowedNetworks)
	if err != nil {
		return err
	}

	a.AllowedNetworks = allowedNetworks
	return nil
}

// pathWithPrefix prepends the short version of config.Kind to path.
func (a *Agent) pathWithPrefix(config resourced_config.Config) string {
	if config.Kind == "reader" {
		return a.pathWithKindPrefix("r", config)
	} else if config.Kind == "writer" {
		return a.pathWithKindPrefix("w", config)
	} else if config.Kind == "executor" {
		return a.pathWithKindPrefix("x", config)
	}
	return config.Path
}

// pathWithKindPrefix is common function called by pathWithReaderPrefix or pathWithWriterPrefix
func (a *Agent) pathWithKindPrefix(kind string, input interface{}) string {
	prefix := "/" + kind

	switch v := input.(type) {
	case resourced_config.Config:
		return prefix + v.Path
	case string:
		if strings.HasPrefix(v, prefix+"/") {
			return v
		} else {
			return prefix + v
		}
	}
	return ""
}

// Run executes a reader/writer config.
// Run will save reader data as JSON in local db.
func (a *Agent) Run(config resourced_config.Config) (output []byte, err error) {
	if config.GoStruct != "" && config.Kind == "reader" {
		output, err = a.runGoStructReader(config)
	} else if config.GoStruct != "" && config.Kind == "writer" {
		output, err = a.runGoStructWriter(config)
	} else if config.GoStruct != "" && config.Kind == "executor" {
		output, err = a.runGoStructExecutor(config)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":              err.Error(),
			"config.GoStruct":    config.GoStruct,
			"config.Path":        config.Path,
			"config.Interval":    config.Interval,
			"config.Kind":        config.Kind,
			"config.ReaderPaths": fmt.Sprintf("%s", config.ReaderPaths),
		}).Error("Failed to execute runGoStructReader/runGoStructWriter/runGoStructExecutor")
	}

	err = a.saveRun(config, output, err)

	return output, err
}

// initGoStructReader initialize and return IReader.
func (a *Agent) initGoStructReader(config resourced_config.Config) (readers.IReader, error) {
	return readers.NewGoStructByConfig(config)
}

// initGoStructWriter initialize and return IWriter.
func (a *Agent) initGoStructWriter(config resourced_config.Config) (writers.IWriter, error) {
	writer, err := writers.NewGoStructByConfig(config)
	if err != nil {
		return nil, err
	}

	// Set configs data.
	writer.SetConfigs(a.Configs)

	// Get readers data.
	readersData := make(map[string][]byte)

	for _, readerPath := range config.ReaderPaths {
		readerJsonBytes, err := a.GetRunByPath(a.pathWithKindPrefix("r", readerPath))
		if err == nil {
			readersData[readerPath] = readerJsonBytes
		}
	}

	writer.SetReadersDataInBytes(readersData)

	return writer, err
}

// initResourcedMasterWriter initialize ResourceD Master specific IWriter.
func (a *Agent) initResourcedMasterWriter(config resourced_config.Config) (writers.IWriter, error) {
	var apiPath string

	if config.GoStruct == "ResourcedMasterHost" {
		apiPath = "/api/hosts"
	} else if config.GoStruct == "ResourcedMasterExecutors" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		apiPath = "/api/executors/" + hostname
	}

	urlFromConfigInterface, ok := config.GoStructFields["Url"]
	if !ok || urlFromConfigInterface == nil { // Check if Url is not defined in config
		config.GoStructFields["Url"] = a.GeneralConfig.ResourcedMaster.URL + apiPath

	} else { // Check if Url does not contain apiPath
		urlFromConfig := urlFromConfigInterface.(string)
		if !strings.HasSuffix(urlFromConfig, apiPath) {
			config.GoStructFields["Url"] = a.GeneralConfig.ResourcedMaster.URL + apiPath
		}
	}

	// Check if username is not defined
	// If so, set GeneralConfig.ResourcedMaster.AccessToken as default
	usernameFromConfigInterface, ok := config.GoStructFields["Username"]
	if !ok || usernameFromConfigInterface == nil {
		config.GoStructFields["Username"] = a.GeneralConfig.ResourcedMaster.AccessToken

	}

	return a.initGoStructWriter(config)
}

// initGoStructExecutor initialize and return IExecutor.
func (a *Agent) initGoStructExecutor(config resourced_config.Config) (executors.IExecutor, error) {
	executor, err := executors.NewGoStructByConfig(config)
	if err != nil {
		return nil, err
	}

	executor.SetReadersDataInBytes(a.Db.Data)
	executor.SetTags(a.Tags)
	executor.SetMetadataStorages(a.MetadataStorages)

	return executor, nil
}

// runGoStruct executes Run() fom IReader/IWriter/IExecutor and returns the output.
// Note that IWriter and IExecutor also implement IReader.
func (a *Agent) runGoStruct(readerOrWriterOrExecutor readers.IReader) ([]byte, error) {
	err := readerOrWriterOrExecutor.Run()
	if err != nil {
		errData := make(map[string]string)
		errData["Error"] = err.Error()
		return json.Marshal(errData)
	}

	return readerOrWriterOrExecutor.ToJson()
}

// runGoStructReader executes IReader and returns the output.
func (a *Agent) runGoStructReader(config resourced_config.Config) ([]byte, error) {
	// Initialize IReader
	reader, err := a.initGoStructReader(config)
	if err != nil {
		return nil, err
	}

	return a.runGoStruct(reader)
}

// runGoStructWriter executes IWriter and returns error if exists.
func (a *Agent) runGoStructWriter(config resourced_config.Config) ([]byte, error) {
	var writer writers.IWriter
	var err error

	// Initialize IWriter
	if strings.HasPrefix(config.GoStruct, "ResourcedMaster") {
		writer, err = a.initResourcedMasterWriter(config)
		if err != nil {
			return nil, err
		}

	} else {
		writer, err = a.initGoStructWriter(config)
		if err != nil {
			return nil, err
		}
	}

	err = writer.GenerateData()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":              err.Error(),
			"config.GoStruct":    config.GoStruct,
			"config.Path":        config.Path,
			"config.Interval":    config.Interval,
			"config.Kind":        config.Kind,
			"config.ReaderPaths": fmt.Sprintf("%s", config.ReaderPaths),
		}).Error("Failed to execute writer.GenerateData()")

		return nil, err
	}

	return a.runGoStruct(writer)
}

// runGoStructExecutor executes IExecutor and returns the output.
func (a *Agent) runGoStructExecutor(config resourced_config.Config) ([]byte, error) {
	var executor executors.IExecutor
	var err error

	// Initialize IExecutor
	executor, err = a.initGoStructExecutor(config)
	if err != nil {
		return nil, err
	}

	return a.runGoStruct(executor)
}

// commonData gathers common information for every reader and writer.
func (a *Agent) commonData(config resourced_config.Config) map[string]interface{} {
	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = config.Path

	if config.Interval == "" {
		config.Interval = "1m"
	}
	record["Interval"] = config.Interval

	if config.GoStruct != "" {
		record["GoStruct"] = config.GoStruct
	}

	return record
}

// hostData builds host related information.
func (a *Agent) hostData() (*host.Host, error) {
	h, err := host.NewHostByHostname()
	if err != nil {
		return nil, err
	}

	h.Tags = a.Tags

	return h, nil
}

// saveRun gathers basic, host, and reader/witer information and save them into local storage.
func (a *Agent) saveRun(config resourced_config.Config, output []byte, err error) error {
	// Do not perform save if config.Path is empty.
	if config.Path == "" {
		return nil
	}

	record := a.commonData(config)

	host, err := a.hostData()
	if err != nil {
		return err
	}
	record["Host"] = host

	if err == nil {
		runData := new(interface{})
		err = json.Unmarshal(output, &runData)
		if err != nil {
			return err
		}
		record["Data"] = runData

	} else {
		errMap := make(map[string]string)
		errMap["Error"] = err.Error()
		record["Data"] = errMap
	}

	recordInJson, err := json.Marshal(record)
	if err != nil {
		return err
	}

	a.Db.Set(a.pathWithPrefix(config), recordInJson)

	return err
}

// GetRun returns the JSON data stored in local storage given Config struct.
func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(a.pathWithPrefix(config))
}

// GetRunByPath returns JSON data stored in local storage given path string.
func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	return a.Db.Get(path), nil
}

// RunForever executes Run() in an infinite loop with a sleep of config.Interval.
func (a *Agent) RunForever(config resourced_config.Config) {
	go func(a *Agent, config resourced_config.Config) {
		for {
			a.Run(config)
			libtime.SleepString(config.Interval)
		}
	}(a, config)
}

// RunAllForever executes all readers & writers in an infinite loop.
func (a *Agent) RunAllForever() {
	for _, config := range a.Configs.Readers {
		a.RunForever(config)
	}
	for _, config := range a.Configs.Writers {
		a.RunForever(config)
	}
	for _, config := range a.Configs.Executors {
		a.RunForever(config)
	}
}

// Check if a given IP:PORT is part of an allowed CIDR
func (a *Agent) IsAllowed(address string) bool {
	// Allow all if we allowed networks is not set
	if len(a.AllowedNetworks) == 0 {
		return true
	}

	ip := libstring.GetIP(address)
	if ip == nil {
		return false
	}

	// Check if IP is in one of our allowed networks
	for _, network := range a.AllowedNetworks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
