package main

import (
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
		// Load initialized config.
		config, err := getConfig()
		if err != nil {
			logger.Error("could not load config", "err", err)
			return err
		}

		// Start data source node initially.
		if err := node.InitialStart(config.BinaryPath, config.AddrBookPath, config.Seeds); err != nil {
			logger.Error("initial start failed", "err", err)
			return err
		}

		var currentMode = ""

		for {
			// Request data source node height and KYVE pool height to calculate difference.
			nodeHeight, err := node.GetNodeHeight(0)
			if err != nil {
				logger.Error("could not get node height", "err", err)
				if shutdownErr := node.ShutdownNode(); shutdownErr != nil {
					logger.Error("could not shutdown node process", "err", shutdownErr)
				}
				return err
			}

			poolHeight, err := pool.GetPoolHeight(config.ChainId, config.PoolId, config.FallbackEndpoints)
			if err != nil {
				logger.Error("could not get pool height", "err", err)
				if shutdownErr := node.ShutdownNode(); shutdownErr != nil {
					logger.Error("could not shutdown node process", "err", shutdownErr)
				}
				return err
			}

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight, "max-height", *poolHeight+config.HeightDifferenceMax, "min-height", *poolHeight+config.HeightDifferenceMin)

			// Calculate height difference to enable the correct mode.
			heightDiff := nodeHeight - *poolHeight

			if heightDiff >= config.HeightDifferenceMax {
				if currentMode != "ghost" {
					logger.Info("enabling GhostMode")
					currentMode = "ghost"
				}
				// Data source node has synced far enough, enable or keep Ghost Mode
				if err = node.EnableGhostMode(config.BinaryPath, config.AddrBookPath); err != nil {
					logger.Error("could not enable Ghost Mode", "err", err)

					if shutdownErr := node.ShutdownNode(); shutdownErr != nil {
						logger.Error("could not shutdown node process", "err", shutdownErr)
					}
					return err
				}
			} else if heightDiff < config.HeightDifferenceMax && heightDiff > config.HeightDifferenceMin {
				// No threshold reached, keep current mode
				logger.Info("keeping current Mode", "mode", currentMode, "height-difference", heightDiff)
			} else {
				if currentMode != "normal" {
					logger.Info("enabling NormalMode")
					currentMode = "normal"
				}
				// Difference is < HeightDifferenceMin, Data source needs to catch up, enable or keep Normal Mode
				if err = node.EnableNormalMode(config.BinaryPath, config.AddrBookPath, config.Seeds); err != nil {
					logger.Error("could not enable Normal Mode", "err", err)

					if shutdownErr := node.ShutdownNode(); shutdownErr != nil {
						logger.Error("could not shutdown node process", "err", shutdownErr)
					}
					return err
				}
				// Diff < 0, can't use node as data source
				if heightDiff <= 0 {
					logger.Info("node has not reached pool height yet, can not use it as data source")
				}
			}
			time.Sleep(time.Second * time.Duration(config.Interval))
		}
	},
}
