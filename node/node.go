package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/KYVENetwork/supervysor/types"

	"github.com/KYVENetwork/supervysor/settings"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/supervysor/node/helpers"
)

var logger = log.NewLogger(os.Stdout)

var Process = types.ProcessType{
	Id:        0,
	GhostMode: true,
}

func GetNodeHeight() int {
	if Process.Id == 0 {
		logger.Info("node hasn't started yet. Try again in 5s ...")

		time.Sleep(time.Second * 5)
		return GetNodeHeight()
	}

	response, err := http.Get(types.ABCIEndpoint)

	if err != nil {
		logger.Error("failed to query height. Try again in 5s ...")

		time.Sleep(time.Second * 5)
		return GetNodeHeight()
	} else {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Error("could not read response data", "err", err.Error())
		}

		var resp types.HeightResponse
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

func startNode(initial bool, binaryPath string, addrBookPath string, seeds string) (*os.Process, error) {
	if !initial {
		if err := helpers.MoveAddressBook(false, addrBookPath); err != nil {
			logger.Error("could not move address book", "err", err)
			return nil, err
		}
	}

	if !(Process.Id == 0 && Process.GhostMode) && !initial {
		// TODO(@christopher): Panic and stop all processes
		return nil, nil
	} else {
		cmdPath, err := exec.LookPath(binaryPath)
		if err != nil {
			return nil, fmt.Errorf("could not resolve binary path: %s", err)
		}

		var args []string

		if strings.HasSuffix(binaryPath, "/cosmovisor") {
			args = []string{
				"run",
				"start",
			}
		} else {
			args = []string{
				"start",
			}
		}

		if initial {
			args = append(args, "--p2p-seeds", seeds)
		}

		mergedArgs := append(args, settings.PruningCommands...)

		cmd := exec.Command(cmdPath, mergedArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				logger.Error("could not start Normal Mode process", "err", err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				// Process can only be stopped through an error, not necessary to log an error
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, errors.New("couldn't start running the node")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, fmt.Errorf("could not find started process: %s", err)
		}

		return process, nil
	}
}

func startGhostNode(binaryPath string, addrBookPath string) (*os.Process, error) {
	if err := helpers.MoveAddressBook(true, addrBookPath); err != nil {
		logger.Error("could not move address book", "err", err)
		return nil, err
	}

	logger.Info("address book successfully moved")

	if !(Process.Id == 0 && !Process.GhostMode) {
		// TODO(@christopher): Panic and stop all processes
		return nil, nil
	} else {

		cmdPath, err := exec.LookPath(binaryPath)
		if err != nil {
			return nil, fmt.Errorf("could not resolve binary path: %s", err)
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
				logger.Error("could not start Ghost Node process", "err", err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			err = cmd.Wait()
			if err != nil {
				// Process can only be stopped through an error, which is why we don't need to log it
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, fmt.Errorf("couldn't start running the node")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, fmt.Errorf("could not find started process: %s", err)
		}

		return process, nil
	}
}

func ShutdownNode() error {
	if Process.Id != 0 {
		process, err := os.FindProcess(Process.Id)
		if err != nil {
			return fmt.Errorf("could not find process to shutdown: %s", err)
		}

		if err = process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("could not terminate process: %s", err)
		}

		Process.Id = 0
	}

	return nil
}
