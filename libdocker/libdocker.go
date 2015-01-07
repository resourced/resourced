package libdocker

import (
	dockerclient "github.com/fsouza/go-dockerclient"
	"os"
	"path"
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
