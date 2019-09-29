package cmds

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// BuildVersion represents the current build tag of vehicle. It is set at compile-time with ldflags
var BuildVersion = "dev"

var (
	noColor bool

	gracefulStop = make(chan os.Signal)
)

var rootCmd = &cobra.Command{
	Use:   "vehicle",
	Short: "Vehicle",
	Long:  `A tool to idempotently run tasks on cloud infrastructure utilising instances`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor {
			// output.ColorDisable()
		}
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)
		go func() {
			sig := <-gracefulStop
			fmt.Printf("\ncaught signal: %+v\nstopping...\n", sig)
			// action.Stop()
		}()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Version: BuildVersion,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
}

// Execute runs the CLI
func Execute() error {
	return rootCmd.Execute()
}
