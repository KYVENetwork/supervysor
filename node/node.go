package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	GhostMode *bool
}

var Process = ProcessType{
	Id:        0,
	GhostMode: nil,
}

func GetNodeHeight() int {
	// TODO: Complete error handling

	if Process.Id == 0 {
		fmt.Println("Node hasn't started yet. Try again in 5s ...")
		time.Sleep(time.Second * 5)
		GetNodeHeight()
	}

	fmt.Println("Got ProcessId ", Process.Id, "; start getNodeHeight().")
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

		fmt.Println("Succesfully retrieved node height: ", nodeHeight)

		return nodeHeight
	}
}
