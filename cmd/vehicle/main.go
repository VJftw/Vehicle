package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "vehicle"
	app.Usage = "Runs commands on cloud infrastructure"
	// app.Action = func(c *cli.Context) error {
	// 	fmt.Println("Hello friend!")
	// 	return nil
	// }

	app.Commands = []cli.Command{
		{
			Name:    "drive",
			Aliases: []string{"d"},
			Usage:   "Drives a Vehicle configuration",
			Action:  drive,
		},
		{
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "Validates a Vehicle configuration",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateUUID() string {
	return time.Now().Format(fmt.Sprintf("vehicle-20060102150405"))
}

func drive(c *cli.Context) error {
	fBytes, err := ioutil.ReadFile("vehicle.yml")
	if err != nil {
		return err
	}

	var config vehicle.Config
	err = yaml.Unmarshal(fBytes, &config)
	if err != nil {
		return err
	}

	uuid := generateUUID()

	cloud := NewCloud(vehicle.NewAWS(&config.AWSConfig, uuid))

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	var wg sync.WaitGroup

	wg.Add(1)
	go cloud.Run(quit, &wg)
	<-quit
	cloud.Stop()
	wg.Wait()

	return nil
}
