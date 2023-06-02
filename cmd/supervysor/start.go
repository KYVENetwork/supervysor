package main

import (
	"time"

	"github.com/KYVENetwork/supervysor/pool"

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

		poolId, err := cast.ToIntE(args[1])
		if err != nil {
			return err
		}

		argSeeds, err := cast.ToStringE(args[2])
		if err != nil {
			return err
		}

		if _, err := node.InitialStart(argBinaryPath, argSeeds); err != nil {
			return err
		}

		config, err := getConfig()
		if err != nil {
			logger.Error("couldn't load config")
			return err
		}

		chainId := config.ChainId
		interval := config.Interval
		heightDifferenceMax := config.HeightDifferenceMax
		heightDifferenceMin := config.HeightDifferenceMin

		for {
			nodeHeight := node.GetNodeHeight()
			poolHeight, err := pool.GetPoolHeight(chainId, poolId)
			if err != nil {
				logger.Error("couldn't get pool height")
				return err
			}

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight)

			diff := nodeHeight - *poolHeight

			if diff >= heightDifferenceMax {
				node.EnableGhostMode(argBinaryPath)
			} else if diff < heightDifferenceMax && diff > heightDifferenceMin {
				// do nothing
			} else if diff <= heightDifferenceMin && diff > 0 {
				node.DisableGhostMode(argBinaryPath, argSeeds)
			} else {
				logger.Error("negative difference between node and pool heights")
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	},
}
