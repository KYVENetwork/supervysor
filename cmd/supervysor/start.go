package main

import (
	"github.com/KYVENetwork/supervysor/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/KYVENetwork/supervysor/executor"
	"github.com/KYVENetwork/supervysor/pool"
)

func NewMetrics(reg prometheus.Registerer) *types.Metrics {
	m := &types.Metrics{
		PoolHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "pool_height",
			Help:      "Height of the specified KYVE data pool.",
		}),
		NodeHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "node_height",
			Help:      "Height of the running data source node.",
		}),
	}
	reg.MustRegister(m.PoolHeight, m.NodeHeight)
	return m
}

// The startCmd of the supervysor launches and manages the node process using the specified binary.
// It periodically retrieves the heights of the node and the associated KYVE pool, and dynamically adjusts
// the sync mode of the node based on these heights.
var startCmd = &cobra.Command{
	Use:                "start",
	Short:              "Start a supervysed Tendermint node",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, flags []string) error {
		// Create Prometheus registry
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		// Create metrics endpoint
		promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
		http.Handle("/metrics", promHandler)
		err := http.ListenAndServe(":26660", nil)
		if err != nil {
			logger.Error("could not start metrics server", "err", err)
			return err
		}

		// Load initialized config.
		config, err := getConfig()
		if err != nil {
			logger.Error("could not load config", "err", err)
			return err
		}

		e := executor.NewExecutor(&logger, config)

		// Start data source node initially.
		if err := e.InitialStart(flags); err != nil {
			logger.Error("initial start failed", "err", err)
			return err
		}

		currentMode := "normal"

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
			m.NodeHeight.Set(float64(nodeHeight))

			poolHeight, err := pool.GetPoolHeight(config.ChainId, config.PoolId, config.FallbackEndpoints)
			if err != nil {
				logger.Error("could not get pool height", "err", err)
				if shutdownErr := e.Shutdown(); shutdownErr != nil {
					logger.Error("could not shutdown node process", "err", shutdownErr)
				}
				return err
			}
			m.PoolHeight.Set(float64(*poolHeight))

			logger.Info("fetched heights successfully", "node", nodeHeight, "pool", poolHeight, "max-height", *poolHeight+config.HeightDifferenceMax, "min-height", *poolHeight+config.HeightDifferenceMin)

			// Calculate height difference to enable the correct mode.
			heightDiff := nodeHeight - *poolHeight

			if heightDiff >= config.HeightDifferenceMax {
				if currentMode != "ghost" {
					logger.Info("enabling GhostMode")
					currentMode = "ghost"
				} else {
					logger.Info("keeping GhostMode")
				}
				// Data source node has synced far enough, enable or keep Ghost Mode
				if err = e.EnableGhostMode(flags); err != nil {
					logger.Error("could not enable Ghost Mode", "err", err)

					if shutdownErr := e.Shutdown(); shutdownErr != nil {
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
				} else {
					logger.Info("keeping NormalMode")
				}
				// Difference is < HeightDifferenceMin, Data source needs to catch up, enable or keep Normal Mode
				if err = e.EnableNormalMode(flags); err != nil {
					logger.Error("could not enable Normal Mode", "err", err)

					if shutdownErr := e.Shutdown(); shutdownErr != nil {
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
