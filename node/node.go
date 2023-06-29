package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/KYVENetwork/supervysor/types"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/supervysor/node/helpers"
)

var logger log.Logger

var Process = types.ProcessType{
	Id:        0,
	GhostMode: false,
}

// The GetNodeHeight function retrieves the height of the node by querying the ABCI endpoint.
// It uses recursion with a maximum depth of 10 to handle delays or failures.
// It returns the nodeHeight if successful or an error message if the recursion depth reaches the limit (200s).
func GetNodeHeight(logFile string, recursionDepth int) (int, error) {
	logger = helpers.InitLogger(logFile)

	if recursionDepth < 10 {
		if Process.Id == 0 {
			logger.Error(fmt.Sprintf("node hasn't started yet. Try again in 20s ... (%d/10)", recursionDepth+1))

			time.Sleep(time.Second * 20)
			return GetNodeHeight(logFile, recursionDepth+1)
		}

		response, err := http.Get(types.ABCIEndpoint)

		if err != nil {
			logger.Error(fmt.Sprintf("failed to query height. Try again in 20s ... (%d/10)", recursionDepth+1))

			time.Sleep(time.Second * 20)
			return GetNodeHeight(logFile, recursionDepth+1)
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

			return nodeHeight, nil
		}
	} else {
		return 0, fmt.Errorf("could not get node height, exiting ...")
	}
}

// startNode starts the node process in Normal Mode and returns the os.Process object representing
// the running process. It checks if the node is being started initially or not, moves the
// address book if necessary, and sets the appropriate command arguments based on the binaryPath.
func startNode(logFile string, initial bool, binaryPath string, homePath string, seeds string) (*os.Process, error) {
	logger = helpers.InitLogger(logFile)

	addrBookPath := filepath.Join(homePath, "config", "addrbook.json")

	if !initial {
		if err := helpers.MoveAddressBook(false, addrBookPath); err != nil {
			logger.Error("could not move address book", "err", err)
			return nil, err
		}
	}

	// To start the node normally when it's not initially, Process ID needs to be = 0 and GhostMode = true
	if (Process.Id != 0 || !Process.GhostMode) && !initial {
		return nil, fmt.Errorf("process management failed")
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
			args = append(args, "--p2p.seeds", seeds)
		}

		if homePath != "" {
			args = append(args, "--home", homePath)
		}

		cmd := exec.Command(cmdPath, args...)

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

// startGhostNode starts the node process in Ghost Mode and returns the os.Process object
// representing the running process. It moves the address book, checks if the node is already running
// or in Ghost Mode ands sets the appropriate command arguments based on the binaryPath.
// It starts the node without seeds and with a changed laddr, so the node can't continue syncing.
func startGhostNode(logFile string, binaryPath string, homePath string) (*os.Process, error) {
	logger = helpers.InitLogger(logFile)

	addrBookPath := filepath.Join(homePath, "config", "addrbook.json")

	if err := helpers.MoveAddressBook(true, addrBookPath); err != nil {
		logger.Error("could not move address book", "err", err)
		return nil, err
	}

	logger.Info("address book successfully moved")

	// To start the node in GhostMode, Process ID needs to be = 0 and GhostMode = false
	if Process.Id != 0 || Process.GhostMode {
		return nil, fmt.Errorf("process management failed")
	} else {

		cmdPath, err := exec.LookPath(binaryPath)
		if err != nil {
			return nil, fmt.Errorf("could not resolve binary path: %s", err)
		}

		port, err := helpers.GetPort()
		if err != nil {
			return nil, fmt.Errorf("could not find unused port: %s", err)
		}

		laddr := "tcp://0.0.0.0:" + strconv.Itoa(port)

		args := []string{
			"start",
			"--p2p.seeds",
			" ",
			"--p2p.laddr",
			laddr,
		}

		if strings.HasSuffix(binaryPath, "/cosmovisor") {
			args = append([]string{"run"}, args...)
		}

		if homePath != "" {
			args = append(args, "--home", homePath)
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
