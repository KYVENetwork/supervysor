package main

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/helpers"

	"cosmossdk.io/log"

	"github.com/spf13/cobra"
)

var logger = log.NewLogger(os.Stdout)

var Version = ""

var supervysor = &cobra.Command{
	Use:     "supervysor",
	Short:   "Supervysor helps sync a Tendermint node used as a KYVE data source.",
	Version: Version,
}

var LogFilePath string

func main() {
	logsDir, err := helpers.GetLogsDir()
	if err != nil {
		panic(err)
	}
	LogFilePath = filepath.Join(logsDir, time.Now().Format("20060102_150405")+".log")

	file, err := os.OpenFile(LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		panic(err)
	}

	multiLogger := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)

	logger = log.NewCustomLogger(zerolog.New(multiLogger).With().Timestamp().Logger())

	supervysor.AddCommand(initCmd)
	supervysor.AddCommand(startCmd)
	supervysor.AddCommand(versionCmd)

	if err := supervysor.Execute(); err != nil {
		os.Exit(1)
	}
}
