package vehicle

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle/provider/aws"
	"github.com/VJftw/vehicle/pkg/vehicle/provider/docker"
)

// UUIDFunc is an overridable function to return a UUID
var UUIDFunc = func() string {
	return time.Now().Format(fmt.Sprintf("vehicle-20060102150405"))
}

// Config represents a vehicle config
type Config struct {
	Clouds           map[string]Vehicle `json:"clouds"`
	WorkingDirectory string             `json:"workingDirectory"`
	Files            []string           `json:"files"`
	Commands         []string           `json:"commands"`
}

// NewConfig returns a new configuration
func NewConfig() *Config {
	return &Config{
		Clouds:           map[string]Vehicle{},
		WorkingDirectory: "",
		Files:            []string{},
		Commands:         []string{},
	}
}

// UnmarshalJSON provides custom JSON unmarshalling
func (c *Config) UnmarshalJSON(b []byte) error {
	// We don't return any errors from this function so we can show more helpful parse errors
	var objMap map[string]*json.RawMessage
	// We'll store the error (if any) so we can return it if necessary
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		// c = handleBlueprintUnmarshalError(t, err)
		return err
	}

	// Unmarshal Commands
	if _, ok := objMap["commands"]; ok {
		err = json.Unmarshal(*objMap["commands"], &c.Commands)
		// c = handleBlueprintUnmarshalError(t, err)
	}

	// Unmarshal clouds by provider
	c.Clouds = map[string]Vehicle{}
	if v, _ := objMap["clouds"]; v != nil {
		var rawClouds map[string]*json.RawMessage
		err = json.Unmarshal(*v, &rawClouds)
		// c = handleUnmarshalError(c, err)
		if err == nil {
			for id, rawMessage := range rawClouds {
				config, err := unmarshalCloud(id, *rawMessage)
				// c = handleUnmarshalError(c, err)
				if err == nil {
					c.Clouds[id] = config
				}
			}
		}
	}

	return nil
}

func unmarshalCloud(id string, rawMessage []byte) (Vehicle, error) {
	var m map[string]interface{}
	err := json.Unmarshal(rawMessage, &m)
	if err != nil {
		return nil, err
	}

	var c Vehicle
	switch m["provider"] {
	case "aws":
		c = aws.New(fmt.Sprintf("%s-%s", id, UUIDFunc()))
	case "docker":
		c = docker.New(fmt.Sprintf("%s-%s", id, UUIDFunc()))
	default:
		return nil, fmt.Errorf("could not determine provider: %+v", m)
	}

	err = json.Unmarshal(rawMessage, c)
	return c, err
}
