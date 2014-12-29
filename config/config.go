package config

import (
	"github.com/BurntSushi/toml"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"io/ioutil"
	"path"
)

func NewConfig(fullpath string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(fullpath, &config)

	if config.Command != "" {
		config.Command = libstring.ExpandTilde(config.Command)
	}

	return config, err
}

func NewConfigStorage(configReaderDir, configWriterDir string) (*ConfigStorage, error) {
	storage := &ConfigStorage{}
	storage.Readers = make([]Config, 0)
	storage.Writers = make([]Config, 0)

	var err error

	if configReaderDir != "" {
		configReaderDir = libstring.ExpandTilde(configReaderDir)
		readerFiles, err := ioutil.ReadDir(configReaderDir)

		if err == nil {
			for _, f := range readerFiles {
				fullpath := path.Join(configReaderDir, f.Name())

				readerConfig, err := NewConfig(fullpath)
				if err == nil {
					storage.Readers = append(storage.Readers, readerConfig)
				}
			}
		}
	}

	if configWriterDir != "" {
		configWriterDir = libstring.ExpandTilde(configWriterDir)
		writerFiles, err := ioutil.ReadDir(configWriterDir)
		if err == nil {
			for _, f := range writerFiles {
				fullpath := path.Join(configWriterDir, f.Name())

				writerConfig, err := NewConfig(fullpath)
				if err == nil {
					storage.Writers = append(storage.Writers, writerConfig)
				}
			}
		}
	}

	return storage, err
}

type Config struct {
	Command  string
	Path     string
	Interval string
}

func (c *Config) Run() ([]byte, error) {
	processWrapper := libprocess.NewProcessWrapper(c.Command)
	return processWrapper.NewCmd().Output()
}

type ConfigStorage struct {
	Readers []Config
	Writers []Config
}
