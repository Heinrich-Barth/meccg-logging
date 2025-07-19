package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func createDirectory(name string) bool {

	if len(name) == 0 || strings.Contains(name, ".") {
		log.Println("Invalid name")
		return false
	}

	var path = filepath.Join(".", "log")

	_, err := os.Stat(path)
	if err == nil {
		log.Println("Directory already available", name)
		return true
	}

	if !os.IsNotExist(err) {
		log.Println(err)
		return false
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Println("Cannot create log directory", name)
		return false
	}

	return true
}
