package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/KYVENetwork/supervysor/types"

	"cosmossdk.io/log"
)

var logger = log.NewLogger(os.Stdout)

func CheckBinaryPath(path string) error {
	_, err := exec.LookPath(path)
	if err != nil {
		return err
	}
	return nil
}

func CheckFilePath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

func CalculateKeepRecent(maxBundleSize int, uploadInterval int) int {
	return int(
		math.Round(
			float64(maxBundleSize) / float64(uploadInterval) * 60 * 60 * 24 * 2))
}

func CalculateMaxDifference(maxBundleSize int, uploadInterval int) int {
	return int(
		math.Round(
			float64(maxBundleSize) / float64(uploadInterval) * 60 * 60 * 24 * 1))
}

func GetPoolSettings(poolId int, chainId string, fallbackEndpoints string) ([2]int, error) {
	var endpoints []string
	var err error

	if chainId == "korellia" {
		endpoints = types.KorelliaEndpoints
	} else if chainId == "kaon-1" {
		endpoints = types.KaonEndpoints
	} else if chainId == "kyve-1" {
		endpoints = types.MainnetEndpoints
	} else {
		return [2]int{}, fmt.Errorf("unknown chainId")
	}

	for _, endpoint := range append(endpoints, strings.Split(fallbackEndpoints, ",")...) {
		if endpoint != "" {
			if height, err := requestPoolSettings(poolId, endpoint); err == nil {
				return height, err
			}
		}
	}

	return [2]int{}, err
}

func requestPoolSettings(poolId int, endpoint string) ([2]int, error) {
	poolEndpoint := endpoint + "/kyve/query/v1beta1/pool/" + strconv.FormatInt(int64(poolId), 10)

	fmt.Println(poolEndpoint)

	response, err := http.Get(poolEndpoint)
	if err != nil {
		logger.Error("API is not available", err.Error())
		return [2]int{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("got unexpected response", err.Error())
		return [2]int{}, err
	}

	var resp types.SettingsResponse
	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		logger.Error("couldn't unmarshal response", err.Error())
		return [2]int{}, err
	}

	uploadInterval := resp.Pool.Data.UploadInterval
	interval, err := strconv.Atoi(uploadInterval)
	if err != nil {
		logger.Error("couldn't convert uploadInterval to int", err.Error())
		return [2]int{}, err
	}

	maxBundleSize := resp.Pool.Data.MaxBundleSize
	size, err := strconv.Atoi(maxBundleSize)
	if err != nil {
		logger.Error("couldn't convert maxBundleSize to int", err.Error())
		return [2]int{}, err
	}

	return [2]int{size, interval}, nil
}
