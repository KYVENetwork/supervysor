package settings

import (
	"fmt"
	"strconv"

	"github.com/KYVENetwork/supervysor/settings/helpers"
	"github.com/KYVENetwork/supervysor/types"
)

var poolSettings = types.PoolSettingsType{
	MaxBundleSize:  0,
	UploadInterval: 0,
}

// TODO(@christopher): Integrate into config.toml
var Settings = types.SettingsType{
	MaxDifference: 0,
	Seeds:         "",
	Interval:      10,
	KeepEvery:     0,
	KeepRecent:    0,
}

var PruningCommands []string

func InitializeSettings(binaryPath string, addrBookPath string, poolId int, stateRequests bool, seeds string, chainId string) error {
	if err := helpers.CheckBinaryPath(binaryPath); err != nil {
		return fmt.Errorf("could not resolve binary path: %s", err)
	}

	if err := helpers.CheckFilePath(addrBookPath); err != nil {
		return fmt.Errorf("could not resolve address book path: %s", err)
	}

	Settings.Seeds = seeds
	if seeds == "" {
		return fmt.Errorf("seeds are not defined")
	}

	settings, err := helpers.GetPoolSettings(poolId, chainId)
	if err != nil {
		return fmt.Errorf("could not get pool settings: %s", err)
	}
	poolSettings.MaxBundleSize = settings[0]
	poolSettings.UploadInterval = settings[1]

	keepRecent := helpers.CalculateKeepRecent(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if keepRecent == 0 {
		return fmt.Errorf("keep-recent calculation failed, poolSettings are probably not correctly set")
	}
	Settings.KeepRecent = keepRecent

	maxDifference := helpers.CalculateMaxDifference(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if maxDifference == 0 {
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
