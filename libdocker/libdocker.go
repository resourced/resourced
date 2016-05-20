// Package libdocker provides docker related library functions.
package libdocker

import (
	"os"
	"path"
	"strconv"
	"sync"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/resourced/resourced/libstring"
)

var connections map[string]*dockerclient.Client
var connectionsLock = &sync.RWMutex{}

type CompleteDockerContainer struct {
	NiceImageName string `json:"NiceImageName,omitempty" yaml:"NiceImageName,omitempty"`
	Command       string `json:"Command,omitempty" yaml:"Command,omitempty"`
	Status        string `json:"Status,omitempty" yaml:"Status,omitempty"`
	dockerclient.Container
}

type CompleteDockerImage struct {
	RepoTags    []string `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
	VirtualSize int64    `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty"`
	ParentID    string   `json:"ParentId,omitempty" yaml:"ParentId,omitempty"`
	dockerclient.Image
}

// DockerClient returns dockerclient.Client which handles Docker connection.
func DockerClient(endpoint string) (*dockerclient.Client, error) {
	var conn *dockerclient.Client
	var err error

	if endpoint == "" {
		endpoint = os.Getenv("DOCKER_HOST")
		if endpoint == "" {
			endpoint = "unix:///var/run/docker.sock"
		}
	}

	connectionsLock.Lock()
	if connections == nil {
		connections = make(map[string]*dockerclient.Client)
	}
	connectionsLock.Unlock()

	// Do not create connection if one already exist.
	connectionsLock.RLock()
	existingConnection, ok := connections[endpoint]
	connectionsLock.RUnlock()

	if ok && existingConnection != nil {
		return existingConnection, nil
	}

	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	if dockerCertPath != "" {
		cert := path.Join(dockerCertPath, "cert.pem")
		key := path.Join(dockerCertPath, "key.pem")
		ca := path.Join(dockerCertPath, "ca.pem")

		conn, err = dockerclient.NewTLSClient(endpoint, cert, key, ca)
	} else {
		conn, err = dockerclient.NewClient(endpoint)
	}

	if err == nil && conn != nil {
		connectionsLock.Lock()
		connections[endpoint] = conn
		connectionsLock.Unlock()
	}

	return conn, err
}

// InfoAndVersion is a convenience function to fetch info and version data.
func InfoAndVersion(endpoint string) (map[string]interface{}, error) {
	client, err := DockerClient(endpoint)
	if err != nil {
		return nil, err
	}

	version, err := client.Version()
	if err != nil {
		return nil, err
	}

	versionAsMap := version.Map()
	infoAsMap := version.Map()

	data := make(map[string]interface{})

	for key, value := range versionAsMap {
		data[key] = value
	}

	data["Driver"] = make(map[string]interface{})

	for key, value := range infoAsMap {
		if libstring.StringInSlice(key, []string{"NGoroutines", "Containers", "Images", "MemTotal"}) {
			data[key] = version.GetInt64(key)

		} else if key == "NFd" {
			data["NumFileDescriptors"] = version.GetInt64(key)

		} else if key == "NEventsListener" {
			data["NumEventsListeners"] = version.GetInt64(key)

		} else if key == "NCPU" {
			data["NumCPUs"] = version.GetInt64(key)

		} else if libstring.StringInSlice(key, []string{"Debug", "IPv4Forwarding", "MemoryLimit", "SwapLimit"}) {
			data[key] = version.GetBool(key)

		} else if key == "Driver" {
			driverMap := data["Driver"].(map[string]interface{})
			driverMap["Name"] = value

		} else if key == "DriverStatus" {
			tupleSlice := make([][]string, 2)
			version.GetJSON(key, &tupleSlice)

			for _, tuple := range tupleSlice {
				tupleKey := tuple[0]
				tupleValue := tuple[1]

				driverMap := data["Driver"].(map[string]interface{})

				if tupleKey == "Root Dir" {
					driverMap["RootDir"] = tupleValue
				}
				if tupleKey == "Dirs" {
					tupleValueInt64, err := strconv.ParseInt(tupleValue, 10, 64)
					if err == nil {
						driverMap[tupleKey] = tupleValueInt64
					}
				}
			}

		} else if key == "RegistryConfig" {
			registryConfig := make(map[string]interface{})
			err := version.GetJSON(key, &registryConfig)
			if err == nil {
				data[key] = registryConfig
			}

		} else {
			data[key] = value
		}
	}

	return data, nil
}

