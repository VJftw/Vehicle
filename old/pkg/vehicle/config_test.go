package vehicle_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestConfigUnmarshal(t *testing.T) {
	fBytes, err := ioutil.ReadFile("config_test.yml")
	if err != nil {
		t.Error(err)
	}

	var config vehicle.Config
	err = yaml.Unmarshal(fBytes, &config)
	assert.Nil(t, err)

	expectedConfig := vehicle.Config{
		AWSConfig: vehicle.AWSConfig{
			Type: "t3.nano",
			Mounts: []vehicle.MountConfig{
				vehicle.MountConfig{
					Size: "100G",
					Path: "/var/mount",
				},
			},
			Subnet: vehicle.SubnetConfig{
				ID: "subnet-xxxxxxxx",
			},
			SecurityGroups: []vehicle.SecurityGroupsConfig{
				vehicle.SecurityGroupsConfig{
					IDs: []string{"sg-xxxxxxxx"},
				},
			},
			IAMPolicy: "\\{\n\\}\n",
			AMI: vehicle.AMIConfig{
				ID: "ami-xxxxxxxx",
			},
			SSH: vehicle.SSHConfig{
				User:    "ubuntu",
				Port:    22,
				Timeout: 300 * time.Second,
			},
		},
		WorkingDirectory: "/app",
		Files: []string{
			"./aaa.txt",
		},
		Commands: []string{
			"ansible-playbook",
		},
	}

	assert.Equal(t, expectedConfig, config)

}
