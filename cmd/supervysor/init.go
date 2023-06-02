package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/spf13/cobra"
)

type Config struct {
	ChainId             string
	Interval            int
	HeightDifferenceMax int
	HeightDifferenceMin int
}

var (
	chainId             string
	interval            int
	heightDifferenceMax int
	heightDifferenceMin int

	cfg Config
)

func init() {
	initCmd.Flags().StringVar(&chainId, "chain-id", "", "KYVE chain-id")
	if err := initCmd.MarkFlagRequired("chain-id"); err != nil {
		panic(fmt.Errorf("flag 'chain-id' should be required: %w", err))
	}

	initCmd.Flags().IntVar(&interval, "interval", 10, "interval to check height difference in seconds")

	initCmd.Flags().IntVar(&heightDifferenceMax, "height-difference-max", 10000, "max difference of pool-height and node-height to enable Ghost Mode")

	initCmd.Flags().IntVar(&heightDifferenceMin, "height-difference-min", 5000, "min difference of pool-height and node-height to enable Normal Mode")

}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize supverysor.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return InitializeSupervysor()
	},
}

func InitializeSupervysor() error {
	configPath, err := getSupervysorDir()
	if err != nil {
		logger.Error("couldn't get supervysor directory path")
		return err
	}

	if _, err := os.Stat(configPath + "/config.toml"); err == nil {
		logger.Info(fmt.Sprintf("supervysor was already initialized and is editable under %s/config.toml", configPath))
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(configPath, 0o755)
			if err != nil {
				return err
			}
		}
		logger.Info("initializing supverysor...")
		config := Config{
			ChainId:             chainId,
			Interval:            interval,
			HeightDifferenceMax: heightDifferenceMax,
			HeightDifferenceMin: heightDifferenceMin,
		}
		b, err := toml.Marshal(config)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(configPath+"/config.toml", b, 0o755)
		if err != nil {
			logger.Error("couldn't write config file")
			panic(err)
		}

		_, err = getConfig()
		if err != nil {
			logger.Error("couldn't load config file")
			return err
		}

		logger.Info(fmt.Sprintf("successfully initialized: config available at %s/config.toml", configPath))
		return nil
	} else {
		return err
	}
}

func getConfig() (*Config, error) {
	configPath, err := getSupervysorDir()
	if err != nil {
		logger.Error("couldn't get supervysor directory path")
		return nil, err
	}

	data, err := os.ReadFile(configPath + "/config.toml")
	if err != nil {
		logger.Error("couldn't find config. Please initialize again")
		return nil, err
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("couldn't unmarshal config")
		panic(err)
	}

	return &cfg, nil
}

func getSupervysorDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error("couldn't find home directory")
		return "", err
	}

	return filepath.Join(home, ".supervysor"), nil
}
