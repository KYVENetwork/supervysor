package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of supervysor",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info(fmt.Sprintf("version: %s", Version))
	},
}
