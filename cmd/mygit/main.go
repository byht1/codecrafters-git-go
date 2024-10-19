package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
	"github.com/codecrafters-io/git-starter-go/cmd/packages/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	availableCommands := constants.NewAvailableCommands()

	switch command := os.Args[1]; command {
	case availableCommands.Init:
		helpers.ProcessCmdFunc(InitCmd)
	case availableCommands.CatFile:
		helpers.ProcessCmdFunc(CatFileCmd)
	case availableCommands.HashObject:
		helpers.ProcessCmdFunc(HashObjectCmd)
	case availableCommands.LsTree:
		helpers.ProcessCmdFunc(LsTreeCmd)
	case availableCommands.WriteTree:
		helpers.ProcessCmdFunc(WriteTreeCmd)
	case availableCommands.CommitTree:
		helpers.ProcessCmdFunc(CommitTreeCmd)
	case availableCommands.Clone:
		helpers.ProcessCmdFunc(CloneCmd)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
