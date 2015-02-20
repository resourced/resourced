package docker

import (
	"encoding/json"
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/resourced/resourced/libdocker"
)

func NewDockerContainers() *DockerContainers {
	dc := &DockerContainers{}
	dc.Data = make(map[string]*dockerclient.Container)
	return dc
}

// DockerContainers gathers docker containers data.
type DockerContainers struct {
	Data       map[string]*dockerclient.Container
	DockerHost string
}

func (dc *DockerContainers) Run() error {
	containers, err := libdocker.AllInspectedContainers(dc.DockerHost)
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" && container.Config != nil && container.Config.Image != "" {
			dc.Data[container.Config.Image+"-"+container.ID] = container
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (dc *DockerContainers) ToJson() ([]byte, error) {
	return json.Marshal(dc.Data)
}
