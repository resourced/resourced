// Package config provides data structure for storing resourced reader/writer configurations.
package config

import (
	"io/ioutil"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/resourced/resourced/libstring"
)

// NewConfig creates Config struct given fullpath and kind.
func NewConfig(fullpath, kind string) (Config, error) {
	fullpath = libstring.ExpandTildeAndEnv(fullpath)

	var config Config
	_, err := toml.DecodeFile(fullpath, &config)

	if config.Command != "" {
		config.Command = libstring.ExpandTildeAndEnv(config.Command)
	}

	if config.Interval == "" {
		config.Interval = "1m"
	}

	config.Kind = kind

	return config, err
}

// NewConfigStorage creates ConfigStorage struct given configDir.
func NewConfigStorage(configDir string) (*ConfigStorage, error) {
	storage := &ConfigStorage{}
	storage.Readers = make([]Config, 0)
	storage.Writers = make([]Config, 0)
	storage.Executors = make([]Config, 0)

	var err error

	for _, configKind := range []string{"reader", "writer", "executor"} {
		configDir = libstring.ExpandTildeAndEnv(configDir)

		configKindPlural := configKind + "s"

		configFiles, err := ioutil.ReadDir(path.Join(configDir, configKindPlural))

		if err == nil {
			for _, f := range configFiles {
				fullpath := path.Join(path.Join(configDir, configKindPlural), f.Name())

				conf, err := NewConfig(fullpath, configKind)
				if err == nil {
					if configKind == "reader" {
						storage.Readers = append(storage.Readers, conf)
					}
					if configKind == "writer" {
						storage.Writers = append(storage.Writers, conf)
					}
					if configKind == "executor" {
						storage.Executors = append(storage.Executors, conf)
					}
				} else {
					println(err.Error())
				}
			}
		}
	}

	return storage, err
}

// Config is a unit of execution for a reader/writer.
// Reader config defines how to fetch a particular information and its JSON data path.
// Writer config defines how to export the JSON data to a particular destination. E.g. Facts/graphing database.
type Config struct {
	Command        string
	GoStruct       string
	GoStructFields map[string]interface{}
	Path           string
	Interval       string

	// There are 3 kinds: reader, writer, and executor
	Kind string

	// Writer specific fields
	// ReaderPaths defines input data endpoints for a Writer.
	ReaderPaths []string

	// Executor specific fields
	LowThreshold  int64
	HighThreshold int64
	Conditions    []interface{}
}

// ConfigStorage stores all readers and writers configuration.
type ConfigStorage struct {
	Readers   []Config
	Writers   []Config
	Executors []Config
}
