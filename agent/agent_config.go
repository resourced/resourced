package agent

import (
	"os"

	resourced_config "github.com/resourced/resourced/config"
)

// setConfigStorage reads config paths and setup configStorage.
func (a *Agent) setConfigStorage() error {
	configDir := os.Getenv("RESOURCED_CONFIG_DIR")

	configStorage, err := resourced_config.NewConfigStorage(configDir)
	if err == nil {
		a.ConfigStorage = configStorage
	}

	return err
}
