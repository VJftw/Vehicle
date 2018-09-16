package vehicle

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Vehicle - Cloud providers should implement this
type Vehicle interface {
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

// func connectToSSH() {
// 	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", publicIP), sshConfig)
// 	if err != nil {
// 		log.Fatalf("unable to connect: %v", err)
// 	}

// 	defer client.Close()

// 	// Create a session
// 	session, err := client.NewSession()
// 	if err != nil {
// 		log.Fatal("unable to create session: ", err)
// 	}
// 	defer session.Close()
// 	// Set up terminal modes
// 	// modes := ssh.TerminalModes{
// 	// 	ssh.ECHO:          0,     // disable echoing
// 	// 	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
// 	// 	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
// 	// }
// 	// Request pseudo terminal
// 	// if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
// 	// 	log.Fatal("request for pseudo terminal failed: ", err)
// 	// }

// 	stdout, err := session.StdoutPipe()
// 	if err != nil {
// 		panic(fmt.Errorf("Unable to setup stdout for session: %v", err))
// 	}
// 	go io.Copy(os.Stdout, stdout)

// 	stderr, err := session.StderrPipe()
// 	if err != nil {
// 		panic(fmt.Errorf("Unable to setup stderr for session: %v", err))
// 	}
// 	go io.Copy(os.Stderr, stderr)

// 	err = session.Run("whoami")

// 	// create new session for each command
// }
