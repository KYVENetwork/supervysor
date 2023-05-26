package node

import (
	"fmt"
)

func EnableGhostMode() {
	// TODO: Check node status (Ghost or not)
	fmt.Println("Enabling Ghost Mode...")
	shutdownNode()

	// TODO: Change startNode to startGhostNode()
	process, err := startNode()
	if err != nil {
		fmt.Println("Ghost Mode enabling failed.")
	} else {
		if process != nil && process.Pid > 0 {
			fmt.Printf("Ghost Node started (PID: %d)\n", process.Pid)
		} else {
			fmt.Println("Ghost Mode enabling failed.")
		}
	}
}

func DisableGhostMode() {
	// TODO: Check node status (Ghost or not)
	fmt.Println("Disabling Ghost Mode...")

	shutdownNode()

	_, err := startNode()
	if err != nil {
		fmt.Println("Ghost Mode disabling failed.")
	} else {
		fmt.Println("Ghost Mode enabled.")
	}
	fmt.Println("Ghost Mode disabled.")
}
