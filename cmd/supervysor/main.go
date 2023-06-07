package main

import (
	"github.com/KYVENetwork/supervysor/node"
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
		if err := node.ShutdownNode(); err != nil {
			logger.Info("could not shutdown node process", "err", err)
		}
		os.Exit(1)
	}
}
