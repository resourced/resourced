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

// NewConfigStorage creates ConfigStorage struct given configReaderDir and configWriterDir.
func NewConfigStorage(configReaderDir, configWriterDir string) (*ConfigStorage, error) {
	storage := &ConfigStorage{}
	storage.Readers = make([]Config, 0)
	storage.Writers = make([]Config, 0)
	storage.Executors = make([]Config, 0)

	var err error

	if configReaderDir != "" {
		configReaderDir = libstring.ExpandTildeAndEnv(configReaderDir)

		readerFiles, err := ioutil.ReadDir(configReaderDir)

		if err == nil {
			for _, f := range readerFiles {
				fullpath := path.Join(configReaderDir, f.Name())

				readerConfig, err := NewConfig(fullpath, "reader")
				if err == nil {
					storage.Readers = append(storage.Readers, readerConfig)
				}
			}
		}
	}

	if configWriterDir != "" {
		configWriterDir = libstring.ExpandTildeAndEnv(configWriterDir)

		writerFiles, err := ioutil.ReadDir(configWriterDir)
		if err == nil {
			for _, f := range writerFiles {
				fullpath := path.Join(configWriterDir, f.Name())

				writerConfig, err := NewConfig(fullpath, "writer")
				if err == nil {
					storage.Writers = append(storage.Writers, writerConfig)
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
	LowTreshold  int64
	HighTreshold int64
	Conditions   []interface{}
}

// ConfigStorage stores all readers and writers configuration.
type ConfigStorage struct {
	Readers   []Config
	Writers   []Config
	Executors []Config
}
