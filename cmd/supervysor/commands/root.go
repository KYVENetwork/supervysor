package commands

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/KYVENetwork/supervysor/types"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"

	"cosmossdk.io/log"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	abciEndpoint      string
	binary            string
	binaryFlags       string
	cfgFlag           string
	chainId           string
	config            string
	fallbackEndpoints string
	home              string
	metrics           bool
	metricsPort       int
	optOut            bool
	poolId            int
	pruningInterval   int
	seeds             string
	statePruning      bool
	untilHeight       int64

	cfg types.SupervysorConfig

	err        error
	configPath string

	compressionType string
	destPath        string
	maxBackups      int
)

var logger = log.NewLogger(os.Stdout)

var supervysor = &cobra.Command{
	Use:     "supervysor",
	Short:   "Supervysor helps sync a Tendermint node used as a KYVE data source.",
	Version: types.Version,
}

func Execute() {
	logsDir, err := helpers.GetLogsDir()
	if err != nil {
		panic(err)
	}
	logFilePath := filepath.Join(logsDir, time.Now().Format("20060102_150405")+".log")

	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout)
	customConsoleWriter := zerolog.ConsoleWriter{Out: writer}
	customConsoleWriter.FormatCaller = func(i interface{}) string {
		return "\x1b[36m[supervysor]\x1b[0m"
	}
	multiLogger := io.MultiWriter(customConsoleWriter, file)
	logger = log.NewCustomLogger(zerolog.New(multiLogger).With().Timestamp().Logger())

	supervysor.AddCommand(initCmd)
	supervysor.AddCommand(startCmd)
	supervysor.AddCommand(versionCmd)
	supervysor.AddCommand(pruneCmd)
	supervysor.AddCommand(backupCmd)

	if err = supervysor.Execute(); err != nil {
		os.Exit(1)
	}
}
