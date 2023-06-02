package settings

import (
	"fmt"
	"os"
	"strconv"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/supervysor/settings/helpers"
)

var logger = log.NewLogger(os.Stdout)

type PoolSettingsType struct {
	MaxBundleSize  int
	UploadInterval int
}

type SettingsType struct {
	MaxDifference int
	Seeds         string
	Interval      int
	KeepEvery     int
	KeepRecent    int
}

var poolSettings = PoolSettingsType{
	MaxBundleSize:  0,
	UploadInterval: 0,
}

// TODO(@christopher): Integrate into config.toml
var Settings = SettingsType{
	MaxDifference: 0,
	Seeds:         "",
	Interval:      10,
	KeepEvery:     0,
	KeepRecent:    0,
}

var PruningCommands []string

func InitializeSettings(binaryPath string, poolId int, stateRequests bool, seeds string, chainId string) error {
	if err := helpers.CheckBinaryPath(binaryPath); err != nil {
		logger.Error("couldn't resolve binary path", err)
		return err
	}

	Settings.Seeds = seeds
	if seeds == "" {
		return fmt.Errorf("seeds are not defined")
	}

	settings, err := helpers.GetPoolSettings(poolId, chainId)
	if err != nil {
		logger.Error("couldn't get pool settings")
		return err
	}
	poolSettings.MaxBundleSize = settings[0]
	poolSettings.UploadInterval = settings[1]

	keepRecent := helpers.CalculateKeepRecent(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if keepRecent == 0 {
		logger.Error("couldn't calculate keep-recent pruning settings")
		return fmt.Errorf("keep-recent calculation failed, poolSettings are probably not correctly set")
	}
	Settings.KeepRecent = keepRecent

	maxDifference := helpers.CalculateMaxDifference(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if maxDifference == 0 {
		logger.Error("couldn't calculate max-difference pruning settings")
		return fmt.Errorf("max-difference calculation failed, poolSettings are probably not correctly set")
	}
	Settings.MaxDifference = maxDifference

	if stateRequests {
		PruningCommands = []string{
			"--pruning",
			"custom",
			"--pruning-keep-every",
			strconv.Itoa(Settings.KeepEvery),
			"--pruning-keep-recent",
			strconv.Itoa(Settings.KeepRecent),
			"--pruning-interval",
			strconv.Itoa(Settings.Interval),
		}
	} else {
		PruningCommands = []string{
			"--pruning",
			"everything",
			"--min-retain-blocks",
			strconv.Itoa(Settings.KeepRecent),
		}
	}

	return nil
}
