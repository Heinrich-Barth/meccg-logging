package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func getLogFile() string {

	if createDirectory("log") {
		return "./log/application.log"
	}

	return ""
}

var logFile *os.File = nil

func closeLogger() {
	if logFile != nil {
		fmt.Println("Closing log file")

		logFile.Close()
		logFile = nil
	}
}

func setupLogger() {

	LOG_FILE := getLogFile()
	if len(LOG_FILE) == 0 {
		log.Println("Cannot obtain logfile")
	}

	logFile, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Cannot write to logfile at", LOG_FILE)
		log.Panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
