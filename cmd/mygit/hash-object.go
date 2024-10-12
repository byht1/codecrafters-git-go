package main

import (
	"flag"
	"fmt"
	"os"
)

func HashObjectCmd() error {
	catFileCmd := flag.NewFlagSet(AvailableCommands.HashObject, flag.ExitOnError)
	w := catFileCmd.String("w", "default", "description")

	err := catFileCmd.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	hashString, err := CreateBlobObject(*w)
	if err != nil {
		return err
	}

	fmt.Print(hashString)

	return nil
}
