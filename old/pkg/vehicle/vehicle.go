package vehicle

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Vehicle - Cloud providers should implement this
type Vehicle interface {
	GetProvider() string
	GetValidationErrors() []error

	// Resolve returns error and error strings if invalid
	ResolveFuncs() []func() (error, []string)
	StartFuncs() []func() error
	GetSSHInfo() (string, uint16, time.Duration, *ssh.ClientConfig)
	Stop() error
}

// add general purpose running commands over SSH func in here
func WaitForSSH(ip string, port uint16, timeout time.Duration) error {
	start := time.Now()
	for start.Add(timeout).After(time.Now()) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err == nil {
			defer conn.Close()
			return nil
		}
		fmt.Printf("%s not yet up (%2f)\n", ip, time.Now().Sub(start).Seconds())
		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("%s:%d not responsive after %2f seconds", ip, port, time.Now().Sub(start).Seconds())
}
