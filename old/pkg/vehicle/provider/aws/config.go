package aws

import (
	"github.com/VJftw/vehicle/pkg/vehicle/ssh"
)

// Config represents the configuration for an AWS instance
type Config struct {
	Type             string                `json:"type"`
	Mounts           []MountConfig         `json:"mounts"`
	Subnet           SubnetConfig          `json:"subnet"`
	SecurityGroups   []SecurityGroupConfig `json:"securityGroups"`
	IAMPolicy        string                `json:"iamPolicy"`
	AMI              AMIConfig             `json:"ami"`
	SSH              ssh.Config            `json:"ssh"`
	ValidationErrors []error               `json:"validationErrors"`
}

// NewConfig returns a new AWS configuration w/ defaults
func NewConfig() *Config {
	return &Config{
		Type:             "t3.nano",
		Mounts:           []MountConfig{},
		Subnet:           SubnetConfig{},
		SecurityGroups:   []SecurityGroupConfig{},
		IAMPolicy:        "",
		AMI:              AMIConfig{},
		SSH:              ssh.Config{Port: 22, Timeout: 300, User: "aws"},
		ValidationErrors: []error{},
	}
}

// GetValidationErrors returns the configuration validation errors
func (c Config) GetValidationErrors() []error {
	return nil
}

// SubnetConfig represents the subnet configuration for an AWS instance
type SubnetConfig struct {
	ID        string            `json:"id"`
	Tags      map[string]string `json:"tags"`
	CIDRBlock string            `json:"cidrBlock"`
}

// SecurityGroupConfig represents the security group configuration for an AWS instance
type SecurityGroupConfig struct {
	IDs  []string          `json:"ids"`
	Tags map[string]string `json:"tags"`
}

// AMIConfig represents the AMI configuration for an AWS instance
type AMIConfig struct {
	ID     string `json:"id"`
	Filter string `json:"filter"`
}

// MountConfig represents the volume mount configuration for an AWS instance
type MountConfig struct {
	Size string `json:"size"`
	Path string `json:"path"`
}
