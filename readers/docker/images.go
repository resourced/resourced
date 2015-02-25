package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
)

func NewDockerImages() *DockerImages {
	di := &DockerImages{}
	di.Data = make(map[string]*libdocker.CompleteDockerImage)
	return di
}

// DockerImages gathers docker images data.
type DockerImages struct {
	Data       map[string]*libdocker.CompleteDockerImage
	DockerHost string
}

func (di *DockerImages) Run() error {
	images, err := libdocker.AllInspectedImages(di.DockerHost)
	if err != nil {
		return nil
	}

	for _, image := range images {
		if image.ID != "" {
			if len(image.RepoTags) > 0 {
				di.Data[image.RepoTags[0]+"-"+image.ID] = image
			} else {
				di.Data[image.ID] = image
			}
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (di *DockerImages) ToJson() ([]byte, error) {
	return json.Marshal(di.Data)
}
