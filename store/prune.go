package store

import (
	"fmt"
	"strings"

	"github.com/KYVENetwork/supervysor/utils"

	"cosmossdk.io/log"
)

func Prune(home string, untilHeight int64, statePruning, forceCompact bool, logger log.Logger) error {
	config, err := utils.LoadConfig(home)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !forceCompact {
		logger.Info("prune without ForceCompact")
		return PruneAnyDB(home, untilHeight, statePruning, logger)
	}

	if strings.ToLower(config.DBBackend) == "goleveldb" {
		logger.Info("GoLevelDB detected, using ForceCompact")
		return PruneGoLevelDB(home, untilHeight, statePruning, logger)
	} else {
		logger.Info("another DB backend than GoLevelDB detected, ForceCompact not supported")
		return PruneAnyDB(home, untilHeight, statePruning, logger)
	}
}
