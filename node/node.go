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

var processId = 0

func InitialStart(stopChan chan struct{}) (*os.Process, error) {
	process, err := startNode(stopChan)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return process, err
	}

	fmt.Println("Process started: ", process.Pid)

	processId = process.Pid

	return process, nil
}

func GetNodeHeight() int {
	if processId == 0 {
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

	return nodeHeight
}
