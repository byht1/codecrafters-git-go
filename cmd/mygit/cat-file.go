package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
	zlibHelper "github.com/codecrafters-io/git-starter-go/cmd/packages/zlib-helper"
)

func CatFileCmd() error {
	catFileCmd := flag.NewFlagSet(constants.AvailableCommands.CatFile, flag.ExitOnError)
	p := catFileCmd.String("p", constants.PARAM_DEFAULT_VALUE, "description")

	err := catFileCmd.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	return catFile(*p)
}

func catFile(value string) error {
	pathToFile := constants.ObjectPathBuilder(value[0:2], value[2:])
	file, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	byteData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	_, object, err := zlibHelper.ReadObject(byteData)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	parts := strings.Split(string(object), "\x00")
	fmt.Print(parts[1])

	return nil
}
