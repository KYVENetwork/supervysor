package main

import (
	"fmt"

	"github.com/KYVENetwork/supervysor/backup"
	"github.com/KYVENetwork/supervysor/cmd/supervysor/helpers"
	"github.com/spf13/cobra"
)

var (
	compressionType string
	destPath        string
	maxBackups      int
	srcPath         string
)

func init() {
	backupCmd.Flags().StringVar(&srcPath, "src-path", "", "source path of the directory to backup")
	if err := backupCmd.MarkFlagRequired("src-path"); err != nil {
		panic(fmt.Errorf("flag 'src-path' should be required: %w", err))
	}

	backupCmd.Flags().StringVar(&destPath, "dest-path", "", "destination path of the written backup (default '~/.ksync/backups)'")

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

		if destPath == "" {
			d, err := helpers.CreateDestPath(backupDir)
			if err != nil {
				return
			}
			destPath = d
		}

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
