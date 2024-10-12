package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"
)

func CommitTreeCmd() error {
	commitTree := flag.NewFlagSet(AvailableCommands.HashObject, flag.ExitOnError)
	p := commitTree.String("p", PARAM_DEFAULT_VALUE, "Each -p indicates the id of a parent commit object")
	m := commitTree.String("m", PARAM_DEFAULT_VALUE, "A paragraph in the commit log message. This can be given more than once and each <message> becomes its own paragraph")

	err := commitTree.Parse(os.Args[3:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	if *p == PARAM_DEFAULT_VALUE {
		p = nil
	}

	if *m == PARAM_DEFAULT_VALUE {
		return fmt.Errorf("the -m option is optional")
	}

	now := time.Now()
	_, offset := now.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60

	obj := CommitObject{
		Tree:    os.Args[2],
		Parent:  p,
		Message: *m,
		Email:   "example@gmail.com",
		Author:  "IHdPA",
		Date:    fmt.Sprintf("%v %+03d%02d", now.Unix(), offsetHours, offsetMinutes),
	}

	var commit bytes.Buffer

	commit.WriteString(fmt.Sprintf("%v %v\n", TREE_TYPE, obj.Tree))
	if obj.Parent != nil {
		commit.WriteString(fmt.Sprintf("parent %v\n", *obj.Parent))
	}
	commit.WriteString(fmt.Sprintf("author %v <%v> %v\n", obj.Author, obj.Email, obj.Date))
	commit.WriteString(fmt.Sprintf("committer %v <%v> %v\n", obj.Author, obj.Email, obj.Date))
	commit.WriteString(fmt.Sprintln(""))
	commit.WriteString(fmt.Sprintln(obj.Message))

	hashString, err := CreateCommitObject(commit.Bytes())
	if err != nil {
		return err
	}

	fmt.Println(hashString)

	return nil
}
