package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func MoveAddressBook(activateGhostMode bool, addrBookPath string) error {
	if activateGhostMode {
		// TODO(@christopher): Check if addrbook needs to be moved
		parentDir := filepath.Dir(filepath.Dir(addrBookPath))
		filename := filepath.Base(addrBookPath)
		destPath := filepath.Join(parentDir, filename)

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
			return fmt.Errorf("could not delete source address book file: %s", err)
		}
	} else {
		// TODO(@christopher): Check if addrbook needs to be moved
		parentDir := filepath.Dir(filepath.Dir(addrBookPath))
		filename := filepath.Base(addrBookPath)
		sourcePath := filepath.Join(parentDir, filename)

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
			return fmt.Errorf("could not delete source address book file: %s", err)
		}
	}

	return nil
}
