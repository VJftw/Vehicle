package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/VJftw/vehicle/pkg/vehicle"
)

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
		address, port, timeout, _ := c.vehicle.GetSSHInfo()
		vehicle.WaitForSSH(address, port, timeout)
	}

	// Run commands

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
