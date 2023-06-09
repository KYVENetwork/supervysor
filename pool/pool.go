package pool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/KYVENetwork/supervysor/types"
)

func GetPoolHeight(chainId string, poolId int) (*int, error) {
	var poolEndpoint string
	if chainId == "korellia" {
		poolEndpoint = types.KorelliaEndpoint + strconv.FormatInt(int64(poolId), 10)
	} else if chainId == "kaon-1" {
		poolEndpoint = types.KaonEndpoint + strconv.FormatInt(int64(poolId), 10)
	} else if chainId == "kyve-1" {
		poolEndpoint = types.MainnetEndpoint + strconv.FormatInt(int64(poolId), 10)
	} else {
		return nil, fmt.Errorf("unknown chainId (needs to be kyve-1, kaon-1 or korellia)")
	}
	response, err := http.Get(poolEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed requesting KYVE endpoint: %s", err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading KYVE endpoint response: %s", err)
	}

	var resp types.SettingsResponse
	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling KYVE endpoint response: %s", err)
	}

	var poolHeight int
	currentKey := resp.Pool.Data.CurrentKey

	if currentKey == "" {
		startKey := resp.Pool.Data.StartKey
		poolHeight, err = strconv.Atoi(startKey)
		if err != nil {
			return nil, fmt.Errorf("could not convert poolHeight from start_key to int: %s", err)
		}
	} else {
		poolHeight, err = strconv.Atoi(currentKey)
		if err != nil {
			return nil, fmt.Errorf("could not convert poolHeight from current_key to int: %s", err)
		}
	}

	return &poolHeight, err
}
