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

	if config.Interval == "" {
		config.Interval = "1m"
	}

	config.Kind = kind

	return config, err
}

// NewConfigs creates Configs struct given configDir.
func NewConfigs(configDir string) (*Configs, error) {
	storage := &Configs{}
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
				fullpath := path.Join(configDir, configKindPlural, f.Name())

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
	Conditions    string
}

// Configs stores all readers and writers configuration.
type Configs struct {
	Readers   []Config
	Writers   []Config
	Executors []Config
}

// NewGeneralConfig is the constructor for GeneralConfig.
func NewGeneralConfig(configDir string) (GeneralConfig, error) {
	configDir = libstring.ExpandTildeAndEnv(configDir)
	fullpath := path.Join(configDir, "general.toml")

	var config GeneralConfig
	_, err := toml.DecodeFile(fullpath, &config)

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config, err
}

// GeneralConfig stores all other configuration data.
type GeneralConfig struct {
	Addr            string
	LogLevel        string
	AllowedNetworks string

	HTTPS struct {
		CertFile string
		KeyFile  string
	}
	ResourcedMaster struct {
		Url         string
		AccessToken string
	}
}
