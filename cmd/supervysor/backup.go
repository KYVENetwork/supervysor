package main

import (
	"fmt"
	"path/filepath"

	"github.com/KYVENetwork/supervysor/store"

	"github.com/KYVENetwork/supervysor/backup"
	"github.com/KYVENetwork/supervysor/cmd/supervysor/helpers"
	"github.com/spf13/cobra"
)

var (
	compressionType string
	destPath        string
	maxBackups      int
)

func init() {
	backupCmd.Flags().StringVar(&home, "home", "", "path to home directory (e.g. /root/.osmosisd)")
	if err := backupCmd.MarkFlagRequired("home"); err != nil {
		panic(fmt.Errorf("flag 'src-path' should be required: %w", err))
	}

	backupCmd.Flags().StringVar(&destPath, "dest-path", "", "destination path of the written backup (default '~/.supervysor/backups)'")

	backupCmd.Flags().StringVar(&compressionType, "compression", "", "compression type to compress backup directory ['tar.gz', 'zip', '']")

	backupCmd.Flags().IntVar(&maxBackups, "max-backups", 0, "number of kept backups (set 0 to keep all)")
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup data directory",
	Run: func(cmd *cobra.Command, args []string) {
		backupDir, err := helpers.GetBackupDir()
		if err != nil {
			logger.Error("failed to get ksync home directory", "err", err)
			return
		}

		config, err := helpers.LoadConfig(home)
		if err != nil {
			logger.Error("failed to load tendermint config", "err", err)
			return
		}

		// Load block store
		blockStoreDB, blockStore, err := store.GetBlockstoreDBs(config)
		if err != nil {
			logger.Error("failed to get blockstore dbs")
			return
		}
		defer blockStoreDB.Close()

		if destPath == "" {
			logger.Info("height", "h", blockStore.Height())
			d, err := helpers.CreateDestPath(backupDir, blockStore.Height())
			if err != nil {
				logger.Error("could not create destination path", "err", err)
				return
			}
			destPath = d
		}

		// Only backup data directory
		srcPath := filepath.Join(home, "data")

		if err := helpers.ValidatePaths(srcPath, destPath); err != nil {
			return
		}

		logger.Info("starting to copy backup", "from", srcPath, "to", destPath)

		if err := backup.CopyDir(srcPath, destPath); err != nil {
			logger.Error("error copying directory to backup destination", "err", err)
		}

		logger.Info("directory copied successfully")

		if compressionType != "" {
			if err := backup.CompressDirectory(destPath, compressionType); err != nil {
				logger.Error("compression failed", "err", err)
			}

			logger.Info("compressed backup successfully")
		}

		if maxBackups > 0 {
			logger.Info("starting to cleanup backup directory", "path", backupDir)
			if err := backup.ClearBackups(backupDir, maxBackups); err != nil {
				logger.Error("clearing backup directory failed", "err", err)
				return
			}
		}
	},
}
