package pool

import (
	"encoding/json"
	"fmt"
	"github.com/KYVENetwork/supervysor/types"
	"io"
	"net/http"
	"strconv"
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

	currentKey := resp.Pool.Data.CurrentKey
	poolHeight, err := strconv.Atoi(currentKey)
	if err != nil {
		return nil, fmt.Errorf("could not convert poolHeight to int: %s", err)
	}

	return &poolHeight, err
}
