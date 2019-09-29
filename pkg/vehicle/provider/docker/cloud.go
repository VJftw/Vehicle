package docker

import (
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle/provider"
	"golang.org/x/crypto/ssh"
)

// Docker represents the Docker vehicle instance provider
type Docker struct {
	*provider.Base

	containerID string
}

// New returns a new Docker vehicle instance
func New(uuid string) *Docker {
	return &Docker{
		Base: &provider.Base{Provider: "docker", UUID: uuid},
	}
}

// GetSSHInfo returns the ssh configuration for a provisioned Docker instance
func (d *Docker) GetSSHInfo() (string, uint16, time.Duration, *ssh.ClientConfig) {
	return "127.0.0.1", uint16(2222), time.Minute, &ssh.ClientConfig{}
}

// ResolveFuncs returns the functions involved in resolving dynamic resource IDs
func (d *Docker) ResolveFuncs() []func() (error, []string) {
	return []func() (error, []string){}
}

// StartFuncs returns the functions involved in starting an AWS instance
func (d *Docker) StartFuncs() []func() error {
	return []func() error{
		// d.createPrivateKey,
		// d.generateSSHConfig,
		// d.launchInstance,
		// d.waitForInstance,
	}
}

// Stop stops and all AWS resources created by vehicle
func (d *Docker) Stop() error {
	// if d.containerID != nil {

	// }
	// if len(d.privateKey) > 0 {

	// }
	return nil
}
