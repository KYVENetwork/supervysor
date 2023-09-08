package settings

import (
	"fmt"
	"strings"

	"github.com/KYVENetwork/supervysor/settings/helpers"
	"github.com/KYVENetwork/supervysor/types"
)

var poolSettings = types.PoolSettingsType{
	MaxBundleSize:  0,
	UploadInterval: 0,
}

var Settings = types.SettingsType{
	MaxDifference: 0,
	Seeds:         "",
	Interval:      10,
	KeepRecent:    0,
}

// InitializeSettings initializes the required settings for the supervysor. It performs checks on the binaryPath
// and homePath and sets the seeds value required for the node. It retrieves the pool settings, calculates the
// keepRecent and maxDifference values, and sets the pruning settings based on these calculated values.
// If any step encounters an error, it returns the corresponding error message.
func InitializeSettings(binaryPath string, homePath string, poolId int, stateRequests bool, seeds string, chainId string, fallbackEndpoints string) error {
	if err := helpers.CheckBinaryPath(binaryPath); err != nil {
		return fmt.Errorf("could not resolve binary path: %s", err)
	}

	if err := helpers.CheckFilePath(homePath); err != nil {
		return fmt.Errorf("could not resolve home-path: %s", err)
	}

	Settings.Seeds = seeds
	if strings.TrimSpace(seeds) == "" {
		return fmt.Errorf("seeds are not defined")
	}

	settings, err := helpers.GetPoolSettings(poolId, chainId, fallbackEndpoints)
	if err != nil {
		return fmt.Errorf("could not get pool settings: %s", err)
	}

	poolSettings.MaxBundleSize = settings[0]
	poolSettings.UploadInterval = settings[1]

	keepRecent := helpers.CalculateKeepRecent(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if keepRecent == 0 {
		return fmt.Errorf("keep-recent calculation failed, poolSettings are probably not correctly set")
	}

	maxDifference := helpers.CalculateMaxDifference(poolSettings.MaxBundleSize, poolSettings.UploadInterval)

	if maxDifference == 0 {
		return fmt.Errorf("max-difference calculation failed, poolSettings are probably not correctly set")
	}
	Settings.MaxDifference = maxDifference

	if maxDifference > keepRecent {
		return fmt.Errorf("max-difference can not be > keep-recent")
	}

	if err = helpers.SetPruningSettings(homePath, stateRequests, keepRecent, Settings.Interval); err != nil {
		return fmt.Errorf("could not set pruning settings: %s", err)
	}

	return nil
}
