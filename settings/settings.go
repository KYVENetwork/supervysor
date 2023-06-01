package settings

import (
	"fmt"
	"os"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/supervysor/settings/helpers"
)

var logger = log.NewLogger(os.Stdout)

type PoolSettingsType struct {
	MaxBundleSize  int
	UploadInterval int
}

type PruningSettingsType struct {
	Interval   int
	KeepEvery  int
	KeepRecent int
}

var poolSettings = PoolSettingsType{
	MaxBundleSize:  0,
	UploadInterval: 0,
}

var PruningSettings = PruningSettingsType{
	Interval:   10,
	KeepEvery:  0,
	KeepRecent: 0,
}

func InitializeSettings(binaryPath string, poolId int64) error {
	if err := helpers.CheckBinaryPath(binaryPath); err != nil {
		logger.Error("couldn't resolve binary path", err)
		return err
	}

	settings, err := helpers.GetPoolSettings(poolId)
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

	PruningSettings.KeepRecent = keepRecent

	return nil
}
