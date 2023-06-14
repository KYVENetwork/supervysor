package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSupervysorDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %s", err)
	}

	return filepath.Join(home, ".supervysor"), nil
}
