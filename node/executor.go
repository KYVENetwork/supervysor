package node

import (
	"fmt"
	"os"
)

func InitialStart() (int, error) {
	fmt.Println("Starting initially...")
	process, err := startNode(true)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return 0, err
	}

	fmt.Println("Initial process started: ", process.Pid)

	Process.Id = process.Pid
	Process.GhostMode = false

	return process.Pid, nil
}

func EnableGhostMode() {
	if !Process.GhostMode {
		fmt.Println("Enabling Ghost Mode...")
		shutdownNode()

		process, err := startGhostNode()
		if err != nil {
			fmt.Println("Ghost Mode enabling failed.")
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				fmt.Printf("Ghost Node started (PID: %d)\n", process.Pid)
			} else {
				// TODO: Panic and shutdown all processes
				fmt.Println("Ghost Mode enabling failed.")
			}
		}
	} else {
		fmt.Println("Keeping Ghost Mode enabled...")
	}
}

func DisableGhostMode() {
	if Process.GhostMode {
		fmt.Println("Disabling Ghost Mode...")

		shutdownNode()

		process, err := startNode(false)
		if err != nil {
			fmt.Println("Ghost Mode disabling failed.")
		} else {
			if process != nil && process.Pid > 0 {
				Process.Id = process.Pid
				Process.GhostMode = true
				fmt.Printf("Normal Node started (PID: %d)\n", process.Pid)
			} else {
				// TODO: Panic and shutdown all processes
				fmt.Println("Ghost Mode disabling failed.")
			}
		}
	} else {
		fmt.Println("Keeping Normal Mode enabled...")
	}
}
