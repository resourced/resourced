package agent

import (
	resourced_config "github.com/resourced/resourced/config"
	"os"
)

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
