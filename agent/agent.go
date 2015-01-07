package agent

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/go-fsnotify/fsnotify"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	resourced_readers "github.com/resourced/resourced/readers"
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

// setConfigStorage reads config paths and setup configStorage.
func (a *Agent) setConfigStorage() error {
	readerDir := os.Getenv("RESOURCED_CONFIG_READER_DIR")
	writerDir := os.Getenv("RESOURCED_CONFIG_WRITER_DIR")

	configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
	if err == nil {
		a.ConfigStorage = configStorage
	}

	go a.watchConfigDirectories(readerDir, writerDir)

	return err
}

// watchConfigDirectories uses inotify to watch changes on config directories.
func (a *Agent) watchConfigDirectories(readerDir, writerDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(readerDir)
	if err != nil {
		return err
	}

	err = watcher.Add(writerDir)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Remove == fsnotify.Remove) || (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) {
				fmt.Println("Config files changed. Rebuilding ConfigStorage...")

				configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
				if err == nil {
					a.ConfigStorage = configStorage
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				fmt.Printf("Error while watching config files: %v\n", err)
			}
		}
	}
	return nil
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

// Run executes a reader/writer and saves the result as JSON in local db.
func (a *Agent) Run(config resourced_config.Config) (output []byte, err error) {
	if config.Command != "" {
		output, err = a.runCommand(config)
	} else if config.GoStruct != "" {
		output, err = a.runGoStruct(config)
	}

	if err != nil {
		return nil, err
	}

	err = a.saveRun(config, output)
	return output, err
}

// runCommand shells out external program and returns the output.
func (a *Agent) runCommand(config resourced_config.Config) ([]byte, error) {
	cmd := libprocess.NewCmd(config.Command)
	return cmd.Output()
}

// runGoStruct executes IReaderWriter and returns the output.
func (a *Agent) runGoStruct(config resourced_config.Config) ([]byte, error) {
	var reader resourced_readers.IReaderWriter

	reader, err := resourced_readers.NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	err = reader.Run()
	if err != nil {
		return nil, err
	}

	return reader.ToJson()
}

// saveRun gathers default basic information and saves output into local storage.
func (a *Agent) saveRun(config resourced_config.Config, output []byte) error {
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
		return a.dbBucket(tx).Put([]byte(config.Path), recordInJson)
	})

	return err
}

// GetRun returns the JSON data stored in local storage given Config struct.
func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(config.Path)
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
