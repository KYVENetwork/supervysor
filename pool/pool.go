package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KYVENetwork/supervysor/types"
)

// GetPoolHeight retrieves the KYVE pool height by using a list of endpoints (& optionally fallback endpoints)
// based on the provided chain and pool ID.
func GetPoolHeight(chainId string, poolId int, poolEndpoints string) (int, error) {
	var endpoints []string

	if poolEndpoints != "" {
		endpoints = strings.Split(poolEndpoints, ",")
	} else {
		if chainId == "korellia-2" {
			endpoints = types.KorelliaEndpoints
		} else if chainId == "kaon-1" {
			endpoints = types.KaonEndpoints
		} else if chainId == "kyve-1" {
			endpoints = types.MainnetEndpoints
		} else {
			return 0, fmt.Errorf("unknown chainId")
		}
	}

	for i := 0; i <= types.BackoffMaxRetries; i++ {
		delay := time.Duration(math.Pow(2, float64(i))) * time.Second

		for _, endpoint := range endpoints {
			height, err := requestPoolHeight(poolId, endpoint)
			if err == nil {
				return height, nil
			} else {
				fmt.Printf("failed to request pool height from %v: %v\n", endpoint, err)
			}
		}

		if i <= types.BackoffMaxRetries {
			fmt.Printf("retrying to query pool again in %v\n", delay)
			time.Sleep(delay)
		}
	}

	return 0, fmt.Errorf("failed to get pool height from all endpoints")
}

// requestPoolHeight retrieves KYVE pool height by making an GET request to the given endpoint.
func requestPoolHeight(poolId int, endpoint string) (int, error) {
	poolEndpoint := endpoint + "/kyve/query/v1beta1/pool/" + strconv.FormatInt(int64(poolId), 10)

	response, err := http.Get(poolEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed requesting KYVE endpoint: %s", err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("failed reading KYVE endpoint response: %s", err)
	}

	var resp types.SettingsResponse
	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		return 0, fmt.Errorf("failed unmarshalling KYVE endpoint response: %s", err)
	}

	var poolHeight int
	currentKey := resp.Pool.Data.CurrentKey

	if currentKey == "" {
		startKey := resp.Pool.Data.StartKey
		poolHeight, err = strconv.Atoi(startKey)
		if err != nil {
			return 0, fmt.Errorf("could not convert poolHeight from start_key to int: %s", err)
		}
	} else {
		poolHeight, err = strconv.Atoi(currentKey)
		if err != nil {
			return 0, fmt.Errorf("could not convert poolHeight from current_key to int: %s", err)
		}
	}

	return poolHeight, err
}
