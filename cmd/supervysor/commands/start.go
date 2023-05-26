package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"supervysor/node"
	"syscall"
	"time"
)

var (
	poolId int64
)

func init() {
	startCmd.Flags().Int64Var(&poolId, "pool-id", 0, "pool id")
	if err := startCmd.MarkFlagRequired("pool-id"); err != nil {
		panic(fmt.Errorf("flag 'pool-id' should be required: %w", err))
	}

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start supervising node",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := node.InitialStart()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Println("STARTED INITIALLY.")

		// Fire off a goroutine to loop until that channel receives a signal.
		// When a signal is received simply exit the program
		go func() {
			// Setup a channel to receive a signal
			c := make(chan os.Signal, 1)

			// Notify this channel when a SIGINT is received
			signal.Notify(c, os.Interrupt)

			<-c

			fmt.Println("CTRL + C")

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
			os.Exit(0)
		}()

		for {
			nodeHeight := node.GetNodeHeight()

			// poolHeight := pool.GetPoolHeight(poolId)
			poolHeight := 1

			fmt.Println(nodeHeight, poolHeight)

			diff := nodeHeight - poolHeight

			if diff >= 10000 {
				node.EnableGhostMode()
			} else if diff < 10000 && diff > 5000 {
				// do nothing
			} else if diff <= 5000 {
				node.DisableGhostMode()
			} else {
				fmt.Println("Error: negative difference between pool and node height.")
			}

			time.Sleep(time.Minute)
		}
	},
}
