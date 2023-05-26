package node

import (
	"fmt"
	"os"
)

func InitialStart() (int, error) {
	process, err := startNode()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return 0, err
	}

	fmt.Println("Process started: ", process.Pid)

	Process.Id = process.Pid
	*Process.GhostMode = false

	return process.Pid, nil
}

func EnableGhostMode() {
	// TODO: Check node status (Ghost or not)
	fmt.Println("Enabling Ghost Mode...")
	shutdownNode()

	// TODO: Change startNode to startGhostNode()
	process, err := startGhostNode()
	if err != nil {
		fmt.Println("Ghost Mode enabling failed.")
	} else {
		if process != nil && process.Pid > 0 {
			Process.Id = process.Pid
			*Process.GhostMode = true
			fmt.Printf("Ghost Node started (PID: %d)\n", process.Pid)
		} else {
			// TODO: Panic and shutdown all processes
			fmt.Println("Ghost Mode enabling failed.")
		}
	}
}

func DisableGhostMode() {
	// TODO: Check node status (Ghost or not)
	fmt.Println("Disabling Ghost Mode...")

	shutdownNode()

	//_, err := startNode()
	//if err != nil {
	//	fmt.Println("Ghost Mode disabling failed.")
	//} else {
	//	fmt.Println("Ghost Mode enabled.")
	//}
	//fmt.Println("Ghost Mode disabled.")
}
