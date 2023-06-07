package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of supervysor",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Sprintf("supervysor version: %s", Version)
	},
}
