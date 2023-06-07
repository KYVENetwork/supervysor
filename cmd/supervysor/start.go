package main

import (
	"fmt"
	"time"

	"github.com/KYVENetwork/supervysor/pool"

	"github.com/KYVENetwork/supervysor/node"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a supervysed Tendermint node",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfig()
		if err != nil {
			logger.Error("could not load config", "err", err)
			return err
		}

		if _, err := node.InitialStart(config.BinaryPath, config.AddrBookPath, config.Seeds); err != nil {
			logger.Error("initial start failed", "err", err)
			return err
		}

		for {
			nodeHeight := node.GetNodeHeight()
			poolHeight, err := pool.GetPoolHeight(config.ChainId, config.PoolId)
			if err != nil {
				logger.Error("couldn't get pool height", "err", err)
				return err
			}

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight)

			diff := nodeHeight - *poolHeight

			if diff >= config.HeightDifferenceMax {
				if err = node.EnableGhostMode(config.BinaryPath, config.AddrBookPath); err != nil {
					logger.Error("could not enable Ghost Mode", "err", err)
					return err
				}
			} else if diff < config.HeightDifferenceMax && diff > config.HeightDifferenceMin {
				logger.Info("keeping current Mode", "height-difference", diff)
			} else if diff <= config.HeightDifferenceMin && diff > 0 {
				if err = node.DisableGhostMode(config.BinaryPath, config.AddrBookPath, config.Seeds); err != nil {
					logger.Error("could not disable Ghost Mode", "err", err)
					return err
				}
			} else {
				if err = node.ShutodwnProcess(); err != nil {
					logger.Error("could not shutdown process", "err", err)
				}
				return fmt.Errorf("negative difference between node and pool heights")
			}
			time.Sleep(time.Second * time.Duration(config.Interval))
		}
	},
}
