package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"

	"github.com/KYVENetwork/supervysor/utils"

	"golang.org/x/exp/slices"

	"github.com/KYVENetwork/supervysor/settings"
	"github.com/KYVENetwork/supervysor/types"

	"github.com/pelletier/go-toml/v2"

	"github.com/spf13/cobra"
)

func init() {
	initCmd.Flags().StringVarP(&binary, "binary", "b", "", "path to chain binaries or cosmovisor (e.g. /root/go/bin/cosmovisor)")
	if err := initCmd.MarkFlagRequired("binary"); err != nil {
		panic(fmt.Errorf("flag 'binary-path' should be required: %w", err))
	}

	initCmd.Flags().StringVar(&home, "home", "", "path to home directory (e.g. /root/.osmosisd)")

	initCmd.Flags().StringVar(&config, "config", "", "path to config directory (default: ~/.supervysor/")

	initCmd.Flags().IntVar(&poolId, "pool-id", 0, "KYVE pool-id")
	if err := initCmd.MarkFlagRequired("pool-id"); err != nil {
		panic(fmt.Errorf("flag 'pool-id' should be required: %w", err))
	}

	initCmd.Flags().StringVar(&seeds, "seeds", "", "seeds for the node to connect")
	if err := initCmd.MarkFlagRequired("seeds"); err != nil {
		panic(fmt.Errorf("flag 'seeds' should be required: %w", err))
	}

	initCmd.Flags().StringVar(&chainId, "chain-id", "kyve-1", "KYVE chain-id")

	initCmd.Flags().BoolVar(&optOut, "opt-out", false, "disable the collection of anonymous usage data")

	initCmd.Flags().StringVar(&fallbackEndpoints, "fallback-endpoints", "", "additional endpoints to query KYVE pool height")

	initCmd.Flags().IntVar(&pruningInterval, "pruning-interval", 24, "block-pruning interval (hours)")

	initCmd.Flags().BoolVar(&statePruning, "state-pruning", true, "state pruning enabled")

	initCmd.Flags().BoolVar(&metrics, "metrics", false, "exposing Prometheus metrics (true or false)")

	initCmd.Flags().IntVar(&metricsPort, "metrics-port", 26660, "port for metrics server")

	initCmd.Flags().StringVar(&abciEndpoint, "abci-endpoint", "http://127.0.0.1:26657", "ABCI Endpoint to request node information")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize supervysor",
	RunE: func(cmd *cobra.Command, args []string) error {
		supportedChains := []string{"kyve-1", "kaon-1", "korellia", "korellia-2"}
		if !slices.Contains(supportedChains, chainId) {
			logger.Error("specified chain-id is not supported", "chain-id", chainId)
			return fmt.Errorf("not supported chain-id")
		}

		// if no home path was given get the default one
		if home == "" {
			home = helpers.GetHomePathFromBinary(binary)
		}

		if pruningInterval <= 6 {
			logger.Error("pruning-interval should be higher than 6 hours")
		}

		utils.TrackInitEvent(chainId, optOut)

		if err := settings.InitializeSettings(binary, home, poolId, false, seeds, chainId, fallbackEndpoints); err != nil {
			logger.Error("could not initialize settings", "err", err)
			return err
		}
		logger.Info("successfully initialized settings")

		if config == "" {
			configPath, err = helpers.GetSupervysorDir()
			if err != nil {
				logger.Error("could not get supervysor directory path", "err", err)
				return err
			}
		} else {
			configPath = config
		}

		if _, err = os.Stat(configPath + "/config.toml"); err == nil {
			logger.Info(fmt.Sprintf("supervysor was already initialized and is editable under %s/config.toml", configPath))
			return nil
		} else if errors.Is(err, os.ErrNotExist) {
			if _, err = os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
				err = os.Mkdir(configPath, 0o755)
				if err != nil {
					return err
				}
			}
			logger.Info("initializing supverysor...")

			supervysorConfig := types.SupervysorConfig{
				ABCIEndpoint:        abciEndpoint,
				BinaryPath:          binary,
				ChainId:             chainId,
				FallbackEndpoints:   fallbackEndpoints,
				HeightDifferenceMax: settings.Settings.MaxDifference,
				HeightDifferenceMin: settings.Settings.MaxDifference / 2,
				HomePath:            home,
				Interval:            10,
				Metrics:             metrics,
				MetricsPort:         metricsPort,
				PoolId:              poolId,
				PruningInterval:     pruningInterval,
				Seeds:               seeds,
				StatePruning:        statePruning,
				StateRequests:       false,
			}
			b, err := toml.Marshal(supervysorConfig)
			if err != nil {
				logger.Error("could not unmarshal config", "err", err)
				return err
			}

			err = os.WriteFile(configPath+"/config.toml", b, 0o755)
			if err != nil {
				logger.Error("could not write config file", "err", err)
				return err
			}

			_, err = getSupervysorConfig(configPath)
			if err != nil {
				logger.Error("could not load config file", "err", err)
				return err
			}

			logger.Info(fmt.Sprintf("successfully initialized: config available at %s/config.toml", configPath))
			return nil
		} else {
			logger.Error("could not get supervysor directory")
			return err
		}
	},
}

// getSupervysorConfig returns the supervysor config.toml file.
func getSupervysorConfig(configPath string) (*types.SupervysorConfig, error) {
	if !strings.HasSuffix(configPath, "/config.toml") {
		configPath += "/config.toml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not find config. Please initialize again: %s", err)
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unsmarshal config: %s", err)
	}

	return &cfg, nil
}
