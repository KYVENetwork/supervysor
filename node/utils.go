package node

import (
	"fmt"
	"os"
	"os/exec"
)

func startNode() (*os.Process, error) {
	// TODO: Check if process.id is still running
	// TODO: Move filled address book, expose seeds

	app := "osmosisd"
	arg1 := "start"
	arg2 := "--p2p.seeds"

	// TODO: Add exposed seeds from cmd input
	arg3 := ""

	cmdPath, err := exec.LookPath(app)

	cmd := exec.Command(cmdPath, arg1, arg2, arg3)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return cmd.Process, nil
}

func startGhostNode() {

	// TODO: move addressbook, change laddr port, expose no seeds

	app := "osmosisd"
	arg1 := "start"
	arg2 := "--p2p.seeds"
	arg3 := "_"
	arg4 := "--p2p.laddr"
	arg5 := "tcp://0.0.0.0:26658"

	cmd := exec.Command(app, arg1, arg2, arg3, arg4, arg5)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func shutdownNode() {
	// TODO: Expect process.id and shutdown the process
}
