// Package config provides data structure for storing resourced reader/writer configurations.
package config

import (
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libstring"
)

// NewConfigs creates Configs struct given configDir.
func NewConfigs(configDir string) (*Configs, error) {
	configs := &Configs{}
	configs.Readers = make([]Config, 0)
	configs.Writers = make([]Config, 0)
	configs.Executors = make([]Config, 0)
	configs.Loggers = make([]Config, 0)

	var err error

	configDir = libstring.ExpandTildeAndEnv(configDir)

	for _, configKind := range []string{"reader", "writer", "executor", "logger"} {
		configKindPlural := configKind + "s"

		configFiles, err := ioutil.ReadDir(path.Join(configDir, configKindPlural))

		if err == nil {
			for _, f := range configFiles {
				fullpath := path.Join(configDir, configKindPlural, f.Name())

				conf, err := NewConfig(fullpath, configKind)
				if err == nil {
					if configKind == "reader" {
						configs.Readers = append(configs.Readers, conf)
					}
					if configKind == "writer" {
						configs.Writers = append(configs.Writers, conf)
					}
					if configKind == "executor" {
						configs.Executors = append(configs.Executors, conf)
					}
					if configKind == "logger" {
						configs.Loggers = append(configs.Loggers, conf)
					}
				}
			}
		}
	}

	return configs, err
}

// Configs stores all readers, writers, and executors configuration.
type Configs struct {
	Readers   []Config
	Writers   []Config
	Executors []Config
	Loggers   []Config
}

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

// Config is a unit of execution for a reader/writer.
// Reader config defines how to fetch a particular information and its JSON data path.
// Writer config defines how to export the JSON data to a particular destination. E.g. Facts/graphing database.
type Config struct {
	GoStruct       string
	GoStructFields map[string]interface{}
	Path           string
	Interval       string
	Host           *host.Host

	// There are 4 kinds: reader, writer, executor, and log
	Kind string

	// Writer specific fields
	// ReaderPaths defines input data endpoints for a Writer.
	ReaderPaths []string

	// Executor specific fields
	LowThreshold               int64
	HighThreshold              int64
	Conditions                 string
	ResourcedMasterURL         string
	ResourcedMasterAccessToken string
}

// CommonJsonData returns common information for every reader/writer/executor JSON interpretation.
func (c *Config) CommonJsonData() map[string]interface{} {
	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = c.Path

	if c.Interval == "" {
		c.Interval = "1m"
	}
	record["Interval"] = c.Interval

	if c.GoStruct != "" {
		record["GoStruct"] = c.GoStruct
	}

	return record
}

// PathWithPrefix prepends the short version of config.Kind to path.
func (c *Config) PathWithPrefix() string {
	if c.Kind == "reader" {
		return c.PathWithKindPrefix("r", "")
	} else if c.Kind == "writer" {
		return c.PathWithKindPrefix("w", "")
	} else if c.Kind == "executor" {
		return c.PathWithKindPrefix("x", "")
	} else if c.Kind == "logger" {
		return c.PathWithKindPrefix("logs", "")
	}
	return c.Path
}

// PathWithKindPrefix is common prepender function
func (c *Config) PathWithKindPrefix(kind string, input string) string {
	prefix := "/" + kind

	if input != "" {
		if strings.HasPrefix(input, prefix+"/") {
			return input
		} else {
			return prefix + input
		}

	} else {
		return prefix + c.Path
	}

	return ""
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

type ITCPServer interface {
	GetAddr() string
	GetCertFile() string
	GetKeyFile() string
}

type TCPConfig struct {
	Addr     string
	CertFile string
	KeyFile  string
}

func (c TCPConfig) GetAddr() string {
	return c.Addr
}

func (c TCPConfig) GetCertFile() string {
	return c.CertFile
}

func (c TCPConfig) GetKeyFile() string {
	return c.KeyFile
}

type GraphiteConfig struct {
	TCPConfig
	StatsInterval string
	Blacklist     []string
}

type LogReceiverConfig struct {
	TCPConfig
	WriteToMasterInterval string
	AutoPruneLength       int64
}

func (l LogReceiverConfig) GetAutoPruneLength() int64 {
	return l.AutoPruneLength
}

// GeneralConfig stores all other configuration data.
type GeneralConfig struct {
	Addr     string
	LogLevel string
	TTL      int

	HTTPS struct {
		CertFile string
		KeyFile  string
	}
	ResourcedMaster struct {
		URL         string
		AccessToken string
	}
	Graphite    GraphiteConfig
	LogReceiver LogReceiverConfig
}
