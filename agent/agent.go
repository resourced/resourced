package agent

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	resourced_config "github.com/resourced/resourced/config"
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
func (a *Agent) setDb() error {
	var err error

	dbPath := os.Getenv("RESOURCED_DB")
	if dbPath == "" {
		dbPath = "~/resourced/db"

		err = os.MkdirAll(libstring.ExpandTilde("~/resourced"), 0755)
		if err != nil {
			return err
		}
	}

	a.DbPath = libstring.ExpandTilde(dbPath)

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

func (a *Agent) pathWithReaderPrefix(input interface{}) string {
	switch v := input.(type) {
	case resourced_config.Config:
		return "/r" + v.Path
	case string:
		return "/r" + v
	}
	return ""
}

func (a *Agent) pathWithWriterPrefix(input interface{}) string {
	switch v := input.(type) {
	case resourced_config.Config:
		return "/w" + v.Path
	case string:
		return "/w" + v
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
		// Get reader's run data.
		jsonBytes, err := a.GetRunByPath(a.pathWithReaderPrefix(config))
		if err != nil {
			return nil, err
		}

		cmd.Stdin = bytes.NewReader(jsonBytes)
	}

	return cmd.Output()
}

// runGoStructReader executes IReader and returns the output.
func (a *Agent) runGoStructReader(config resourced_config.Config) ([]byte, error) {
	reader, err := resourced_readers.NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	err = reader.Run()
	if err != nil {
		errData := make(map[string]string)
		errData["Error"] = err.Error()
		return json.Marshal(errData)
	}

	return reader.ToJson()
}

// runGoStructWriter executes IWriter and returns error if exists.
func (a *Agent) runGoStructWriter(config resourced_config.Config) ([]byte, error) {
	writer, err := resourced_writers.NewGoStruct(config.GoStruct)
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

	writer.SetReadersData(readersData)

	err = writer.Run()
	if err != nil {
		errData := make(map[string]string)
		errData["Error"] = err.Error()
		return json.Marshal(errData)
	}

	return writer.ToJson()
}

// saveRun gathers default basic information and saves output into local storage.
func (a *Agent) saveRun(config resourced_config.Config, output []byte) error {
	// Do not perform save if config.Path is empty.
	if config.Path == "" {
		return nil
	}

	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = config.Path
	record["Interval"] = config.Interval
	record["Tags"] = a.Tags

	if config.Command != "" {
		record["Command"] = config.Command
	}

	if config.GoStruct != "" {
		record["GoStruct"] = config.GoStruct
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	record["Hostname"] = hostname

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
