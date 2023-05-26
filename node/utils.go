package node

import (
	"os"
	"os/exec"
	"syscall"
)

func shutdownNode() {
	process, err := os.FindProcess(Process.Id)
	if err != nil {
		logger.Error("couldn't find process to shutdown")
		// TODO(@christopher): Panic and shutdown all running processes
	}

	// Terminate the process
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// TODO(@christopher): Panic and shutdown all running processes
		logger.Error("couldn't terminate process: %s\n", err)
		return
	}

	Process.Id = 0

	logger.Info("process terminated successfully")
}

func moveAddressBook() {
	if Process.GhostMode {
		// Move address book to right place, because mode will change from Ghost to Normal
		source := "/root/.osmosisd/addrbook.json"
		destination := "/root/.osmosisd/config/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			logger.Error("couldn't move addrbook.json: %s\n", err)
			return
		}

		logger.Info("address book successfully moved back to %s .", destination)
	} else {
		// Move address book to hidden place, because mode will change from Normal to Ghost
		source := "/root/.osmosisd/config/addrbook.json"
		destination := "/root/.osmosisd/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			logger.Error("couldn't move addrbook.json: %s\n", err)
			return
		}

		logger.Info("address book successfully moved to %s .", destination)
	}
}
