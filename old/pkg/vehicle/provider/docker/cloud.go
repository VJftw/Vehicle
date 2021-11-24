package docker

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/network"

	"github.com/VJftw/vehicle/pkg/vehicle/provider"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
)

// Docker represents the Docker vehicle instance provider
type Docker struct {
	*provider.Base

	dockerClient *client.Client

	containerID string
	networkID   string
	sshPort     uint16

	privateKey      *rsa.PrivateKey
	sshClientConfig *ssh.ClientConfig
}

// New returns a new Docker vehicle instance
func New(uuid string) *Docker {
	dockerClient, _ := client.NewEnvClient()

	return &Docker{
		Base:         &provider.Base{Provider: "docker", UUID: uuid},
		dockerClient: dockerClient,
	}
}

// GetSSHInfo returns the ssh configuration for a provisioned Docker instance
func (d *Docker) GetSSHInfo() (string, uint16, time.Duration, *ssh.ClientConfig) {
	return "127.0.0.1", d.sshPort, 20 * time.Second, d.sshClientConfig
}

// ResolveFuncs returns the functions involved in resolving dynamic resource IDs
func (d *Docker) ResolveFuncs() []func() (error, []string) {
	return []func() (error, []string){}
}

// StartFuncs returns the functions involved in starting an Docker instance
func (d *Docker) StartFuncs() []func() error {
	return []func() error{
		d.createPrivateKey,
		d.generateSSHConfig,
		d.launchInstance,
	}
}

// Stop stops and all Docker resources created by vehicle
func (d *Docker) Stop() error {
	if d.containerID != "" {
		stopTimeout, _ := time.ParseDuration("1s")
		err := d.dockerClient.ContainerStop(
			context.Background(),
			d.containerID,
			&stopTimeout,
		)
		if err != nil {
			return err
		}
		err = d.dockerClient.ContainerRemove(
			context.Background(),
			d.containerID,
			types.ContainerRemoveOptions{RemoveVolumes: true},
		)
		if err != nil {
			return err
		}
	}
	// if len(d.privateKey) > 0 {

	// }
	return nil
}

func (d *Docker) createPrivateKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	d.privateKey = privateKey

	return nil
}

func (d *Docker) generateSSHConfig() error {
	fmt.Println("generating ssh config")
	// privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(d.privateKey)}
	// var bufferedPrivateKey bytes.Buffer
	// privateKeyWriter := bufio.NewWriter(&bufferedPrivateKey)
	// if err := pem.Encode(privateKeyWriter, privateKeyPEM); err != nil {
	// 	return err
	// }

	// signer, err := ssh.ParsePrivateKey(bufferedPrivateKey.Bytes())
	// if err != nil {
	// 	return err
	// }

	d.sshClientConfig = &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("vehicle"),
			// ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(keyPairOut.)
	}
	return nil
}

func (d *Docker) launchInstance() error {
	fmt.Println("launching instance")

	containerConfig := &container.Config{
		Image: "c83e7a004c1b",
		// Cmd:   []string(dR.Command),
		// Env: []string{},
	}

	createResp, err := d.dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		&container.HostConfig{
			PublishAllPorts: true,
		},
		&network.NetworkingConfig{},
		"vehicle-name",
	)

	if err != nil {
		return err
	}

	d.containerID = createResp.ID

	err = d.dockerClient.ContainerStart(
		context.Background(),
		d.containerID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	inspect, err := d.dockerClient.ContainerInspect(context.Background(), d.containerID)
	if err != nil {
		return err
	}

	sshPortStr := inspect.NetworkSettings.NetworkSettingsBase.Ports["22/tcp"][0].HostPort
	sshPort64, err := strconv.ParseUint(sshPortStr, 10, 16)

	d.sshPort = uint16(sshPort64)

	return nil

}

// GetValidationErrors returns the configuration validation errors
func (d Docker) GetValidationErrors() []error {
	return nil
}
