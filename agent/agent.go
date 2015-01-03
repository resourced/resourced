package agent

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	resourced_readers "github.com/resourced/resourced/readers"
	"os"
	"time"
)

// NewAgent is the constructor fot Agent struct.
func NewAgent() (*Agent, error) {
	agent := &Agent{}

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
}

// setConfigStorage reads config paths and setup configStorage.
func (a *Agent) setConfigStorage() error {
	readerDir := os.Getenv("RESOURCED_CONFIG_READER_DIR")
	writerDir := os.Getenv("RESOURCED_CONFIG_WRITER_DIR")

	configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
	if err == nil {
		a.ConfigStorage = configStorage
	}
	return err
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

// DbBucket returns the boltdb bucket.
func (a *Agent) DbBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("resources"))
}

// Run executes a particular reader/writer config and saves the result on local storage.
func (a *Agent) Run(config resourced_config.Config) (output []byte, err error) {
	if config.Command != "" {
		output, err = a.RunCommand(config)
	} else if config.GoStruct != "" {
		output, err = a.RunGoStruct(config)
	}

	if err != nil {
		return nil, err
	}

	err = a.SaveRun(config, output)
	return output, err
}

// RunCommand shells out external program and returns the output.
func (a *Agent) RunCommand(config resourced_config.Config) ([]byte, error) {
	cmd := libprocess.NewCmd(config.Command)
	return cmd.Output()
}

// RunCommand executes IReaderWriter and returns the output.
func (a *Agent) RunGoStruct(config resourced_config.Config) ([]byte, error) {
	var reader resourced_readers.IReaderWriter

	// TODO(didip): We need reflection, this is so ghetto!
	if config.GoStruct == "NetworkInterfaces" {
		reader = resourced_readers.NewNetworkInterfaces()
	} else if config.GoStruct == "Df" {
		reader = resourced_readers.NewDf()
	} else if config.GoStruct == "Memory" {
		reader = resourced_readers.NewMemory()
	} else if config.GoStruct == "Ps" {
		reader = resourced_readers.NewPs()
	} else if config.GoStruct == "LoadAvg" {
		reader = resourced_readers.NewLoadAvg()
	} else if config.GoStruct == "Uptime" {
		reader = resourced_readers.NewUptime()
	} else if config.GoStruct == "Meminfo" {
		reader = resourced_readers.NewMeminfo()
	} else {
		err := errors.New("GoStruct is undefined.")
		return nil, err
	}

	err := reader.Run()
	if err != nil {
		return nil, err
	}

	return reader.ToJson()
}

// SaveRun gathers default basic information and saves output into local storage.
func (a *Agent) SaveRun(config resourced_config.Config, output []byte) error {
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
		return a.DbBucket(tx).Put([]byte(config.Path), recordInJson)
	})

	return err
}

// GetRun returns JSON data stored in local storage given Config struct.
func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(config.Path)
}

// GetRunByPath returns JSON data stored in local storage given path string.
func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	var data []byte

	a.Db.View(func(tx *bolt.Tx) error {
		data = a.DbBucket(tx).Get([]byte(path))
		return nil
	})

	return data, nil
}

// RunForever executes Run() in a loop with a sleep of config.Interval.
func (a *Agent) RunForever(config resourced_config.Config) {
	go func(a *Agent, config resourced_config.Config) {
		for {
			a.Run(config)
			libtime.SleepString(config.Interval)
		}
	}(a, config)
}

// RunAllForever executes all readers & writers in a loop.
func (a *Agent) RunAllForever() {
	for _, config := range a.ConfigStorage.Readers {
		a.RunForever(config)
	}
	for _, config := range a.ConfigStorage.Writers {
		a.RunForever(config)
	}
}
