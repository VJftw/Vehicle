package vehicle

import (
	"fmt"
	"time"
)

type MountConfig struct {
	Size string `json:"size" yaml:"size"`
	Path string `json:"path" yaml:"path"`
}

func (m *MountConfig) UnmarshalYAMLConfig(configMap interface{}) error {
	if y, ok := configMap.(map[interface{}]interface{}); ok {
		switch x := y["size"].(type) {
		case string:
			m.Size = x
			break
		}
		switch x := y["path"].(type) {
		case string:
			m.Path = x
			break
		}
		return nil
	}

	return fmt.Errorf("unsupported map config")
}

type SubnetConfig struct {
	ID        string            `json:"id" yaml:"-"`
	Tags      map[string]string `json:"tags" yaml:"tags"`
	CIDRBlock string            `json:"cidrBlock" yaml:"cidr_block"`
}

type SecurityGroupsConfig struct {
	IDs  []string          `json:"ids" yaml:"-"`
	Tags map[string]string `json:"tags" yaml:"tags"`
}

type AMIConfig struct {
	ID     string `json:"id" yaml:"-"`
	Filter string `json:"filter" yaml:"filter"`
}

type SSHConfig struct {
	User    string        `json:"user" yaml:"user"`
	Port    uint16        `json:"port" yaml:"port"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func (s *SSHConfig) UnmarshalYAMLConfig(configMap interface{}) error {
	if y, ok := configMap.(map[interface{}]interface{}); ok {
		switch x := y["user"].(type) {
		case string:
			s.User = x
			break
		}
		switch x := y["port"].(type) {
		case int:
			s.Port = uint16(x)
			break
		default:
			s.Port = 22
		}
		switch x := y["timeout"].(type) {
		case int:
			s.Timeout = time.Duration(x) * time.Second
			break
		default:
			s.Timeout = 300 * time.Second
		}
		return nil
	}

	return fmt.Errorf("unsupported map config")
}

type AWSConfig struct {
	Type           string                 `json:"type" yaml:"type"`
	Mounts         []MountConfig          `json:"mounts" yaml:"mounts"`
	Subnet         SubnetConfig           `json:"subnet" yaml:"subnet"`
	SecurityGroups []SecurityGroupsConfig `json:"securityGroups" yaml:"security_groups"`
	IAMPolicy      string                 `json:"iamPolicy" yaml:"iam_policy"`
	AMI            AMIConfig              `json:"ami" yaml:"ami"`
	SSH            SSHConfig              `json:"ssh" yaml:"ssh"`
}

func (i *AWSConfig) UnmarshalYAMLConfig(configMap interface{}) error {
	if m, ok := configMap.(map[interface{}]interface{}); ok {
		switch x := m["type"].(type) {
		case string:
			i.Type = x
			break
		}

		i.Subnet = SubnetConfig{}
		switch x := m["subnet"].(type) {
		case string:
			i.Subnet.ID = x
			break
		case interface{}:
			break
		}

		i.SecurityGroups = []SecurityGroupsConfig{}
		switch x := m["security_groups"].(type) {
		case []interface{}:
			for _, sg := range x {
				switch y := sg.(type) {
				case string:
					i.SecurityGroups = append(i.SecurityGroups, SecurityGroupsConfig{
						IDs: []string{y},
					})
					break
				}
			}
			break
		}

		i.Mounts = []MountConfig{}
		switch x := m["mounts"].(type) {
		case []interface{}:
			for _, mount := range x {
				switch y := mount.(type) {
				case interface{}:
					j := MountConfig{}
					j.UnmarshalYAMLConfig(y)
					i.Mounts = append(i.Mounts, j)
					break
				}
			}
			break
		}

		switch x := m["iam_policy"].(type) {
		case string:
			i.IAMPolicy = x
			break
		}

		i.AMI = AMIConfig{}
		switch x := m["ami"].(type) {
		case string:
			i.AMI.ID = x
			break
		case interface{}:
			// i.AMI.UnmarshalYAML(x)
			break
		default:
			fmt.Println("invalid type for ami")
		}

		switch x := m["ssh"].(type) {
		case map[interface{}]interface{}:
			i.SSH = SSHConfig{}
			i.SSH.UnmarshalYAMLConfig(x)
			break
		default:
			fmt.Println("invalid type for ssh")
		}

		return nil
	}
	return fmt.Errorf("incompatbile interface")
}

type Config struct {
	AWSConfig        AWSConfig `json:"aws" yaml:"aws"`
	WorkingDirectory string    `json:"workingDirectory" yaml:"working_directory"`
	Files            []string  `json:"files" yaml:"files"`
	Commands         []string  `json:"commands" yaml:"commands"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var configMap map[string]interface{}
	err := unmarshal(&configMap)
	if err != nil {
		return err
	}

	switch x := configMap["working_directory"].(type) {
	case string:
		c.WorkingDirectory = x
		break
	}

	c.Files = []string{}
	switch x := configMap["files"].(type) {
	case []interface{}:
		for _, f := range x {
			if val, ok := f.(string); ok {
				c.Files = append(c.Files, val)
			}
		}
		break
	}

	c.Commands = []string{}
	switch x := configMap["commands"].(type) {
	case []interface{}:
		for _, command := range x {
			if val, ok := command.(string); ok {
				c.Commands = append(c.Commands, val)
			}
		}
		break
	}

	switch x := configMap["aws"].(type) {
	case interface{}:
		c.AWSConfig = AWSConfig{}
		c.AWSConfig.UnmarshalYAMLConfig(x)
	}

	return nil
}
