package main

import (
	"fmt"
	"os"
	"os/signal"
	"supervysor/cmd/supervysor/commands"
	"supervysor/node"
	"syscall"
)

func main() {
	commands.Execute()

	// Setup a channel to receive a signal
	done := make(chan os.Signal, 1)

	// Notify this channel when a SIGINT is received
	signal.Notify(done, os.Interrupt)

	// Fire off a goroutine to loop until that channel receives a signal.
	// When a signal is received simply exit the program
	go func() {
		for _ = range done {
			os.Exit(0)
		}

		if node.ProcessId != 0 {
			process, err := os.FindProcess(node.ProcessId)
			if err != nil {
				fmt.Printf("Fehler beim Finden des Prozesses: %v\n", err)
				os.Exit(1)
			}

			err = process.Signal(syscall.SIGTERM)
			if err != nil {
				fmt.Printf("Fehler beim Beenden des Prozesses: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Prozess erfolgreich beendet.")
		}
	}()
}
