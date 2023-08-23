package store

import (
	"fmt"
	"os"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/helpers"
	dbm "github.com/tendermint/tm-db"
)

func PruneBlocks(home string, untilHeight int64, logger log.Logger) error {
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

	logger.Info(fmt.Sprintf("Pruned %d blocks, new base height is %d", blocks, blockStore.Base()))

	return nil
}
