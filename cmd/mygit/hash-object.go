package main

import (
	"flag"
	"fmt"
	"os"
)

func HashObjectCmd() error {
	hashObjectCmd := flag.NewFlagSet(AvailableCommands.HashObject, flag.ExitOnError)
	w := hashObjectCmd.String("w", PARAM_DEFAULT_VALUE, "description")

	err := hashObjectCmd.Parse(os.Args[2:])
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
