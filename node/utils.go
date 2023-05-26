package node

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func startNode() (*os.Process, error) {

	moveAddressBook()

	if !(Process.Id == 0 && *Process.GhostMode == true) {
		// TODO: Panic and stop all processes
		return nil, nil
	} else {

		app := "osmosisd"
		arg1 := "start"
		arg2 := "--p2p.seeds"
		arg3 := "21d7539792ee2e0d650b199bf742c56ae0cf499e@162.55.132.230:2000,44ff091135ef2c69421eacfa136860472ac26e60@65.21.141.212:2000,ec4d3571bf709ab78df61716e47b5ac03d077a1a@65.108.43.26:2000,4cb8e1e089bdf44741b32638591944dc15b7cce3@65.108.73.18:2000,f515a8599b40f0e84dfad935ba414674ab11a668@osmosis.blockpane.com:26656,6bcdbcfd5d2c6ba58460f10dbcfde58278212833@osmosis.artifact-staking.io:26656,24841abfc8fbd401d8c86747eec375649a2e8a7e@osmosis.pbcups.org:26656,77bb5fb9b6964d6e861e91c1d55cf82b67d838b5@bd-osmosis-seed-mainnet-us-01.bdnodes.net:26656,3243426ab56b67f794fa60a79cc7f11bc7aa752d@bd-osmosis-seed-mainnet-eu-02.bdnodes.net:26656,ebc272824924ea1a27ea3183dd0b9ba713494f83@osmosis-mainnet-seed.autostake.com:26716,7c66126b64cd66bafd9ccfc721f068df451d31a3@osmosis-seed.sunshinevalidation.io:9393"
		cmdPath, err := exec.LookPath(app)
		if err != nil {
			return nil, err
		}

		cmd := exec.Command(cmdPath, arg1, arg2, arg3)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, fmt.Errorf("Couldn't start running the node.")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, err
		}

		return process, nil
	}
}

func startGhostNode() (*os.Process, error) {

	moveAddressBook()

	if !(Process.Id == 0 && *Process.GhostMode == false) {
		// TODO: Panic and stop all processes
		return nil, nil
	} else {

		app := "osmosisd"
		arg1 := "start"
		arg2 := "--p2p.seeds"
		arg3 := " "
		arg4 := "--p2p.laddr"

		// TODO: Find unused port
		arg5 := "tcp://0.0.0.0:26658"

		cmdPath, err := exec.LookPath(app)
		if err != nil {
			fmt.Println("Couldn't find /.osmosid")
			return nil, err
		}

		cmd := exec.Command(cmdPath, arg1, arg2, arg3, arg4, arg5)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		processIDChan := make(chan int)

		go func() {
			err := cmd.Start()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
				return
			}

			processIDChan <- cmd.Process.Pid

			// Wait for process end
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err)
				processIDChan <- -1
			}
		}()

		processID := <-processIDChan

		if processID == -1 {
			return nil, fmt.Errorf("Couldn't start running the node.")
		}

		process, err := os.FindProcess(processID)
		if err != nil {
			return nil, err
		}

		return process, nil
	}
}

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
	if *Process.GhostMode {
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
