package main

import (
	"fmt"
	"os"

	"github.com/VJftw/vehicle/cmd/vehicle/cmds"
)

func main() {
	if err := cmds.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

// func generateUUID() string {
// 	return time.Now().Format(fmt.Sprintf("vehicle-20060102150405"))
// }

// func drive(c *cli.Context) error {
// 	fBytes, err := ioutil.ReadFile("vehicle.yml")
// 	if err != nil {
// 		return err
// 	}

// 	var config vehicle.Config
// 	err = yaml.Unmarshal(fBytes, &config)
// 	if err != nil {
// 		return err
// 	}

// 	uuid := generateUUID()

// 	cloud := NewCloud(vehicle.NewAWS(&config.AWSConfig, uuid))

// 	quit := make(chan os.Signal)
// 	signal.Notify(quit, os.Interrupt)
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go cloud.Run(quit, &wg)
// 	<-quit
// 	cloud.Stop()
// 	wg.Wait()

// 	return nil
// }
