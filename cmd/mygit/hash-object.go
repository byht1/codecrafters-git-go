package main

import (
	"compress/zlib"
	"crypto/sha1"
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

	return hashObject(*w)
}

func hashObject(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading a file: %v", err)
	}

	header := fmt.Sprintf("blob %d\000", len(data))
	fullData := append([]byte(header), data...)

	hash := sha1.Sum(fullData)
	hashString := fmt.Sprintf("%x", hash)

	objectDir := ObjectPathBuilder(hashString[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return fmt.Errorf("error creating a directory: %v", err)
	}

	objectPath := ObjectPathBuilder(hashString[:2], hashString[2:])
	objectFile, err := os.Create(objectPath)
	if err != nil {
		return fmt.Errorf("error when creating an object: %v", err)
	}
	defer objectFile.Close()

	zlibWriter := zlib.NewWriter(objectFile)
	defer zlibWriter.Close()

	_, err = zlibWriter.Write(fullData)
	if err != nil {
		return fmt.Errorf("error when writing data: %v", err)
	}

	// fmt.Println("The Git object is saved as:", hashString)
	fmt.Print(hashString)

	return nil
}