// AllContainers is a convenience function to fetch a slice of all containers data.
func AllContainers(endpoint string) ([]dockerclient.APIContainers, error) {
	client, err := DockerClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.ListContainers(dockerclient.ListContainersOptions{})
}

// AllInspectedContainers is a convenience function to fetch a slice of all inspected containers data.
func AllInspectedContainers(endpoint string) ([]*CompleteDockerContainer, error) {
	client, err := DockerClient(endpoint)
	if err != nil {
		return nil, err
	}

	shortDescContainers, err := client.ListContainers(dockerclient.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	containersChan := make(chan *CompleteDockerContainer)
	var wg sync.WaitGroup

	for _, shortDescContainer := range shortDescContainers {
		container := &CompleteDockerContainer{}
		container.NiceImageName = shortDescContainer.Image
		container.Command = shortDescContainer.Command
		container.Status = shortDescContainer.Status

		wg.Add(1)

		go func(container *CompleteDockerContainer) {
			defer wg.Done()

			fullDescContainer, err := client.InspectContainer(shortDescContainer.ID)
			if err == nil && fullDescContainer != nil {
				container.ID = fullDescContainer.ID
				container.Created = fullDescContainer.Created
				container.Path = fullDescContainer.Path
				container.Args = fullDescContainer.Args
				container.Config = fullDescContainer.Config
				container.State = fullDescContainer.State
				container.Image = fullDescContainer.Image
				container.NetworkSettings = fullDescContainer.NetworkSettings
				container.SysInitPath = fullDescContainer.SysInitPath
				container.ResolvConfPath = fullDescContainer.ResolvConfPath
				container.HostnamePath = fullDescContainer.HostnamePath
				container.HostsPath = fullDescContainer.HostsPath
				container.Name = fullDescContainer.Name
				container.Driver = fullDescContainer.Driver
				container.Volumes = fullDescContainer.Volumes
				container.VolumesRW = fullDescContainer.VolumesRW
				container.HostConfig = fullDescContainer.HostConfig

				containersChan <- container
			}
		}(container)
	}

	containers := make([]*CompleteDockerContainer, 0)

	go func() {
		for container := range containersChan {
			containers = append(containers, container)
		}
	}()

	wg.Wait()
	close(containersChan)

	return containers, nil
}

// AllImages is a convenience function to fetch a slice of all images data.
func AllImages(endpoint string) ([]dockerclient.APIImages, error) {
	client, err := DockerClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.ListImages(dockerclient.ListImagesOptions{})
}

// AllInspectedImages is a convenience function to fetch a slice of all inspected images data.
func AllInspectedImages(endpoint string) ([]*CompleteDockerImage, error) {
	client, err := DockerClient(endpoint)
	if err != nil {
		return nil, err
	}

	shortDescImages, err := client.ListImages(dockerclient.ListImagesOptions{})
	if err != nil {
		return nil, err
	}

	imagesChan := make(chan *CompleteDockerImage)
	var wg sync.WaitGroup

	for _, shortDescImage := range shortDescImages {
		img := &CompleteDockerImage{}
		img.ID = shortDescImage.ID
		img.RepoTags = shortDescImage.RepoTags
		img.VirtualSize = shortDescImage.VirtualSize
		img.ParentID = shortDescImage.ParentID

		wg.Add(1)

		go func(img *CompleteDockerImage) {
			defer wg.Done()

			fullDescImage, err := client.InspectImage(img.ID)
			if err == nil && fullDescImage != nil {
				img.Parent = fullDescImage.Parent
				img.Comment = fullDescImage.Comment
				img.Created = fullDescImage.Created
				img.Container = fullDescImage.Container
				img.ContainerConfig = fullDescImage.ContainerConfig
				img.DockerVersion = fullDescImage.DockerVersion
				img.Author = fullDescImage.Author
				img.Config = fullDescImage.Config
				img.Architecture = fullDescImage.Architecture
				img.Size = fullDescImage.Size

				imagesChan <- img
			}
		}(img)
	}

	images := make([]*CompleteDockerImage, 0)

	go func() {
		for image := range imagesChan {
			images = append(images, image)
		}
	}()

	wg.Wait()
	close(imagesChan)

	return images, nil
}
