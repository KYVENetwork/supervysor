package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

var ProcessId = 0

func InitialStart() (int, error) {
	process, err := startNode()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return 0, err
	}

	fmt.Println("Process started: ", process.Pid)

	ProcessId = process.Pid

	return process.Pid, nil
}

func GetNodeHeight() int {
	if ProcessId == 0 {
		fmt.Println("Node hasn't started yet.")
		time.Sleep(time.Second * 5)
		GetNodeHeight()
	}
	// TODO: Query from locally running node (-> not configurable)
	abciEndpoint := "http://localhost:26657/abci_info?"
	response, err := http.Get(abciEndpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

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

	fmt.Println("GetNodeHeight: ", nodeHeight)

	return nodeHeight
}
