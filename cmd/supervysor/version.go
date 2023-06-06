package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of KSYNC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ksync version: %s\n", getVersion())
	},
}

func getVersion() string {
	version, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to get ksync version")
	}

	return strings.TrimSpace(version.Main.Version)
}
