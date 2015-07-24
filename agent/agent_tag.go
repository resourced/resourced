package agent

import (
	"bufio"
	"github.com/resourced/resourced/libstring"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (a *Agent) setTags() error {
	a.Tags = make(map[string]string)

	configDir := os.Getenv("RESOURCED_CONFIG_DIR")
	if configDir == "" {
		return nil
	}

	configDir = libstring.ExpandTildeAndEnv(configDir)

	tagFiles, err := ioutil.ReadDir(path.Join(configDir, "tags"))
	if err != nil {
		return nil
	}

	for _, f := range tagFiles {
		fullpath := path.Join(configDir, "tags", f.Name())

		file, err := os.Open(fullpath)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tagsPerLine := strings.Split(scanner.Text(), ",")
			for _, tagKeyValue := range tagsPerLine {
				keyValue := strings.Split(tagKeyValue, "=")
				if len(keyValue) >= 2 {
					a.Tags[keyValue[0]] = strings.Join(keyValue[1:], "=")
				}
			}
		}
	}

	return nil
}
