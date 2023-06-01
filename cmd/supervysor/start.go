package main

import (
	"time"

	"github.com/KYVENetwork/supervysor/settings"

	"github.com/KYVENetwork/supervysor/node"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [binary-path] [pool-id] [seeds]",
	Short: "Start a supervysed Tendermint node.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		argBinaryPath, err := cast.ToStringE(args[0])
		if err != nil {
			return err
		}

		poolId, err := cast.ToUint64E(args[1])
		if err != nil {
			return err
		}

		argSeeds, err := cast.ToStringE(args[2])
		if err != nil {
			return err
		}

		if err := settings.InitializeSettings(argBinaryPath, int64(poolId)); err != nil {
			return err
		}

		if _, err := node.InitialStart(argBinaryPath, argSeeds); err != nil {
			return err
		}

		for {
			nodeHeight := node.GetNodeHeight()
			poolHeight := 1 // TODO(@christopher): Replace with real height.

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight)

			diff := nodeHeight - poolHeight

			if diff >= 1000 {
				node.EnableGhostMode(argBinaryPath)
			} else if diff < 1000 && diff > 500 {
				// do nothing
			} else if diff <= 500 && diff > 0 {
				node.DisableGhostMode(argBinaryPath, argSeeds)
			} else {
				logger.Error("negative difference between node and pool heights")
			}

			time.Sleep(time.Minute / 6)
		}
	},
}
