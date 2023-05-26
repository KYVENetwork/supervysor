package pool

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"cosmossdk.io/log"
)

var logger = log.NewLogger(os.Stdout)

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
		logger.Error(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err.Error())
	}

	var resp Response
	err = json.Unmarshal([]byte(responseData), &resp)
	if err != nil {
		logger.Error(err.Error())
	}

	currentKey := resp.Pool.Data.CurrentKey
	poolHeight, err := strconv.Atoi(currentKey)
	if err != nil {
		logger.Error("couldn't convert poolHeight to str", err.Error())
	}

	return poolHeight
}
