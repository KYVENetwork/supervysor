package main

import (
	cmd "github.com/KYVENetwork/supervysor/cmd/supervysor/commands"
	"github.com/KYVENetwork/supervysor/types"
)

var version string

func init() {
	types.Version = version
}

func main() {
	cmd.Execute()
}
