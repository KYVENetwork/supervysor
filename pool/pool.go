package pool

import (
	"encoding/json"
	"fmt"
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

func GetPoolHeight(chainId string, poolId int64) (*int, error) {
	var poolEndpoint string
	if chainId == "korellia" {
		poolEndpoint = "https://api.korellia.kyve.network/kyve/query/v1beta1/pool/" + strconv.FormatInt(poolId, 10)
	} else if chainId == "kaon-1" {
		poolEndpoint = "https://api-eu-1.kaon.kyve.network/kyve/query/v1beta1/pool/" + strconv.FormatInt(poolId, 10)
	} else if chainId == "kyve-1" {
		poolEndpoint = "https://api-eu-1.kyve.network/kyve/query/v1beta1/pool/" + strconv.FormatInt(poolId, 10)
	} else {
		return nil, fmt.Errorf("unknown chainId (needs to be kyve-1, kaon-1 or korellia)")
	}
	response, err := http.Get(poolEndpoint)
	if err != nil {
		logger.Error("failed requesting KYVE endpoint", "err", err.Error())
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("failed reading KYVE endpoint response", "err", err.Error())
		return nil, err
	}

	var resp Response
	err = json.Unmarshal([]byte(responseData), &resp)
	if err != nil {
		logger.Error("failed unmarshalling KYVE endpoint response", "err", err.Error())
		return nil, err
	}

	currentKey := resp.Pool.Data.CurrentKey
	poolHeight, err := strconv.Atoi(currentKey)
	if err != nil {
		logger.Error("could not convert poolHeight to int", "err", err.Error())
		return nil, err
	}

	// TODO(@christopher): Remove for real testing environment
	poolHeight = 0

	return &poolHeight, err
}
