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
	hash, err := writeTree("./")
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

	var treeObj []TreeObject

	for _, file := range files {
		if file.Name() == ROOT_DIR {
			continue
		}

		newObj := TreeObject{
			Name: file.Name(),
		}

		processObject(pathToDir, file, &newObj)

		treeObj = append(treeObj, newObj)
	}

	var tree bytes.Buffer

	for _, obj := range treeObj {
		entry := fmt.Sprintf("%s %s\000", obj.Mode, obj.Name)
		tree.WriteString(entry)

		shaBytes, err := hex.DecodeString(obj.Hash)
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
