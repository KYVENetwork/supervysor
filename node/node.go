package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

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
	// TODO: Complete error handling

	if Process.Id == 0 {
		fmt.Println("Node hasn't started yet. Try again in 5s ...")
		time.Sleep(time.Second * 5)
		GetNodeHeight()
	}

	abciEndpoint := "http://localhost:26657/abci_info?"
	response, err := http.Get(abciEndpoint)

	if err != nil {
		fmt.Print(err.Error())
		fmt.Println("Failed to query height. Try again in 5s ...")
		time.Sleep(time.Second * 5)
		return GetNodeHeight()
	} else {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}

		var resp Response
		err = json.Unmarshal([]byte(responseData), &resp)
		if err != nil {
			fmt.Println("Error:", err)
		}

		lastBlockHeight := resp.Result.Response.LastBlockHeight

		nodeHeight, err := strconv.Atoi(lastBlockHeight)
		if err != nil {
			fmt.Println("Error during conversion", err)
		}

		return nodeHeight
	}
}

func startNode(initial bool) (*os.Process, error) {
	if !initial {
		moveAddressBook()
	}

	if !(Process.Id == 0 && Process.GhostMode) && !initial {
		// TODO: Panic and stop all processes
		return nil, nil
	} else {

		app := "osmosisd"
		arg1 := "start"
		arg2 := "--p2p.seeds"
		arg3 := "21d7539792ee2e0d650b199bf742c56ae0cf499e@162.55.132.230:2000,44ff091135ef2c69421eacfa136860472ac26e60@65.21.141.212:2000,ec4d3571bf709ab78df61716e47b5ac03d077a1a@65.108.43.26:2000,4cb8e1e089bdf44741b32638591944dc15b7cce3@65.108.73.18:2000,f515a8599b40f0e84dfad935ba414674ab11a668@osmosis.blockpane.com:26656,6bcdbcfd5d2c6ba58460f10dbcfde58278212833@osmosis.artifact-staking.io:26656,24841abfc8fbd401d8c86747eec375649a2e8a7e@osmosis.pbcups.org:26656,77bb5fb9b6964d6e861e91c1d55cf82b67d838b5@bd-osmosis-seed-mainnet-us-01.bdnodes.net:26656,3243426ab56b67f794fa60a79cc7f11bc7aa752d@bd-osmosis-seed-mainnet-eu-02.bdnodes.net:26656,ebc272824924ea1a27ea3183dd0b9ba713494f83@osmosis-mainnet-seed.autostake.com:26716,7c66126b64cd66bafd9ccfc721f068df451d31a3@osmosis-seed.sunshinevalidation.io:9393"

		cmdPath, err := exec.LookPath(app)
		if err != nil {
			return nil, err
		}

		cmd := exec.Command(cmdPath, arg1, arg2, arg3)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, fmt.Errorf("Couldn't start running the node.")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, err
		}

		return process, nil
	}
}

func startGhostNode() (*os.Process, error) {
	moveAddressBook()

	if !(Process.Id == 0 && !Process.GhostMode) {
		// TODO: Panic and stop all processes
		return nil, nil
	} else {

		app := "osmosisd"
		arg1 := "start"
		arg2 := "--p2p.seeds"
		arg3 := " "
		arg4 := "--p2p.laddr"

		// TODO: Find unused port
		arg5 := "tcp://0.0.0.0:26658"

		cmdPath, err := exec.LookPath(app)
		if err != nil {
			fmt.Println("Couldn't find /.osmosid")
			return nil, err
		}

		cmd := exec.Command(cmdPath, arg1, arg2, arg3, arg4, arg5)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, fmt.Errorf("Couldn't start running the node.")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, err
		}

		return process, nil
	}
}
