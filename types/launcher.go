package types

import (
	"cosmossdk.io/log"
)

type Launcher struct {
	Logger log.Logger
	Cfg    *Config
}

func NewLauncher(logger *log.Logger, cfg *Config) (Launcher, error) {
	return Launcher{Logger: *logger, Cfg: cfg}, nil
}
