package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/helpers"

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

		supervysorDir, err := helpers.GetSupervysorDir()
		if err != nil {
			logger.Error("could not load supervysor directory")
			return err
		}

		logFile, err := os.OpenFile(supervysorDir+"/logs/"+time.Now().Format("20060102_150405")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			logger.Error(fmt.Sprintf("could not open log file: %v", err))
			return err
		}
		defer logFile.Close()

		fileLogger := log.New(logFile, "", log.LstdFlags)

		// Start data source node initially.
		if err := node.InitialStart(config.BinaryPath, config.AddrBookPath, config.Seeds); err != nil {
			fileLogger.Printf("initial start failed: %s", err)
			logger.Error("initial start failed", "err", err)
			return err
		}

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

			fileLogger.Printf(fmt.Sprintf("fetched heights successfully; node=%d pool=%d current-target-height=%d", nodeHeight, poolHeight, *poolHeight+config.HeightDifferenceMax))
			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight, "current-target-height", *poolHeight+config.HeightDifferenceMax)

			// Calculate height difference to enable the correct mode.
			heightDiff := nodeHeight - *poolHeight

			if heightDiff >= config.HeightDifferenceMax {
				// Data source node has synced far enough, enable or keep Ghost Mode
				if err = node.EnableGhostMode(config.BinaryPath, config.AddrBookPath); err != nil {
					fileLogger.Printf(fmt.Sprintf("could not enable Ghost Mode err=%s", err))
					logger.Error("could not enable Ghost Mode", "err", err)

					if shutdownErr := node.ShutdownNode(); shutdownErr != nil {
						logger.Error("could not shutdown node process", "err", shutdownErr)
					}
					return err
				}
			} else if heightDiff < config.HeightDifferenceMax && heightDiff > config.HeightDifferenceMin {
				// No threshold reached, keep current mode
				fileLogger.Printf(fmt.Sprintf("keeping current Mode: height-difference=%d", heightDiff))
				logger.Info("keeping current Mode", "height-difference", heightDiff)
			} else {
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
					fileLogger.Printf("node has not reached pool height yet, can not use it as data source")
					logger.Info("node has not reached pool height yet, can not use it as data source")
				}
			}
			time.Sleep(time.Second * time.Duration(config.Interval))
		}
	},
}
