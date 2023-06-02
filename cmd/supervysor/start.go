package main

import (
	"time"

	"github.com/KYVENetwork/supervysor/pool"

	"github.com/KYVENetwork/supervysor/node"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [binary-path] [pool-id] [seeds]",
	Short: "Start a supervysed Tendermint node.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfig()
		if err != nil {
			logger.Error("could not load config", "err", err)
			return err
		}

		if _, err := node.InitialStart(config.BinaryPath, config.Seeds); err != nil {
			return err
		}

		for {
			nodeHeight := node.GetNodeHeight()
			poolHeight, err := pool.GetPoolHeight(config.ChainId, config.PoolId)
			if err != nil {
				logger.Error("couldn't get pool height")
				return err
			}

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight)

			diff := nodeHeight - *poolHeight

			if diff >= config.HeightDifferenceMax {
				node.EnableGhostMode(config.BinaryPath)
			} else if diff < config.HeightDifferenceMax && diff > config.HeightDifferenceMin {
				// do nothing
			} else if diff <= config.HeightDifferenceMin && diff > 0 {
				node.DisableGhostMode(config.BinaryPath, config.Seeds)
			} else {
				logger.Error("negative difference between node and pool heights")
			}

			time.Sleep(time.Second * time.Duration(config.Interval))
		}
	},
}
