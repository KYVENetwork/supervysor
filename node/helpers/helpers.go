package helpers

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"cosmossdk.io/log"

	"github.com/rs/zerolog"
)

var logger = log.NewLogger(os.Stdout)

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
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// InitLogger initializes the file logging for the current log file.
func InitLogger(logFile string) log.Logger {
	File, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		panic(err)
	}

	multiLogger := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, File)

	logger = log.NewCustomLogger(zerolog.New(multiLogger).With().Timestamp().Logger())

	return logger
}

// MoveAddressBook is responsible for moving an address book file from one location to another,
// making it not visible in the GhostMode and visible in NormalMode.
func MoveAddressBook(activateGhostMode bool, addrBookPath string) error {
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
			err := srcFile.Close()
			if err != nil {
				panic("could not close file")
			}
		}(srcFile)

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("could not create new address book file: %s", err)
		}
		defer func(destFile *os.File) {
			err := destFile.Close()
			if err != nil {
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
			err := srcFile.Close()
			if err != nil {
				panic("could not close file")
			}
		}(srcFile)

		destFile, err := os.Create(addrBookPath)
		if err != nil {
			return fmt.Errorf("could not create new address book file: %s", err)
		}
		defer func(destFile *os.File) {
			err := destFile.Close()
			if err != nil {
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
