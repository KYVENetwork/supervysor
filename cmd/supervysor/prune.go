package main

import (
	"fmt"

	"github.com/KYVENetwork/supervysor/store"
	"github.com/spf13/cobra"
)

var untilHeight int64

func init() {
	pruneCmd.Flags().StringVar(&home, "home", "", "home directory")
	if err := pruneCmd.MarkFlagRequired("home"); err != nil {
		panic(fmt.Errorf("flag 'home' should be required: %w", err))
	}

	pruneCmd.Flags().Int64Var(&untilHeight, "until-height", 0, "prune blocks until this height (excluding)")
	if err := pruneCmd.MarkFlagRequired("until-height"); err != nil {
		panic(fmt.Errorf("flag 'until-height' should be required: %w", err))
	}
}

var pruneCmd = &cobra.Command{
	Use:   "prune-blocks",
	Short: "Prune blocks until a specific height",
	Run: func(cmd *cobra.Command, args []string) {
		if err := store.PruneBlocks(home, untilHeight, logger); err != nil {
			logger.Error(err.Error())
		}
	},
}
