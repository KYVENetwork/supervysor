package helpers

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"cosmossdk.io/log"
)

// GetPort resolves an unused TCP address.
func GetPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err = l.Close()
		if err != nil {
			panic(err)
		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}

// MoveAddressBook is responsible for moving an address book file from one location to another,
// making it not visible in the GhostMode and visible in NormalMode.
func MoveAddressBook(activateGhostMode bool, addrBookPath string, log log.Logger) error {
	if activateGhostMode {
		parentDir := filepath.Dir(filepath.Dir(addrBookPath))
		filename := filepath.Base(addrBookPath)
		destPath := filepath.Join(parentDir, filename)

		if _, err := os.Stat(destPath); err == nil {
			if _, err = os.Stat(addrBookPath); err == nil {
				err = os.Remove(addrBookPath)
				if err != nil {
					return fmt.Errorf("could not remove address book file: %s", err)
				}
			}
			return nil
		}

		srcFile, err := os.Open(addrBookPath)
		if err != nil {
			return fmt.Errorf("could not open source address book file: %s", err)
		}
		defer func(srcFile *os.File) {
			err = srcFile.Close()
			if err != nil {
				log.Error("could not close src file")
				panic("could not close file")
			}
		}(srcFile)

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("could not create new address book file: %s", err)
		}
		defer func(destFile *os.File) {
			err = destFile.Close()
			if err != nil {
				log.Error("could not close dest file")
				panic("could not close file")
			}
		}(destFile)

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("could not copy address book into new address book: %s", err)
		}

		err = os.Remove(addrBookPath)
		if err != nil {
			return fmt.Errorf("could not remove source address book file: %s", err)
		}
	} else {
		parentDir := filepath.Dir(filepath.Dir(addrBookPath))
		filename := filepath.Base(addrBookPath)
		sourcePath := filepath.Join(parentDir, filename)

		if _, err := os.Stat(addrBookPath); err == nil {
			if _, err = os.Stat(sourcePath); err == nil {
				err = os.Remove(sourcePath)
				if err != nil {
					return fmt.Errorf("could not remove address book file: %s", err)
				}
			}
			return nil
		}

		srcFile, err := os.Open(sourcePath)
		if err != nil {
			return fmt.Errorf("could not open source address book file: %s", err)
		}
		defer func(srcFile *os.File) {
			err = srcFile.Close()
			if err != nil {
				log.Error("could not close srcFile")
				panic("could not close file")
			}
		}(srcFile)

		destFile, err := os.Create(addrBookPath)
		if err != nil {
			return fmt.Errorf("could not create new address book file: %s", err)
		}
		defer func(destFile *os.File) {
			err = destFile.Close()
			if err != nil {
				log.Error("could not close destFile")
				panic("could not close file")
			}
		}(destFile)

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("could not copy address book into new address book: %s", err)
		}

		err = os.Remove(sourcePath)
		if err != nil {
			return fmt.Errorf("could not remove source address book file: %s", err)
		}
	}

	return nil
}

func SplitArgs(argsString string) []string {
	// Split the string by spaces
	split := strings.Fields(argsString)

	var args []string
	var currentArg string

	for _, part := range split {
		// If the current part starts with "--" or "-", consider it as a new argument
		if strings.HasPrefix(part, "--") || strings.HasPrefix(part, "-") {
			// If there was a previous argument, add it to the result
			if currentArg != "" {
				args = append(args, currentArg)
			}

			// Start a new argument
			currentArg = part
		} else {
			// If the current part doesn't start with "--" or "-", append it to the current argument
			currentArg += " " + part
		}
	}

	// Add the last argument to the result
	if currentArg != "" {
		args = append(args, currentArg)
	}

	return args
}
