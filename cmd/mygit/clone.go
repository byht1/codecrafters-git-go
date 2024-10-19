package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
	gitObject "github.com/codecrafters-io/git-starter-go/cmd/packages/git-object"
	"github.com/codecrafters-io/git-starter-go/cmd/packages/helpers"
	zlibHelper "github.com/codecrafters-io/git-starter-go/cmd/packages/zlib-helper"
)

type DeltaObject struct {
	baseObject string
	data       []byte
}

type CloneArgs struct {
	Repo string
	Dir  string
}

func readUint32BigEndian(bytes []byte) uint32 {
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}

func readPktLine(blob []byte) (int, []byte, error) {
	pktLength := blob[:4]
	blob = blob[4:]

	dst := [2]byte{}
	_, err := hex.Decode(dst[:], pktLength)
	if err != nil {
		return 0, nil, err
	}

	size := uint16(dst[0])<<8 | uint16(dst[1])

	// pkt-line 0000
	if size == 0 {
		return 4, []byte{}, nil
	}

	if len(blob) < int(size)-4 {
		return 4, nil, errors.New("error reading pkt line")
	}

	data := blob[:size-4]

	// strip trailing linefeed, if it exists
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	return int(size), data, nil
}

func getObjectName(pktLines [][]byte) (string, error) {
	// skip the first pktLine (001e# service=git-upload-pack)
	for _, pktLine := range pktLines[1:] {
		if len(pktLine) == 0 {
			continue
		}

		var hash, ref string
		fmt.Sscanf(string(pktLine), "%s %s", &hash, &ref)
		if ref == "refs/heads/master" || ref == "refs/heads/main" {
			return hash, nil
		}
	}

	return "", errors.New("invalid pktLines")
}

func getPackFile(cloneUrl string) ([]byte, string, error) {
	response, err := http.Get(fmt.Sprintf("%s/info/refs?service=git-upload-pack", cloneUrl))
	if err != nil {
		return nil, "", err
	}

	discoveryBuffer := bytes.Buffer{}
	io.Copy(&discoveryBuffer, response.Body)
	discovery := discoveryBuffer.Bytes()

	pktLines := [][]byte{}

	for len(discovery) > 0 {
		n, data, err := readPktLine(discovery)
		if err != nil {
			return nil, "", err
		}

		discovery = discovery[n:]
		pktLines = append(pktLines, data)
	}

	objectName, err := getObjectName(pktLines)
	if err != nil {
		return nil, "", err
	}

	buffer := bytes.NewBufferString(fmt.Sprintf("0032want %s\n00000009done\n", objectName))
	response, err = http.Post(fmt.Sprintf("%s/git-upload-pack", cloneUrl), "application/x-git-upload-pack-request", buffer)
	if err != nil {
		return nil, "", err
	}

	packFileBuffer := bytes.Buffer{}
	io.Copy(&packFileBuffer, response.Body)
	packFile := packFileBuffer.Bytes()
	n, _, err := readPktLine(packFile) // read 0008NAK
	if err != nil {
		return nil, "", err
	}

	packFile = packFile[n:]
	return packFile, objectName, nil
}

func readSize(packFile []byte) (size uint64, used int, err error) {
	data := packFile[used]
	used++

	size = uint64(data & 0x7F)
	shift := 7

	for data&0x80 != 0 {
		if len(packFile) <= used || 64 <= shift {
			return 0, 0, errors.New("bad size")
		}

		data = packFile[used]
		used++

		size += uint64(data&0x7F) << shift
		shift += 7
	}
	return size, used, nil
}

func checkoutPackFile(packFile []byte) error {
	if len(packFile) < 32 {
		return errors.New("bad pack file")
	}

	checksum := packFile[len(packFile)-20:]
	packFile = packFile[:len(packFile)-20]
	expected := sha1.Sum(packFile)

	if !bytes.Equal(checksum, expected[:]) {
		return errors.New("invalid pack file checksum")
	}

	if !bytes.Equal(packFile[0:4], []byte("PACK")) {
		return errors.New("invalid pack file header")
	}

	version := readUint32BigEndian(packFile[4:8])

	if version != 2 && version != 3 {
		return errors.New("invalid pack file version")
	}

	return nil
}

