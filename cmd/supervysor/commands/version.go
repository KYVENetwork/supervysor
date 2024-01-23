package commands

import (
	"github.com/KYVENetwork/supervysor/utils"
	"github.com/spf13/cobra"
)

func init() {
	versionCmd.Flags().BoolVar(&optOut, "opt-out", false, "disable the collection of anonymous usage data")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of supervysor",
	Run: func(cmd *cobra.Command, args []string) {
		utils.TrackVersionEvent(optOut)

		logger.Info(Version)
	},
}
