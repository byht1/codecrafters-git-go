package gitObject

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
)

type HashSumType [20]byte

func CreateObject(objectType string, data []byte) (HashSumType, string, error) {
	header := fmt.Sprintf("%v %d\000", objectType, len(data))
	fullData := append([]byte(header), data...)

	hash := sha1.Sum(fullData)
	hashString := fmt.Sprintf("%x", hash)

	objectDir := constants.ObjectPathBuilder(hashString[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return hash, "", fmt.Errorf("error creating a directory: %v", err)
	}

	objectPath := constants.ObjectPathBuilder(hashString[:2], hashString[2:])
	objectFile, err := os.Create(objectPath)
	if err != nil {
		return hash, "", fmt.Errorf("error when creating an object: %v", err)
	}
	defer objectFile.Close()

	zlibWriter := zlib.NewWriter(objectFile)
	defer zlibWriter.Close()

	_, err = zlibWriter.Write(fullData)
	if err != nil {
		return hash, "", fmt.Errorf("error when writing data: %v", err)
	}

	return hash, hashString, nil
}

func CreateBlobObject(pathToFile string) (HashSumType, string, error) {
	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return [20]byte{}, "", fmt.Errorf("error reading a file: %v", err)
	}

	return CreateObject(BLOB_TYPE, data)
}

func CreateTreeObject(data []byte) (HashSumType, string, error) {
	return CreateObject(TREE_TYPE, data)
}

func CreateCommitObject(data []byte) (HashSumType, string, error) {
	return CreateObject(COMMIT_TYPE, data)
}
