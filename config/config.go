package config

import (
	"github.com/BurntSushi/toml"
	"github.com/resourced/resourced/libstring"
	"io/ioutil"
	"os"
	"path"
)

// NewConfig creates Config struct given fullpath and kind.
func NewConfig(fullpath, kind string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(fullpath, &config)

	if config.Command != "" {
		config.Command = libstring.ExpandTilde(config.Command)
		config.Command = os.ExpandEnv(config.Command)
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

	var err error

	if configReaderDir != "" {
		configReaderDir = libstring.ExpandTilde(configReaderDir)
		configReaderDir = os.ExpandEnv(configReaderDir)

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
		configWriterDir = libstring.ExpandTilde(configWriterDir)
		configWriterDir = os.ExpandEnv(configWriterDir)

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
	Command  string
	GoStruct string
	Path     string
	Interval string

	// There are only 2 kinds: reader and writer
	Kind string
}

// ConfigStorage stores all readers and writers configuration.
type ConfigStorage struct {
	Readers []Config
	Writers []Config
}
