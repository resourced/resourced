package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
)

func NewDockerInfoVersion() *DockerInfoVersion {
	dc := &DockerInfoVersion{}
	return dc
}

// DockerInfoVersion gathers docker containers data.
type DockerInfoVersion struct {
	Data       map[string]interface{}
	DockerHost string
}

// Run fetches info and version data.
func (dc *DockerInfoVersion) Run() error {
	infoAndVersion, err := libdocker.InfoAndVersion(dc.DockerHost)
	if err != nil {
		return nil
	}

	dc.Data = infoAndVersion

	return nil
}

// ToJson serialize Data field to JSON.
func (dc *DockerInfoVersion) ToJson() ([]byte, error) {
	return json.Marshal(dc.Data)
}
