package store

import (
	"fmt"
	"os"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"

	"cosmossdk.io/log"

	dbm "github.com/tendermint/tm-db"
)

func Prune(home string, untilHeight int64, statePruning bool, logger log.Logger) error {
	config, err := helpers.LoadConfig(home)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	blockStoreDB, blockStore, err := GetBlockstoreDBs(config)
	defer func(blockStoreDB dbm.DB) {
		err = blockStoreDB.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(0)
		}
	}(blockStoreDB)

	if err != nil {
		panic(fmt.Errorf("failed to load blockstore db: %w", err))
	}

	base := blockStore.Base()

	logger.Info("blockstore base", "base", base)

	if untilHeight < base {
		logger.Error("Error: base height is higher than prune height", "base height", base, "until height", untilHeight)
		os.Exit(0)
	}

	blocks, err := blockStore.PruneBlocks(untilHeight)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	if statePruning {
		stateStoreDB, stateStore, err := GetStateDBs(config)
		defer func(stateStoreDB dbm.DB) {
			err = stateStoreDB.Close()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(0)
			}
		}(stateStoreDB)

		if err = stateStore.PruneStates(base, untilHeight); err != nil {
			logger.Error(err.Error())
			os.Exit(0)
		}

		logger.Info(fmt.Sprintf("Pruned %d blocks and the state until %d, new base height is %d", blocks, untilHeight, blockStore.Base()))
	} else {
		logger.Info(fmt.Sprintf("Pruned %d blocks until %d, new base height is %d", blocks, untilHeight, blockStore.Base()))
	}

	return nil
}
