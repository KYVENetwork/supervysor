package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Pool struct {
		Data struct {
			CurrentKey string `json:"current_key"`
		} `json:"data"`
	} `json:"pool"`
}

func GetPoolHeight(poolId int64) int {
	poolEndpoint := "https://api.korellia.kyve.network/kyve/query/v1beta1/pool/" + strconv.FormatInt(poolId, 10)
	response, err := http.Get(poolEndpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var resp Response
	err = json.Unmarshal([]byte(responseData), &resp)
	if err != nil {
		fmt.Println("Error:", err)
	}

	currentKey := resp.Pool.Data.CurrentKey

	poolHeight, err := strconv.Atoi(currentKey)

	if err != nil {
		fmt.Println("Error during conversion", err)
	}

	return poolHeight
}
