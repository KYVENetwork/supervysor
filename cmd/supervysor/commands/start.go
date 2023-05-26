package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"supervysor/node"
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
