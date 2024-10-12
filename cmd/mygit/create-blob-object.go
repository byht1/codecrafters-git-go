package main

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
)

func createObject(objectType string, data []byte) (string, error) {
	header := fmt.Sprintf("%v %d\000", objectType, len(data))
	fullData := append([]byte(header), data...)

	hash := sha1.Sum(fullData)
	hashString := fmt.Sprintf("%x", hash)

	objectDir := ObjectPathBuilder(hashString[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return "", fmt.Errorf("error creating a directory: %v", err)
	}

	objectPath := ObjectPathBuilder(hashString[:2], hashString[2:])
	objectFile, err := os.Create(objectPath)
	if err != nil {
		return "", fmt.Errorf("error when creating an object: %v", err)
	}
	defer objectFile.Close()

	zlibWriter := zlib.NewWriter(objectFile)
	defer zlibWriter.Close()

	_, err = zlibWriter.Write(fullData)
	if err != nil {
		return "", fmt.Errorf("error when writing data: %v", err)
	}

	return hashString, nil
}

func CreateBlobObject(pathToFile string) (string, error) {
	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return "", fmt.Errorf("error reading a file: %v", err)
	}

	return createObject(BLOB_TYPE, data)
}

func CreateTreeObject(data []byte) (string, error) {
	return createObject(TREE_TYPE, data)
}

func CreateCommitObject(data []byte) (string, error) {
	return createObject(COMMIT_TYPE, data)
}