func processGitObject(packFile []byte, used *int, size uint64, objectType gitObject.ObjectCode) error {
	read, object, err := zlibHelper.ReadObject(packFile[*used:])
	*used += read
	if err != nil {
		return err
	}

	if int(size) != len(object) {
		return errors.New("bad object header length")
	}

	objectTypeStr := map[gitObject.ObjectCode]string{
		gitObject.OBJ_COMMIT_CODE: gitObject.COMMIT_TYPE,
		gitObject.OBJ_TREE_CODE:   gitObject.TREE_TYPE,
		gitObject.OBJ_BLOB_CODE:   gitObject.BLOB_TYPE,
		gitObject.OBJ_TAG_CODE:    gitObject.TAG_TYPE,
	}[objectType]

	_, _, err = gitObject.CreateObject(objectTypeStr, object)
	if err != nil {
		return err
	}

	return nil
}

func processOfcDelta(packFile []byte, used *int, size uint64) error {
	_, read, err := readSize(packFile[*used:])
	*used += read
	if err != nil {
		return err
	}

	read, object, err := zlibHelper.ReadObject(packFile[*used:])
	*used += read
	if err != nil {
		return err
	}

	if int(size) != len(object) {
		return errors.New("bad object header length")
	}

	return errors.New("cant handle ofs-delta object")
}

func processRefDelta(packFile []byte, used *int, size uint64, deltaObjects *[]DeltaObject) error {
	hash := packFile[*used : *used+20]
	*used += 20

	read, object, err := zlibHelper.ReadObject(packFile[*used:])
	*used += read

	if err != nil {
		return err
	}

	if int(size) != len(object) {
		return errors.New("bad object header length")
	}

	*deltaObjects = append(*deltaObjects, DeltaObject{baseObject: hex.EncodeToString(hash), data: object})

	return nil
}

func writePackFile(packFile []byte) error {
	err := checkoutPackFile(packFile)
	if err != nil {
		return err
	}

	used := 8
	numObjects := readUint32BigEndian(packFile[used:])
	used += 4

	var objectsRead uint32
	deltaObjects := []DeltaObject{}
	packFile = packFile[:len(packFile)-20]

	for used < len(packFile) {
		objectsRead++

		size, objectType, read, err := gitObject.ReadObjectHeader(packFile[used:])
		used += read

		if err != nil {
			return err
		}

		switch objectType {
		case gitObject.OBJ_COMMIT_CODE, gitObject.OBJ_TREE_CODE, gitObject.OBJ_BLOB_CODE, gitObject.OBJ_TAG_CODE:
			err := processGitObject(packFile, &used, size, objectType)
			if err != nil {
				return err
			}
		case gitObject.OBJ_OFS_DELTA_CODE:
			err := processOfcDelta(packFile, &used, size)
			if err != nil {
				return err
			}
		case gitObject.OBJ_REF_DELTA_CODE:
			err := processRefDelta(packFile, &used, size, &deltaObjects)
			if err != nil {
				return err
			}
		default:
			return errors.New("invalid object type")
		}

	}

	if numObjects != objectsRead {
		return errors.New("bad object count")
	}

	for len(deltaObjects) > 0 {
		unAddedDeltaObjects := []DeltaObject{}
		added := false

		for _, delta := range deltaObjects {
			if objectExists(delta.baseObject) {
				added = true
				baseObject, objectType, err := gitObject.OpenObject(delta.baseObject)
				if err != nil {
					return err
				}

				err = writeDeltaObject(baseObject, delta.data, objectType)
				if err != nil {
					return err
				}

			} else {
				unAddedDeltaObjects = append(unAddedDeltaObjects, delta)
			}
		}

		if !added {
			return errors.New("bad delta objects")
		}

		deltaObjects = unAddedDeltaObjects
	}

	return nil
}

