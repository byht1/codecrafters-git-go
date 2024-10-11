package main

import (
	"fmt"
	"os"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	NewAvailableCommands()

	switch command := os.Args[1]; command {
	case AvailableCommands.Init:
		ProcessCmdFunc(InitCmd)
	case AvailableCommands.CatFile:
		ProcessCmdFunc(CatFileCmd)
	case AvailableCommands.HashObject:
		ProcessCmdFunc(HashObjectCmd)
	case AvailableCommands.LsTree:
		ProcessCmdFunc(LsTreeCmd)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
