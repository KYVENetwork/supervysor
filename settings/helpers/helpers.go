package helpers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

// CalculateKeepRecent calculates the value for keepRecent, which is relevant for the pruning settings.
// It ensures that data that doesn't need to be stored anymore is pruned only after it has been validated.
// The calculation is based on the KYVE pool settings, and it ensures that blocks are stored for 3 days in advance.
func CalculateKeepRecent(maxBundleSize int, uploadInterval int) int {
	return int(
		math.Round(
			float64(maxBundleSize) / float64(uploadInterval) * 60 * 60 * 24 * 3))
}

// CalculateMaxDifference calculates a crucial threshold for the supervysor architecture.
// When the node is ahead of the pool by this value, the syncing process will halt in Ghost Mode.
// Once the node is again halfway within this value, the normal syncing process continues until
// reaching the threshold again.
func CalculateMaxDifference(maxBundleSize int, uploadInterval int) int {
	return int(
		math.Round(
			float64(maxBundleSize) / float64(uploadInterval) * 60 * 60 * 24 * 2))
}

// GetPoolSettings retrieves KYVE pool settings by using a list of endpoints (& optionally fallback endpoints)
// based on the provided chain and pool ID.
func GetPoolSettings(poolId int, chainId string, poolEndpoints string) ([2]int, error) {
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
			return [2]int{}, fmt.Errorf("unknown chainId")
		}
	}

	for _, endpoint := range endpoints {
		height, err := requestPoolSettings(poolId, endpoint)
		if err == nil {
			return height, err
		} else {
			fmt.Printf("failed to request pool settings from %v: %v", endpoint, err)
		}
	}

	return [2]int{}, fmt.Errorf("failed to get pool settings from all endpoints")
}

// SetPruningSettings updates the pruning settings in the app.toml file of the given home directory.
// It reads the current file, modifies the relevant lines and writes the updated lines back to the file.
func SetPruningSettings(homePath string, stateRequests bool, keepRecent int, interval int) error {
	configPath := filepath.Join(homePath, "config", "app.toml")

	file, err := os.OpenFile(configPath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var updatedLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line contains pruning settings and set new pruning settings
		if stateRequests {
			if strings.Contains(line, "pruning =") {
				line = "pruning = \"" + "custom" + "\""
			} else if strings.Contains(line, "pruning-keep-recent =") {
				line = "pruning-keep-recent = \"" + strconv.Itoa(keepRecent) + "\""
			} else if strings.Contains(line, "pruning-interval =") {
				line = "pruning-interval = \"" + strconv.Itoa(interval) + "\""
			} else if strings.Contains(line, "min-retain-blocks =") {
				line = "min-retain-blocks = 0"
			} else if strings.Contains(line, "snapshot-interval =") {
				line = "snapshot-interval = " + strconv.Itoa(0)
			}
		} else {
			if strings.Contains(line, "pruning =") {
				line = "pruning = \"" + "custom" + "\""
			} else if strings.Contains(line, "pruning-keep-recent =") {
				line = "pruning-keep-recent = \"" + "10" + "\""
			} else if strings.Contains(line, "pruning-interval =") {
				line = "pruning-interval = \"" + "100" + "\""
			} else if strings.Contains(line, "min-retain-blocks =") {
				line = "min-retain-blocks = 0"
			} else if strings.Contains(line, "snapshot-interval =") {
				line = "snapshot-interval = " + strconv.Itoa(0)
			}
		}
		updatedLines = append(updatedLines, line)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	if err = writeUpdatedConfig(configPath, updatedLines); err != nil {
		return err
	}

	return nil
}

// requestPoolSettings retrieves KYVE pool settings by making an GET request to the given endpoint.
// It reads the response, extracts the relevant settings information and returns it.
func requestPoolSettings(poolId int, endpoint string) ([2]int, error) {
	poolEndpoint := endpoint + "/kyve/query/v1beta1/pool/" + strconv.FormatInt(int64(poolId), 10)

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

// writeUpdatedConfig is responsible for writing the updated pruning settings to a given config file.
func writeUpdatedConfig(configPath string, pruningSettings []string) error {
	updatedFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer updatedFile.Close()

	writer := bufio.NewWriter(updatedFile)
	for _, line := range pruningSettings {
		if _, err = fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	return writer.Flush()
}
