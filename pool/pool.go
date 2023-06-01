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
			CurrentKey     string `json:"current_key"`
			UploadInterval string `json:"upload_interval"`
			MaxBundleSize  string `json:"max_bundle_size"`
		} `json:"data"`
	} `json:"pool"`
}

func GetPoolHeight(poolId int64) (*int, error) {
	poolEndpoint := "https://api.korellia.kyve.network/kyve/query/v1beta1/pool/" + strconv.FormatInt(poolId, 10)
	response, err := http.Get(poolEndpoint)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var resp Response
	err = json.Unmarshal([]byte(responseData), &resp)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	currentKey := resp.Pool.Data.CurrentKey
	poolHeight, err := strconv.Atoi(currentKey)
	if err != nil {
		logger.Error("couldn't convert poolHeight to int", err.Error())
		return nil, err
	}

	return &poolHeight, err
}
