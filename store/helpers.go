package store

import (
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/tendermint/tendermint/state"
	tmStore "github.com/tendermint/tendermint/store"
	db "github.com/tendermint/tm-db"
)

func PruneGoLevelDB(home string, untilHeight int64, statePruning bool, logger log.Logger) error {
	o := opt.Options{
		DisableSeeksCompaction: true,
	}

	// Get BlockStore
	blockStoreDB, err := db.NewGoLevelDBWithOpts("blockstore", filepath.Join(home, "data"), &o)
	if err != nil {
		return fmt.Errorf("failed to create blockStoreDB: %w", err)
	}
	blockStore := tmStore.NewBlockStore(blockStoreDB)

	// Get StateStore
	stateDB, err := db.NewGoLevelDBWithOpts("state", filepath.Join(home, "data"), &o)
	if err != nil {
		return fmt.Errorf("failed to create stateStoreDB: %w", err)
	}
	stateStore := state.NewStore(stateDB)

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

	if err = blockStoreDB.ForceCompact(nil, nil); err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	if statePruning {
		err = stateStore.PruneStates(base, untilHeight)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(0)
		}

		if err = stateDB.ForceCompact(nil, nil); err != nil {
			logger.Error(err.Error())
			os.Exit(0)
		}

		logger.Info(fmt.Sprintf("Pruned %d blocks and the state until %d, new base height is %d", blocks, untilHeight, blockStore.Base()))
	} else {
		logger.Info(fmt.Sprintf("Pruned %d blocks until %d, new base height is %d", blocks, untilHeight, blockStore.Base()))
	}

	return nil
}

func PruneAnyDB(home string, untilHeight int64, statePruning bool, logger log.Logger) error {
	config, err := helpers.LoadConfig(home)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	blockStoreDB, blockStore, err := GetBlockstoreDBs(config)
	defer func(blockStoreDB db.DB) {
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
		defer func(stateStoreDB db.DB) {
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
