package types

import (
	"cosmossdk.io/log"
)

type Launcher struct {
	Logger  log.Logger
	Cfg     *Config
	Process ProcessType
}

func NewLauncher(logger *log.Logger, cfg *Config) *Launcher {
	return &Launcher{Logger: *logger, Cfg: cfg, Process: ProcessType{Id: 0, GhostMode: false}}
}
