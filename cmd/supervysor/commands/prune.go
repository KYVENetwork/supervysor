package commands

import (
	"fmt"

	"github.com/KYVENetwork/supervysor/utils"

	"github.com/KYVENetwork/supervysor/store"
	"github.com/spf13/cobra"
)

func init() {
	pruneCmd.Flags().StringVar(&home, "home", "", "path to home directory (e.g. /root/.osmosisd)")
	if err := pruneCmd.MarkFlagRequired("home"); err != nil {
		panic(fmt.Errorf("flag 'home' should be required: %w", err))
	}

	pruneCmd.Flags().Int64Var(&untilHeight, "until-height", 0, "prune until specified height (excluding)")
	if err := pruneCmd.MarkFlagRequired("until-height"); err != nil {
		panic(fmt.Errorf("flag 'until-height' should be required: %w", err))
	}

	pruneCmd.Flags().BoolVar(&statePruning, "state-pruning", true, "enable state pruning")

	startCmd.Flags().BoolVar(&forceCompact, "force-compact", false, "prune with ForceCompact enabled")

	pruneCmd.Flags().BoolVar(&optOut, "opt-out", false, "disable the collection of anonymous usage data")
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune blocks and optionally state from base height until a specific height",
	Run: func(cmd *cobra.Command, args []string) {
		utils.TrackPruneEvent(optOut)

		if err := store.Prune(home, untilHeight, statePruning, forceCompact, logger); err != nil {
			logger.Error(err.Error())
		}
	},
}
