package main

import (
	"os"
	"os/signal"
	"supervysor/cmd/supervysor/commands"
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
	}()
}
