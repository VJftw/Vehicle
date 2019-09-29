package cmds

import (
	"fmt"

	"github.com/VJftw/vehicle/pkg/vehicle"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "inits a vehicle configuration",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// path := getAbsPath(args[0])

		newConfig := vehicle.NewConfig()
		y, err := yaml.Marshal(newConfig)
		if err != nil {
			return err
		}

		fmt.Println(string(y))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
