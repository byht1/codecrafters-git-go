package main

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func CatFileCmd() error {
	catFileCmd := flag.NewFlagSet(AvailableCommands.CatFile, flag.ExitOnError)
	p := catFileCmd.String("p", "default", "description")

	err := catFileCmd.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	return catFile(*p)
}

func catFile(value string) error {
	pathToFile := path.Join(".git/objects", value[0:2], value[2:])
	file, err := os.Open(pathToFile)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	r, err := zlib.NewReader(io.Reader(file))
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	defer r.Close()

	s, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	parts := strings.Split(string(s), "\x00")
	fmt.Print(parts[1])

	return nil
}
