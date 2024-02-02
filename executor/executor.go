package executor

import (
	"fmt"
	"time"

	"github.com/KYVENetwork/supervysor/store"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/supervysor/node"
	"github.com/KYVENetwork/supervysor/types"
)

type Executor struct {
	Logger  log.Logger
	Cfg     *types.SupervysorConfig
	Process types.ProcessType
}

func NewExecutor(logger *log.Logger, cfg *types.SupervysorConfig) *Executor {
	return &Executor{Logger: *logger, Cfg: cfg, Process: types.ProcessType{Id: -1, GhostMode: false}}
}

// InitialStart initiates the node by starting it in the initial mode.
func (e *Executor) InitialStart(flags string) error {
	e.Logger.Info("starting initially")
	process, err := node.StartNode(e.Cfg, e.Logger, &e.Process, true, false, flags)
	if err != nil {
		return fmt.Errorf("could not start node initially: %s", err)
	}

	e.Logger.Info("initial process started", "pId", process.Pid)

	e.Process.Id = process.Pid
	e.Process.GhostMode = false

	return nil
}

// EnableGhostMode activates the Ghost Mode by starting the node in GhostMode if it is not already enabled.
// If not, it shuts down the node running in NormalMode, initiates the GhostMode and updates the process ID
// and GhostMode upon success.
func (e *Executor) EnableGhostMode(flags string) error {
	if !e.Process.GhostMode {
		if err := node.ShutdownNode(&e.Process); err != nil {
			e.Logger.Error("could not shutdown node", "err", err)
		}
		e.Logger.Info("successfully shut down node", "mode", "normal")

		time.Sleep(time.Second * time.Duration(10))

		process, err := node.StartGhostNode(e.Cfg, e.Logger, &e.Process, false, flags)
		if err != nil {
			return fmt.Errorf("Ghost Mode enabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				e.Process.Id = process.Pid
				e.Process.GhostMode = true
				e.Logger.Info("node started in Ghost Mode")
			} else {
				return fmt.Errorf("enabling Ghost Mode failed: process is not defined")
			}
		}
	}
	return nil
}

// EnableNormalMode enables the Normal Mode by starting the node in NormalMode if it is not already enabled.
// If the GhostMode is active, it shuts down the node, starts the NormalMode with the provided parameters
// and updates the process ID and GhostMode upon success.
func (e *Executor) EnableNormalMode(flags string) error {
	if e.Process.GhostMode {
		if err := node.ShutdownNode(&e.Process); err != nil {
			e.Logger.Error("could not shutdown node", "err", err)
		}
		e.Logger.Info("successfully shut down node", "mode", "ghost")

		time.Sleep(time.Second * time.Duration(10))

		process, err := node.StartNode(e.Cfg, e.Logger, &e.Process, false, false, flags)
		if err != nil {
			return fmt.Errorf("Ghost Mode disabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				e.Process.Id = process.Pid
				e.Process.GhostMode = false
				e.Logger.Info("Node started in Normal Mode", "pId", process.Pid)
			} else {
				return fmt.Errorf("Ghost Mode disabling failed: process is not defined")
			}
		}
	}
	return nil
}

func (e *Executor) PruneData(homePath string, pruneHeight int, statePruning, forceCompact bool, flags string) error {
	if err := e.Shutdown(); err != nil {
		e.Logger.Error("could not shutdown node process", "err", err)
		return err
	}
	err := store.Prune(homePath, int64(pruneHeight)-1, statePruning, forceCompact, e.Logger)
	if err != nil {
		e.Logger.Error("could not prune, exiting")
		return err
	}

	time.Sleep(time.Second * time.Duration(30))

	if e.Process.GhostMode {
		process, err := node.StartGhostNode(e.Cfg, e.Logger, &e.Process, true, flags)
		if err != nil {
			return fmt.Errorf("Ghost Mode enabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				e.Process.Id = process.Pid
				e.Process.GhostMode = true
				e.Logger.Info("node started in GhostMode after pruning")
			} else {
				return fmt.Errorf("enabling Ghost Mode failed: process is not defined")
			}
		}
	} else {
		process, err := node.StartNode(e.Cfg, e.Logger, &e.Process, false, true, flags)
		if err != nil {
			return fmt.Errorf("Ghost Mode disabling failed: %s", err)
		} else {
			if process != nil && process.Pid > 0 {
				e.Process.Id = process.Pid
				e.Process.GhostMode = false
				e.Logger.Info("Node started in Normal Mode after pruning", "pId", process.Pid)
			} else {
				return fmt.Errorf("GhostMode disabling failed: process is not defined")
			}
		}
	}
	return nil
}

func (e *Executor) GetHeight() (int, error) {
	return node.GetNodeHeight(e.Logger, &e.Process, e.Cfg.ABCIEndpoint)
}

func (e *Executor) Shutdown() error {
	return node.ShutdownNode(&e.Process)
}
