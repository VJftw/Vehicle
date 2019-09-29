package ssh

import "time"

// Config represents the SSHConfig for an instance
type Config struct {
	User    string        `json:"user"`
	Port    uint16        `json:"port"`
	Timeout time.Duration `json:"timeout"`
}
