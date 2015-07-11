package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("DockerContainers", NewDockerContainers)
}

// NewDockerContainers is DockerContainers constructor.
func NewDockerContainers() readers.IReader {
	dc := &DockerContainers{}
	dc.Data = make(map[string]*libdocker.CompleteDockerContainer)
	return dc
}

// DockerContainers gathers docker containers data.
type DockerContainers struct {
	Data       map[string]*libdocker.CompleteDockerContainer
	DockerHost string
}

func (dc *DockerContainers) Run() error {
	containers, err := libdocker.AllInspectedContainers(dc.DockerHost)
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" && container.NiceImageName != "" {
			dc.Data[container.NiceImageName+"-"+container.ID] = container
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (dc *DockerContainers) ToJson() ([]byte, error) {
	return json.Marshal(dc.Data)
}
