package main

import (
	"os"

	"cosmossdk.io/log"

	"github.com/spf13/cobra"
)

var logger = log.NewLogger(os.Stdout)

var Version = ""

var supervysor = &cobra.Command{
	Use:     "supervysor",
	Short:   "Supervysor helps sync a Tendermint node used as a KYVE data source.",
	Version: Version,
}

func main() {
	supervysor.AddCommand(initCmd)
	supervysor.AddCommand(startCmd)
	supervysor.AddCommand(versionCmd)

	if err := supervysor.Execute(); err != nil {
		os.Exit(1)
	}
}
