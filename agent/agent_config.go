package agent

import (
	"os"

	resourced_config "github.com/resourced/resourced/config"
)

// setConfigs reads config paths and setup configStorage.
func (a *Agent) setConfigs() error {
	configDir := os.Getenv("RESOURCED_CONFIG_DIR")

	configStorage, err := resourced_config.NewConfigs(configDir)
	if err != nil {
		return err
	}
	a.Configs = configStorage

	generalConfig, err := resourced_config.NewGeneralConfig(configDir)
	if err != nil {
		return err
	}
	a.GeneralConfig = generalConfig

	return err
}
