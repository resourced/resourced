// Package agent runs readers, writers, and HTTP server.
package agent

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	resourced_config "github.com/resourced/resourced/config"
	resourced_host "github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	resourced_readers "github.com/resourced/resourced/readers"
	resourced_writers "github.com/resourced/resourced/writers"
	"os"
	"strings"
	"time"
)

// NewAgent is the constructor fot Agent struct.
func NewAgent() (*Agent, error) {
	agent := &Agent{}

	agent.setTags()

	err := agent.setConfigStorage()
	if err != nil {
		return nil, err
	}

	err = agent.setDb()
	if err != nil {
		return nil, err
	}

	return agent, err
}

// Agent struct carries most of the functionality of ResourceD.
// It collects information through readers and serve them up as HTTP+JSON.
type Agent struct {
	ConfigStorage *resourced_config.ConfigStorage
	DbPath        string
	Db            *bolt.DB
	Tags          []string
}

// setTags store RESOURCED_TAGS data to Tags field.
func (a *Agent) setTags() {
	a.Tags = make([]string, 0)

	tags := os.Getenv("RESOURCED_TAGS")
	if tags != "" {
		tagsSlice := strings.Split(tags, ",")
		a.Tags = make([]string, len(tagsSlice))

		for i, tag := range tagsSlice {
			a.Tags[i] = strings.TrimSpace(tag)
		}
	}
}

// setDb configures the local storage.
// The base path to local storage is defined in RESOURCED_DB.
func (a *Agent) setDb() error {
	var err error

	dbPath := os.Getenv("RESOURCED_DB")
	if dbPath == "" {
		dbPath = "~/resourced/db"

		err = os.MkdirAll(libstring.ExpandTildeAndEnv("~/resourced"), 0755)
		if err != nil {
			return err
		}
	}

	a.DbPath = libstring.ExpandTildeAndEnv(dbPath)

	a.Db, err = bolt.Open(a.DbPath, 0644, nil)
	if err != nil {
		return err
	}

	// Create "resources" bucket
	a.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("resources"))
		return nil
	})

	return err
}

// dbBucket returns the boltdb bucket.
func (a *Agent) dbBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("resources"))
}

// pathWithPrefix prepends the short version of config.Kind to path.
func (a *Agent) pathWithPrefix(config resourced_config.Config) string {
	if config.Kind == "reader" {
		return a.pathWithReaderPrefix(config)
	} else if config.Kind == "writer" {
		return a.pathWithWriterPrefix(config)
	}
	return config.Path
}

// pathWithReaderPrefix conveniently assign /r prefix to path.
func (a *Agent) pathWithReaderPrefix(input interface{}) string {
	switch v := input.(type) {
	case resourced_config.Config:
		return "/r" + v.Path
	case string:
		if strings.HasPrefix(v, "/r/") {
			return v
		} else {
			return "/r" + v
		}
	}
	return ""
}

// pathWithWriterPrefix conveniently assign /w prefix to path.
func (a *Agent) pathWithWriterPrefix(input interface{}) string {
	switch v := input.(type) {
	case resourced_config.Config:
		return "/w" + v.Path
	case string:
		if strings.HasPrefix(v, "/w/") {
			return v
		} else {
			return "/w" + v
		}
	}
	return ""
}

// Run executes a reader/writer config.
// Run will save reader data as JSON in local db.
func (a *Agent) Run(config resourced_config.Config) (output []byte, err error) {
	if config.Command != "" {
		output, err = a.runCommand(config)
	} else if config.GoStruct != "" && config.Kind == "reader" {
		output, err = a.runGoStructReader(config)
	} else if config.GoStruct != "" && config.Kind == "writer" {
		output, err = a.runGoStructWriter(config)
	}
	if err != nil {
		return output, err
	}

	err = a.saveRun(config, output)

	return output, err
}

// runCommand shells out external program and returns the output.
func (a *Agent) runCommand(config resourced_config.Config) ([]byte, error) {
	cmd := libprocess.NewCmd(config.Command)

	if config.Kind == "writer" {
		// Get readers data.
		readersData := make(map[string]interface{})

		for _, readerPath := range config.ReaderPaths {
			readerJsonBytes, err := a.GetRunByPath(a.pathWithReaderPrefix(readerPath))

			if err == nil {
				var data interface{}
				err := json.Unmarshal(readerJsonBytes, &data)
				if err == nil {
					readersData[readerPath] = data
				}
			}
		}

		readersDataJsonBytes, err := json.Marshal(readersData)
		if err != nil {
			return nil, err
		}

		cmd.Stdin = bytes.NewReader(readersDataJsonBytes)
	}

	return cmd.Output()
}

