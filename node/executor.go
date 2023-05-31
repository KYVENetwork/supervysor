package node

import (
	"os"
)

func InitialStart(seeds string) (int, error) {
	logger.Info("starting initially")
	process, err := startNode(true, seeds)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
		return 0, err
	}

	logger.Info("initial process started, PID: ", process.Pid)

	Process.Id = process.Pid
	Process.GhostMode = false

	return process.Pid, nil
}

func EnableGhostMode() {
	if !Process.GhostMode {
		logger.Info("enabling Ghost Mode")
		shutdownNode()

		process, err := startGhostNode()
		if err != nil {
			logger.Error("Ghost Mode enabling failed", err)
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				logger.Info("Ghost Node started (PID: %d)\n", process.Pid)
			} else {
				// TODO(@christopher): Panic and shutdown all processes
				logger.Error("Ghost Mode enabling failed.")
			}
		}
	} else {
		logger.Info("keeping Ghost Mode enabled")
	}
}

func DisableGhostMode(seeds string) {
	if Process.GhostMode {
		logger.Info("disabling Ghost Mode")

		shutdownNode()

		process, err := startNode(false, seeds)
		if err != nil {
			logger.Error("Ghost Mode disabling failed", err)
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				logger.Info("Normal Node started (PID: %d)\n", process.Pid)
			} else {
				// TODO(@christopher): Panic and shutdown all processes
				logger.Error("Ghost Mode disabling failed")
			}
		}
	} else {
		logger.Info("keeping Normal Mode enabled")
	}
}
