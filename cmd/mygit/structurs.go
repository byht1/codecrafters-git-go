package main

type TreeObject struct {
	Mode string
	Type string
	Hash string
	Name string
}

type CommitObject struct {
	Tree    string
	Parent  *string
	Author  string
	Email   string
	Date    string
	Message string
}
