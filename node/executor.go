package node

import (
	"fmt"
)

func InitialStart(binaryPath string, addrBookPath string, seeds string) (int, error) {
	logger.Info("starting initially")
	process, err := startNode(true, binaryPath, addrBookPath, seeds)
	if err != nil {
		return 0, fmt.Errorf("could not start node initially: %s", err)
	}

	logger.Info("initial process started", "pId", process.Pid)

	Process.Id = process.Pid
	Process.GhostMode = false

	return process.Pid, nil
}

func EnableGhostMode(binaryPath string, addrBookPath string) error {
	if !Process.GhostMode {
		logger.Info("enabling Ghost Mode")

		if err := ShutdownNode(); err != nil {
			logger.Error("could not shutdown node", "err", err)
		}

		process, err := startGhostNode(binaryPath, addrBookPath)
		if err != nil {
			return fmt.Errorf("Ghost Mode enabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				logger.Info("node started in Ghost Mode")
			} else {
				return fmt.Errorf("Ghost Mode enabling failed: process is not defined")
			}
		}
	} else {
		logger.Info("keeping Ghost Mode enabled")
	}
	return nil
}

func DisableGhostMode(binaryPath string, addrBookPath string, seeds string) error {
	if Process.GhostMode {
		logger.Info("disabling Ghost Mode")

		if err := ShutdownNode(); err != nil {
			logger.Error("could not shutdown node", "err", err)
		}

		process, err := startNode(false, binaryPath, addrBookPath, seeds)
		if err != nil {
			return fmt.Errorf("Ghost Mode disabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				logger.Info("Node started in Normal Mode", "pId", process.Pid)
			} else {
				return fmt.Errorf("Ghost Mode disabling failed: process is not defined")
			}
		}
	} else {
		logger.Info("keeping Normal Mode enabled")
	}
	return nil
}
