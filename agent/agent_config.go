package agent

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
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

	go a.watchConfigDirectories(readerDir, writerDir)

	return err
}

// watchConfigDirectories uses inotify to watch changes on config directories.
func (a *Agent) watchConfigDirectories(readerDir, writerDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(readerDir)
	if err != nil {
		return err
	}

	err = watcher.Add(writerDir)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Remove == fsnotify.Remove) || (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) {
				fmt.Println("Config files changed. Rebuilding ConfigStorage...")

				configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
				if err == nil {
					a.ConfigStorage = configStorage
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				fmt.Printf("Error while watching config files: %v\n", err)
			}
		}
	}
	return nil
}
