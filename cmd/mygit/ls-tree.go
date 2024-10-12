package main

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

func LsTreeCmd() error {
	var hash string

	lsTreeCmd := flag.NewFlagSet(AvailableCommands.LsTree, flag.ExitOnError)
	nameOnly := lsTreeCmd.Bool("name-only", false, "show only names")

	err := lsTreeCmd.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	args := lsTreeCmd.Args()
	if len(args) > 0 {
		hash = args[0]
	}

	_, err = readTree(hash, *nameOnly)
	if err != nil {
		return err
	}

	return nil
}

func readTree(hash string, isNameOnly bool) (string, error) {
	pathToGitObject := ObjectPathBuilder(hash[:2], hash[2:])
	file, err := os.Open(pathToGitObject)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	r, err := zlib.NewReader(io.Reader(file))
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	defer r.Close()

	var contents bytes.Buffer
	_, err = io.Copy(&contents, r)
	if err != nil {
		return "", err
	}
	contents.ReadBytes('\x00')

	var objCollection []TreeObject

	for {
		var newObj TreeObject

		mode, err := contents.ReadBytes(' ')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		objectant, err := contents.ReadBytes('\x00')
		if err != nil {
			return "", err
		}

		sha := make([]byte, 20)
		_, err = contents.Read(sha)
		if err != nil {
			return "", err
		}

		newObj.Mode = string(bytes.TrimSpace(mode))
		newObj.Name = string(bytes.TrimRight(objectant, "\x00"))
		newObj.Hash = hex.EncodeToString(sha)

		if newObj.Mode == RAW_TREE_MODE {
			newObj.Mode = TREE_MODE
			newObj.Type = TREE_TYPE
		} else {
			newObj.Type = BLOB_TYPE
		}

		objCollection = append(objCollection, newObj)
	}

	for _, obj := range objCollection {
		switch {
		case isNameOnly:
			fmt.Println(obj.Name)

		default:
			fmt.Printf("%v %v %v %v\n", obj.Mode, obj.Type, obj.Hash, obj.Name)
		}
	}

	return "", nil
}
