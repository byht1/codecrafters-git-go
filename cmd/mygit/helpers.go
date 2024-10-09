package main

import "log"

func ProcessCmdFunc(fn func() error) {
	err := fn()
	if err != nil {
		log.Fatalln(err)
	}
}
