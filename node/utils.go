package node

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func shutdownNode() {
	process, err := os.FindProcess(Process.Id)
	if err != nil {
		fmt.Println("Error finding process.")
		// TODO: Panic and shutdown all running processes
	}

	// Terminate the process
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// TODO: Panic and shutdown all running processes
		fmt.Printf("Error terminating process: %s\n", err)
		return
	}

	Process.Id = 0

	fmt.Println("Process terminated successfully.")
}

func moveAddressBook() {
	if Process.GhostMode {
		// Move address book to right place, because mode will change from Ghost to Normal
		source := "/root/.osmosisd/addrbook.json"
		destination := "/root/.osmosisd/config/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Couldn't move addrbook.json: %s\n", exitError.Stderr)
			} else {
				fmt.Printf("Couldn't move addrbook.json: %s\n", err.Error())
			}
			return
		}
		fmt.Printf("Addressbook sucessfully moved back to %s .", destination)
	} else {
		// Move address book to hidden place, because mode will change from Normal to Ghost
		source := "/root/.osmosisd/config/addrbook.json"
		destination := "/root/.osmosisd/ "

		cmd := exec.Command("mv", source, destination)

		err := cmd.Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Couldn't move addrbook.json: %s\n", exitError.Stderr)
			} else {
				fmt.Printf("Couldn't move addrbook.json: %s\n", err.Error())
			}
			return
		}
		fmt.Printf("Addressbook sucessfully moved to %s .", destination)
	}
}
