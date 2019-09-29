package ssh

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnectToSSH(address string, port uint16, timeout time.Duration, sshConfig *ssh.ClientConfig) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), sshConfig)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	defer session.Close()
	// Set up terminal modes
	// modes := ssh.TerminalModes{
	// 	ssh.ECHO:          0,     // disable echoing
	// 	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	// 	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	// }
	// Request pseudo terminal
	// if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
	// 	log.Fatal("request for pseudo terminal failed: ", err)
	// }

	stdout, err := session.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("Unable to setup stdout for session: %v", err))
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		panic(fmt.Errorf("Unable to setup stderr for session: %v", err))
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("whoami")

	// create new session for each command
}
