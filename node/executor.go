package node

import (
	"fmt"

	"github.com/KYVENetwork/supervysor/node/helpers"
)

// InitialStart initiates the node by starting it in the initial mode.
func InitialStart(logFile string, binaryPath string, homePath string, seeds string) error {
	logger = helpers.InitLogger(logFile)

	logger.Info("starting initially")
	process, err := startNode(logFile, true, binaryPath, homePath, seeds)
	if err != nil {
		return fmt.Errorf("could not start node initially: %s", err)
	}

	logger.Info("initial process started", "pId", process.Pid)

	Process.Id = process.Pid
	Process.GhostMode = false

	return nil
}

// EnableGhostMode activates the Ghost Mode by starting the node in GhostMode if it is not already enabled.
// If not, it shuts down the node running in NormalMode, initiates the GhostMode and updates the process ID
// and GhostMode upon success.
func EnableGhostMode(logFile string, binaryPath string, homePath string) error {
	logger = helpers.InitLogger(logFile)

	if !Process.GhostMode {
		if err := ShutdownNode(); err != nil {
			logger.Error("could not shutdown node", "err", err)
		}

		process, err := startGhostNode(logFile, binaryPath, homePath)
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
	}
	return nil
}

// EnableNormalMode enables the Normal Mode by starting the node in NormalMode if it is not already enabled.
// If the GhostMode is active, it shuts down the node, starts the NormalMode with the provided parameters
// and updates the process ID and GhostMode upon success.
func EnableNormalMode(logFile string, binaryPath string, homePath string, seeds string) error {
	logger = helpers.InitLogger(logFile)

	if Process.GhostMode {
		if err := ShutdownNode(); err != nil {
			logger.Error("could not shutdown node", "err", err)
		}

		process, err := startNode(logFile, false, binaryPath, homePath, seeds)
		if err != nil {
			return fmt.Errorf("Ghost Mode disabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = false
				logger.Info("Node started in Normal Mode", "pId", process.Pid)
			} else {
				return fmt.Errorf("Ghost Mode disabling failed: process is not defined")
			}
		}
	}
	return nil
}
