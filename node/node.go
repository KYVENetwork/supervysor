package node

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/KYVENetwork/supervysor/settings"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/supervysor/node/helpers"
)

var logger = log.NewLogger(os.Stdout)

type Response struct {
	Result struct {
		Response struct {
			LastBlockHeight string `json:"last_block_height"`
		} `json:"response"`
	} `json:"result"`
}

type ProcessType struct {
	Id        int
	GhostMode bool
}

var Process = ProcessType{
	Id:        0,
	GhostMode: true,
}

func GetNodeHeight() int {
	// TODO(@christopher) Complete error handling

	if Process.Id == 0 {
		logger.Info("node hasn't started yet. Try again in 5s ...")

		time.Sleep(time.Second * 5)
		GetNodeHeight()
	}

	abciEndpoint := "http://localhost:26657/abci_info?"
	response, err := http.Get(abciEndpoint)

	if err != nil {
		logger.Error("failed to query height. Try again in 5s ...")

		time.Sleep(time.Second * 5)
		return GetNodeHeight()
	} else {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Error("could not read response data", "err", err.Error())
		}

		var resp Response
		err = json.Unmarshal(responseData, &resp)
		if err != nil {
			logger.Error("could not unmarshal JSON", "err", err.Error())
		}

		lastBlockHeight := resp.Result.Response.LastBlockHeight
		nodeHeight, err := strconv.Atoi(lastBlockHeight)
		if err != nil {
			logger.Error("could not convert lastBlockHeight to str", "err", err.Error())
		}

		return nodeHeight
	}
}

func startNode(initial bool, binaryPath string, seeds string) (*os.Process, error) {
	if !initial {
		helpers.MoveAddressBook(Process.GhostMode)
	}

	if !(Process.Id == 0 && Process.GhostMode) && !initial {
		// TODO(@christopher): Panic and stop all processes
		return nil, nil
	} else {
		// TODO(@christopher): Support pruning choice between default, everything and custom pruning setting with recommended settings on default (keeping 1 week backup based on pool settings).
		// TODO(@christopher): Support pruning for state requests (e.g. spot-price -> pricing integration).

		cmdPath, err := exec.LookPath(binaryPath)
		if err != nil {
			logger.Error("could not resolve binary path", "err", err)
			return nil, err
		}

		var args []string

		if strings.HasSuffix(binaryPath, "/cosmovisor") {
			args = []string{
				"run",
				"start",
				"--p2p.seeds",
				seeds,
			}
		} else {
			args = []string{
				"start",
				"--p2p.seeds",
				seeds,
			}
		}

		mergedArgs := append(args, settings.PruningCommands...)

		cmd := exec.Command(cmdPath, mergedArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				logger.Error("could not start Normal Mode process", "err", err.Error())
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				// Process can only be stopped through an error, which is why we don't need to log it
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, errors.New("couldn't start running the node")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			logger.Error("could not find started process", "err", err.Error())
			return nil, err
		}

		return process, nil
	}
}

func startGhostNode(binaryPath string) (*os.Process, error) {
	helpers.MoveAddressBook(Process.GhostMode)

	if !(Process.Id == 0 && !Process.GhostMode) {
		// TODO(@christopher): Panic and stop all processes
		return nil, nil
	} else {

		cmdPath, err := exec.LookPath(binaryPath)
		if err != nil {
			logger.Error("could not resolve binary path", "err", err)
			return nil, err
		}

		var args []string

		if strings.HasSuffix(binaryPath, "/cosmovisor") {
			args = []string{
				"run",
				"start",
				"--p2p.seeds",
				" ",
				"--p2p.laddr",
				// TODO(@christopher): Find unused port
				"tcp://0.0.0.0:26658",
			}
		} else {
			args = []string{
				"start",
				"--p2p.seeds",
				" ",
				"--p2p.laddr",
				// TODO(@christopher): Find unused port
				"tcp://0.0.0.0:26658",
			}
		}

		cmd := exec.Command(cmdPath, args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				logger.Error("could not start Ghost Node process", "err", err.Error())
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				// Process can only be stopped through an error, which is why we don't need to log it
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, errors.New("couldn't start running the node")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			logger.Error("could not find started process")
			return nil, err
		}

		return process, nil
	}
}

func shutdownNode() {
	process, err := os.FindProcess(Process.Id)
	if err != nil {
		logger.Error("could not find process to shutdown", "err", err)
		// TODO(@christopher): Panic and shutdown all running processes
	}

	// Terminate the process
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// TODO(@christopher): Panic and shutdown all running processes
		logger.Error("could not terminate process", "err", err)
		return
	}

	Process.Id = 0

	logger.Info("process terminated successfully")
}