func writeDeltaObject(baseObject, deltaObject []byte, objectType string) error {
	used := 0
	baseSize, read, err := readSize(deltaObject[used:])
	if err != nil {
		return err
	}
	used += read

	if len(baseObject) != int(baseSize) {
		return errors.New("bad delta header")
	}

	expectedSize, read, err := readSize(deltaObject[used:])
	if err != nil {
		return err
	}
	used += read

	buffer := bytes.Buffer{}

	for used < len(deltaObject) {
		opcode := deltaObject[used]
		used++

		if opcode&0x80 != 0 {
			var argument uint64

			for bit := 0; bit < 7; bit++ {
				if opcode&(1<<bit) != 0 {
					argument += uint64(deltaObject[used]) << (bit * 8)
					used++
				}
			}

			offset := argument & 0xFFFFFFFF
			size := (argument >> 32) & 0xFFFFFF

			if size == 0 {
				size = 0x10000
			}

			buffer.Write(baseObject[offset : offset+size])

		} else {
			size := int(opcode & 0x7F)
			buffer.Write(deltaObject[used : used+size])
			used += size
		}
	}

	undeltifiedObject := buffer.Bytes()

	if int(expectedSize) != len(undeltifiedObject) {
		return errors.New("bad delta header")
	}

	_, _, err = gitObject.CreateObject(objectType, undeltifiedObject)
	if err != nil {
		return err
	}

	return nil
}

func objectExists(hash string) bool {
	path := constants.ObjectPathBuilder(hash[:2], hash[2:])
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func checkoutCommit(commitHash string) error {
	commit, objectType, err := gitObject.OpenObject(commitHash)
	if err != nil {
		return err
	}

	if objectType != gitObject.COMMIT_TYPE {
		return errors.New("object not a commit")
	}
	treeHash := commit[5:45]

	err = checkoutTree(string(treeHash), ".")
	return err
}

func checkoutTree(treeHash, dir string) error {
	os.MkdirAll(dir, 0755)

	tree, objectType, err := gitObject.OpenObject(treeHash)
	if err != nil {
		return err
	}

	if objectType != gitObject.TREE_TYPE {
		return errors.New("object not a tree")
	}

	for len(tree) > 0 {
		used, mode, name, hash := helpers.ReadTreeEntry(tree)
		tree = tree[used:]

		hashStr := hex.EncodeToString(hash[:])
		fullPath := fmt.Sprintf("%s/%s", dir, name)

		if mode == gitObject.RAW_TREE_MODE {
			err = checkoutTree(hashStr, fullPath)
			if err != nil {
				return err
			}
		}

		if mode == gitObject.FILE_MODE_1 || mode == gitObject.FILE_MODE_2 {
			blob, objectType, err := gitObject.OpenObject(hashStr)
			if err != nil {
				return err
			}

			if objectType != gitObject.BLOB_TYPE {
				return errors.New("object not a blob")
			}

			os.WriteFile(fullPath, blob, 0644) // currently ignoring mode
		}
	}

	return nil
}

func checkoutArgs() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: mygit <repo> <your_dir>\n")
		os.Exit(1)
	}
}

func getArgs() CloneArgs {
	args := CloneArgs{
		Repo: os.Args[2],
		Dir:  "test",
	}

	if len(os.Args) == 4 {
		args.Dir = os.Args[3]
	}

	return args
}

func createDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}

func CloneCmd() error {
	checkoutArgs()
	args := getArgs()

	err := createDir(args.Dir)
	if err != nil {
		return err
	}

	InitCmd()

	packFile, commit, err := getPackFile(args.Repo)
	if err != nil {
		return err
	}
	err = writePackFile(packFile)
	if err != nil {
		return err
	}

	err = checkoutCommit(commit)
	if err != nil {
		return err
	}

	return nil
}
