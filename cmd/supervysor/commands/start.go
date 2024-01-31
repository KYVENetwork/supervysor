package commands

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"

	"github.com/KYVENetwork/supervysor/utils"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/spf13/cobra"

	"github.com/KYVENetwork/supervysor/executor"
	"github.com/KYVENetwork/supervysor/pool"
)

func init() {
	startCmd.Flags().StringVar(&cfgFlag, "config", "", "path to config directory (e.g. ~/.supervysor/)")

	startCmd.Flags().BoolVar(&statePruning, "state-pruning", true, "enable state pruning")

	startCmd.Flags().StringVar(&binaryFlags, "flags", "", "flags for the underlying binary (e.g. '--address, ')")

	startCmd.Flags().BoolVar(&optOut, "opt-out", false, "disable the collection of anonymous usage data")
}

// The startCmd of the supervysor launches and manages the node process using the specified binary.
// It periodically retrieves the heights of the node and the associated KYVE pool, and dynamically adjusts
// the sync mode of the node based on these heights.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a supervysed Tendermint node",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgFlag == "" {
			configPath, err = helpers.GetSupervysorDir()
			if err != nil {
				logger.Error("could not get supervysor directory path", "err", err)
				return err
			}
		} else {
			configPath = cfgFlag
		}

		// Load initialized config.
		supervysorConfig, err := getSupervysorConfig(configPath)
		if err != nil {
			logger.Error("could not load config", "err", err)
			return err
		}

		utils.TrackStartEvent(supervysorConfig.ChainId, optOut)

		metrics := supervysorConfig.Metrics

		// Create Prometheus registry
		reg := prometheus.NewRegistry()
		m := helpers.NewMetrics(reg)

		if metrics {
			go func() {
				err := helpers.StartMetricsServer(reg, supervysorConfig.MetricsPort)
				if err != nil {
					panic(err)
				}
			}()
		}

		e := executor.NewExecutor(&logger, supervysorConfig)

		// Start data source node initially.
		if err := e.InitialStart(binaryFlags); err != nil {
			logger.Error("initial start failed", "err", err)
			return err
		}

		currentMode := "normal"

		if metrics {
			go func() {
				for {
					dbSize, err := helpers.GetDirectorySize(filepath.Join(supervysorConfig.HomePath, "data"))
					if err != nil {
						logger.Error("could not get data directory size; will not expose metrics", "err", err)
					} else {
						m.DataDirSize.Set(dbSize)
					}

					time.Sleep(time.Second * time.Duration(120))
				}
			}()
		}

		var pruningCount float64 = 0
		for {
			// Request data source node height and KYVE pool height to calculate difference.
			nodeHeight, err := e.GetHeight()
			if err != nil {
				logger.Error("could not get node height", "err", err)
				if shutdownErr := e.Shutdown(); shutdownErr != nil {
					logger.Error("could not shutdown node process", "err", shutdownErr)
				}
				return err
			}
			if metrics {
				m.NodeHeight.Set(float64(nodeHeight))
			}

			poolHeight, err := pool.GetPoolHeight(supervysorConfig.ChainId, supervysorConfig.PoolId, supervysorConfig.PoolEndpoints)
			if err != nil {
				logger.Error("could not get pool height", "err", err)
				if shutdownErr := e.Shutdown(); shutdownErr != nil {
					logger.Error("could not shutdown node process", "err", shutdownErr)
				}
				return err
			}
			if metrics {
				m.PoolHeight.Set(float64(poolHeight))
			}

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight, "max-height", poolHeight+supervysorConfig.HeightDifferenceMax, "min-height", poolHeight+supervysorConfig.HeightDifferenceMin)

			if supervysorConfig.PruningInterval != 0 {
				logger.Info("current pruning count", "pruning-count", fmt.Sprintf("%.2f", pruningCount), "pruning-threshold", supervysorConfig.PruningInterval)
				if pruningCount > float64(supervysorConfig.PruningInterval) && nodeHeight > 0 {
					if currentMode == "ghost" {
						pruneHeight := poolHeight
						if nodeHeight < poolHeight {
							pruneHeight = nodeHeight
						}
						logger.Info("pruning after node shutdown", "until-height", pruneHeight)

						err = e.PruneData(supervysorConfig.HomePath, pruneHeight-1, supervysorConfig.StatePruning, binaryFlags)
						if err != nil {
							logger.Error("could not prune", "err", err)
							return err
						}
					} else {
						if nodeHeight < poolHeight {
							logger.Info("pruning after node shutdown", "until-height", nodeHeight)

							err = e.PruneData(supervysorConfig.HomePath, nodeHeight-1, supervysorConfig.StatePruning, binaryFlags)
							if err != nil {
								logger.Error("could not prune", "err", err)
								return err
							}
						}
					}
					pruningCount = 0
				}
			}

			// Calculate height difference to enable the correct mode.
			heightDiff := nodeHeight - poolHeight

			if metrics {
				m.MaxHeight.Set(float64(poolHeight + supervysorConfig.HeightDifferenceMax))
				m.MinHeight.Set(float64(poolHeight + supervysorConfig.HeightDifferenceMin))
			}

			if heightDiff >= supervysorConfig.HeightDifferenceMax {
				if currentMode != "ghost" {
					logger.Info("enabling GhostMode")
				} else {
					logger.Info("keeping GhostMode")
				}
				// Data source node has synced far enough, enable or keep Ghost Mode
				if err = e.EnableGhostMode(binaryFlags); err != nil {
					logger.Error("could not enable Ghost Mode", "err", err)

					if shutdownErr := e.Shutdown(); shutdownErr != nil {
						logger.Error("could not shutdown node process", "err", shutdownErr)
					}
					return err
				}
				currentMode = "ghost"
			} else if heightDiff < supervysorConfig.HeightDifferenceMax && heightDiff > supervysorConfig.HeightDifferenceMin {
				// No threshold reached, keep current mode
				logger.Info("keeping current Mode", "mode", currentMode, "height-difference", heightDiff)
			} else {
				if currentMode != "normal" {
					logger.Info("enabling NormalMode")
				} else {
					logger.Info("keeping NormalMode")
				}
				// Difference is < HeightDifferenceMin, Data source needs to catch up, enable or keep Normal Mode
				if err = e.EnableNormalMode(binaryFlags); err != nil {
					logger.Error("could not enable Normal Mode", "err", err)

					if shutdownErr := e.Shutdown(); shutdownErr != nil {
						logger.Error("could not shutdown node process", "err", shutdownErr)
					}
					return err
				}
				currentMode = "normal"

				// Diff < 0, can't use node as data source
				if heightDiff <= 0 {
					logger.Info("node has not reached pool height yet, can not use it as data source")
				}
			}
			pruningCount = pruningCount + float64(supervysorConfig.Interval)/60/60
			time.Sleep(time.Second * time.Duration(supervysorConfig.Interval))
		}
	},
}