// processJson shells out external program to mangle JSON and save the new JSON on writer.
func (a *Agent) processJson(config resourced_config.Config, writer resourced_writers.IWriter) error {
	processorPath := writer.GetJsonProcessor()
	if processorPath == "" {
		return nil
	}

	cmd := libprocess.NewCmd(processorPath)

	readersData := writer.GetReadersData()

	readersDataJsonBytes, err := json.Marshal(readersData)
	if err != nil {
		return err
	}

	cmd.Stdin = bytes.NewReader(readersDataJsonBytes)

	newJsonDataBytes, err := cmd.Output()
	if err != nil {
		return err
	}

	var newJsonData map[string]interface{}
	err = json.Unmarshal(newJsonDataBytes, &newJsonData)
	if err != nil {
		println("is unmarshalling yield error? " + err.Error())
		return err
	}

	writer.SetReadersData(newJsonData)

	return err
}

// initGoStructReader initialize and return IReader.
func (a *Agent) initGoStructReader(config resourced_config.Config) (resourced_readers.IReader, error) {
	return resourced_readers.NewGoStructByConfig(config)
}

// initGoStructWriter initialize and return IWriter.
func (a *Agent) initGoStructWriter(config resourced_config.Config) (resourced_writers.IWriter, error) {
	writer, err := resourced_writers.NewGoStructByConfig(config)
	if err != nil {
		return nil, err
	}

	// Get readers data.
	readersData := make(map[string][]byte)

	for _, readerPath := range config.ReaderPaths {
		readerJsonBytes, err := a.GetRunByPath(a.pathWithReaderPrefix(readerPath))
		if err == nil {
			readersData[readerPath] = readerJsonBytes
		}
	}

	writer.SetReadersDataInBytes(readersData)

	return writer, err
}

// runGoStruct executes IReader/IWriter and returns the output.
// Note that IWriter also implements IReader
func (a *Agent) runGoStruct(readerOrWriter resourced_readers.IReader) ([]byte, error) {
	err := readerOrWriter.Run()
	if err != nil {
		errData := make(map[string]string)
		errData["Error"] = err.Error()
		return json.Marshal(errData)
	}

	return readerOrWriter.ToJson()
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
	// Initialize IWriter
	writer, err := a.initGoStructWriter(config)
	if err != nil {
		return nil, err
	}

	err = a.processJson(config, writer)
	if err != nil {
		return nil, err
	}

	return a.runGoStruct(writer)
}

// hostData builds host related information.
func (a *Agent) hostData() (*resourced_host.Host, error) {
	host, err := resourced_host.NewHostByHostname()
	if err != nil {
		return nil, err
	}

	host.Tags = a.Tags

	interfacesReader := resourced_readers.NewNetInterfaces()
	if interfacesReader.Run() == nil {
		host.NetworkInterfaces = make(map[string]map[string]interface{})

		for iface, stats := range interfacesReader.Data {
			host.NetworkInterfaces[iface] = make(map[string]interface{})
			host.NetworkInterfaces[iface]["HardwareAddress"] = stats.HardwareAddr
			host.NetworkInterfaces[iface]["IPAddresses"] = make([]string, len(stats.Addrs))

			for i, addr := range stats.Addrs {
				ipAddresses := host.NetworkInterfaces[iface]["IPAddresses"].([]string)
				ipAddresses[i] = addr.Addr
			}
		}
	}

	return host, nil
}

// saveRun gathers basic, host, and reader/witer information and save them into local storage.
func (a *Agent) saveRun(config resourced_config.Config, output []byte) error {
	// Do not perform save if config.Path is empty.
	if config.Path == "" {
		return nil
	}

	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = config.Path
	record["Interval"] = config.Interval

	if config.Command != "" {
		record["Command"] = config.Command
	}

	if config.GoStruct != "" {
		record["GoStruct"] = config.GoStruct
	}

	host, err := a.hostData()
	if err != nil {
		return err
	}
	record["Host"] = host

	runData := make(map[string]interface{})
	err = json.Unmarshal(output, &runData)
	if err != nil {
		return err
	}
	record["Data"] = runData

	recordInJson, err := json.Marshal(record)
	if err != nil {
		return err
	}

	err = a.Db.Update(func(tx *bolt.Tx) error {
		return a.dbBucket(tx).Put([]byte(a.pathWithPrefix(config)), recordInJson)
	})

	return err
}

// GetRun returns the JSON data stored in local storage given Config struct.
func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(a.pathWithPrefix(config))
}

// GetRunByPath returns JSON data stored in local storage given path string.
func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	var data []byte

	a.Db.View(func(tx *bolt.Tx) error {
		data = a.dbBucket(tx).Get([]byte(path))
		return nil
	})

	return data, nil
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
	for _, config := range a.ConfigStorage.Readers {
		a.RunForever(config)
	}
	for _, config := range a.ConfigStorage.Writers {
		a.RunForever(config)
	}
}
