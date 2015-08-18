package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/resourced/resourced/libstring"
)

// NewDefaultConfigs provide default config setup.
// This function is called on first boot.
func NewDefaultConfigs(configDir string) error {
	configDir = libstring.ExpandTildeAndEnv(configDir)

	// Create configDir if it does not exist
	if _, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(configDir, 0755)
			if err != nil {
				return err
			}
		}
	}

	// Create subdirectories
	for _, subdirConfigs := range []string{"readers", "writers", "executors", "tags"} {
		subdirPath := path.Join(configDir, subdirConfigs)

		if _, err := os.Stat(subdirPath); err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(subdirPath, 0755)
				if err != nil {
					return err
				}
			}
		}
	}

	// Create default tags
	defaultTagsTemplate := `GOOS=%v
uname=%v
`

	unameBytes, err := exec.Command("uname", "-a").CombinedOutput()
	if err != nil {
		return err
	}

	uname := strings.TrimSpace(string(unameBytes))

	err = ioutil.WriteFile(path.Join(configDir, "tags", "default"), []byte(fmt.Sprintf(defaultTagsTemplate, runtime.GOOS, uname)), 0755)
	if err != nil {
		return err
	}

	// Create a default general.toml
	generalToml := `# Addr is the host and port of ResourceD Agent HTTP/S server
Addr = "localhost:55555"

# Valid LogLevel are: debug, info, warning, error, fatal, panic
LogLevel = "info"

[HTTPS]
CertFile = ""
KeyFile = ""

[ResourcedMaster]
# Url is the root endpoint to Resourced Master
URL = "http://localhost:55655"

# General purpose AccessToken, it will be used when AccessToken is not defined elsewhere.
AccessToken = "{access-token}"
`

	err = ioutil.WriteFile(path.Join(configDir, "general.toml"), []byte(generalToml), 0644)
	if err != nil {
		return err
	}

	return nil
}
