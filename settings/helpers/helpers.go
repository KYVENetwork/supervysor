package helpers

import (
	"github.com/KYVENetwork/supervysor/settings"
	"math"
	"os/exec"
)

func CheckBinaryPath(path string) error {
	_, err := exec.LookPath(path)
	if err != nil {
		return err
	}
	return nil
}

func CalculateKeepRecent(poolSettings settings.PoolSettingsType) int {
	return int(
		math.Round(
			float64(poolSettings.MaxBundleSize) / float64(poolSettings.UploadInterval) * 60 * 60 * 24 * 7))
}
