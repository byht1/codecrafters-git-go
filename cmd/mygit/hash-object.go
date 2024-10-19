package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
	gitObject "github.com/codecrafters-io/git-starter-go/cmd/packages/git-object"
)

func HashObjectCmd() error {
	hashObjectCmd := flag.NewFlagSet(constants.AvailableCommands.HashObject, flag.ExitOnError)
	w := hashObjectCmd.String("w", constants.PARAM_DEFAULT_VALUE, "description")

	err := hashObjectCmd.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	_, hashString, err := gitObject.CreateBlobObject(*w)
	if err != nil {
		return err
	}

	fmt.Print(hashString)

	return nil
}
