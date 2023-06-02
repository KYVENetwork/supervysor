package helpers

import (
	"os"
	"os/exec"

	"cosmossdk.io/log"
)

var logger = log.NewLogger(os.Stdout)

func MoveAddressBook(ghostMode bool) {
	if ghostMode {
		// Move address book to right place, because mode will change from Ghost to Normal
		source := "/root/.osmosisd/addrbook.json"
		destination := "/root/.osmosisd/config/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			logger.Error("could not move addrbook.json", "error", err)
			return
		}

		logger.Info("address book successfully moved back", "destination", destination)
	} else {
		// Move address book to hidden place, because mode will change from Normal to Ghost
		source := "/root/.osmosisd/config/addrbook.json"
		destination := "/root/.osmosisd/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			logger.Error("could not move addrbook.json", "err", err)
			return
		}

		logger.Info("address book moved", "destination", destination)
	}
}
