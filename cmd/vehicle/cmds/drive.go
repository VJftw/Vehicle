package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle"
	"github.com/VJftw/vehicle/pkg/vehicle/ssh"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var driveCmd = &cobra.Command{
	Use:     "drive",
	Aliases: []string{"d"},
	Short:   "Drives a vehicle configuration",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := getAbsPath(args[0])
		fBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var config vehicle.Config
		err = yaml.Unmarshal(fBytes, &config)
		if err != nil {
			return err
		}

		// uuid := generateUUID()

		clouds := []*cloud{}
		for _, v := range config.Clouds {
			clouds = append(clouds, NewCloud(v))
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		var wg sync.WaitGroup

		for _, c := range clouds {
			wg.Add(1)
			go c.Run(quit, &wg)
		}
		<-quit
		for _, c := range clouds {
			c.Stop()
		}
		wg.Wait()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(driveCmd)
}

func generateUUID() string {
	return time.Now().Format(fmt.Sprintf("vehicle-20060102150405"))
}

type cloud struct {
	vehicle vehicle.Vehicle

	stop chan os.Signal
	run  bool
}

func NewCloud(v vehicle.Vehicle) *cloud {
	return &cloud{vehicle: v, run: true}
}

func (c *cloud) Run(stop chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { c.run = false }()
	defer c.cleanUp(stop)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// Start instance
	for _, startFunc := range c.vehicle.StartFuncs() {
		if !c.run {
			break
		}
		if err := startFunc(); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
	}

	// Wait for SSH
	if c.run {
		address, port, timeout, sshConfig := c.vehicle.GetSSHInfo()
		vehicle.WaitForSSH(address, port, timeout)
		// Run commands
		ssh.ConnectToSSH(address, port, timeout, sshConfig)
	}

}

func (c *cloud) Stop() {
	fmt.Println("stopping")
	c.run = false
}

func (c *cloud) cleanUp(stop chan os.Signal) {
	// Clean up
	if err := c.vehicle.Stop(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	if c.run {
		stop <- os.Interrupt
	}
}

func getAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}
	return path
}
