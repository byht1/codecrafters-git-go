package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path"
)

func WriteTreeCmd() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	hash, err := writeTree(currentDir)
	if err != nil {
		return err
	}

	fmt.Println(hash)

	return nil
}

func writeTree(pathToDir string) (string, error) {
	files, err := os.ReadDir(pathToDir)
	if err != nil {
		return "", err
	}

	var tree bytes.Buffer

	for _, file := range files {
		if file.Name() == ROOT_DIR {
			continue
		}

		newObj := TreeObject{
			Name: file.Name(),
		}

		processObject(pathToDir, file, &newObj)

		entry := fmt.Sprintf("%s %s\000", newObj.Mode, newObj.Name)
		tree.WriteString(entry)

		shaBytes, err := hex.DecodeString(newObj.Hash)
		if err != nil {
			return "", err
		}

		if len(shaBytes) != 20 {
			return "", fmt.Errorf("the SHA1 hash should be 20 bytes, but %d bytes were received", len(shaBytes))
		}
		tree.Write(shaBytes)
	}

	hashString, err := CreateTreeObject(tree.Bytes())
	if err != nil {
		return "", err
	}

	return hashString, nil
}

func processObject(pathToObject string, file fs.DirEntry, obj *TreeObject) error {
	pathToFile := path.Join(pathToObject, obj.Name)

	if file.IsDir() {
		return processDir(pathToFile, obj)
	}

	return processFile(pathToFile, obj)
}

func processDir(pathToFile string, obj *TreeObject) error {
	hash, err := writeTree(pathToFile)
	if err != nil {
		return err
	}

	obj.Hash = hash
	obj.Type = TREE_TYPE
	obj.Mode = RAW_TREE_MODE

	return nil
}

func processFile(pathToFile string, obj *TreeObject) error {
	info, err := os.Stat(pathToFile)
	if err != nil {
		return err
	}

	hash, err := CreateBlobObject(pathToFile)
	if err != nil {
		return err
	}

	obj.Mode = fmt.Sprintf("%v%o", PREFIX_MODE, info.Mode().Perm())
	obj.Type = BLOB_TYPE
	obj.Hash = hash

	return nil
}
