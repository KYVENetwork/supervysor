package main

import (
	"time"

	"github.com/KYVENetwork/supervysor/node"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a supervysed Tendermint node.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := cast.ToUint64E(args[0])
		if err != nil {
			return err
		}

		if _, err := node.InitialStart(); err != nil {
			return err
		}

		for {
			nodeHeight := node.GetNodeHeight()
			poolHeight := 1 // TODO(@christopher): Replace with real height.

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight)

			diff := nodeHeight - poolHeight

			if diff >= 1000 {
				node.EnableGhostMode()
			} else if diff < 1000 && diff > 500 {
				// do nothing
			} else if diff <= 500 && diff > 0 {
				node.DisableGhostMode()
			} else {
				logger.Error("negative difference between node and pool heights")
			}

			time.Sleep(time.Minute / 6)
		}
	},
}
