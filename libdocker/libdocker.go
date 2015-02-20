// Package libdocker provides docker related library functions.
package libdocker

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	"os"
	"path"
	"sync"
)

// DockerClient returns dockerclient.Client which handles Docker connection.
func DockerClient() (*dockerclient.Client, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	if dockerCertPath != "" {
		cert := path.Join(dockerCertPath, "cert.pem")
		key := path.Join(dockerCertPath, "key.pem")
		ca := path.Join(dockerCertPath, "ca.pem")

		return dockerclient.NewTLSClient(dockerHost, cert, key, ca)
	} else {
		return dockerclient.NewClient(dockerHost)
	}
}

// AllContainers is a convenience function to fetch a slice of all containers data.
func AllContainers() ([]dockerclient.APIContainers, error) {
	client, err := DockerClient()
	if err != nil {
		return nil, err
	}

	return client.ListContainers(dockerclient.ListContainersOptions{})
}

// AllInspectedContainers is a convenience function to fetch a slice of all inspected containers data.
func AllInspectedContainers() ([]*dockerclient.Container, error) {
	client, err := DockerClient()
	if err != nil {
		return nil, err
	}

	shortDescContainers, err := client.ListContainers(dockerclient.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	containersChan := make(chan *dockerclient.Container)
	var wg sync.WaitGroup

	for _, shortDescContainer := range shortDescContainers {
		wg.Add(1)

		go func(shortDescContainer dockerclient.APIContainers) {
			defer wg.Done()

			fullDescContainer, err := client.InspectContainer(shortDescContainer.ID)
			if err == nil && fullDescContainer != nil {
				containersChan <- fullDescContainer
			}
		}(shortDescContainer)
	}

	containers := make([]*dockerclient.Container, 0)

	go func() {
		for container := range containersChan {
			containers = append(containers, container)
		}
	}()

	wg.Wait()

	return containers, nil
}

// AllImages is a convenience function to fetch a slice of all images data.
func AllImages() ([]dockerclient.APIImages, error) {
	client, err := DockerClient()
	if err != nil {
		return nil, err
	}

	return client.ListImages(dockerclient.ListImagesOptions{})
}

// AllInspectedImages is a convenience function to fetch a slice of all inspected images data.
func AllInspectedImages() ([]*dockerclient.Image, error) {
	client, err := DockerClient()
	if err != nil {
		return nil, err
	}

	shortDescImages, err := client.ListImages(dockerclient.ListImagesOptions{})
	if err != nil {
		return nil, err
	}

	imagesChan := make(chan *dockerclient.Image)
	var wg sync.WaitGroup

	for _, shortDescImage := range shortDescImages {
		wg.Add(1)

		go func(shortDescImage dockerclient.APIImages) {
			defer wg.Done()

			fullDescImage, err := client.InspectImage(shortDescImage.ID)
			if err == nil && fullDescImage != nil {
				imagesChan <- fullDescImage
			}
		}(shortDescImage)
	}

	images := make([]*dockerclient.Image, 0)

	go func() {
		for image := range imagesChan {
			images = append(images, image)
		}
	}()

	wg.Wait()

	return images, nil
}
