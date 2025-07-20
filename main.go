package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func createEmptyEnv() {

	f, err := os.Create(".env")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, errW := f.WriteString("URL=")
	if errW != nil {
		log.Fatal(errW)
	}

}

func main() {

	setupLogger()

	defer func() {
		log.Println("Closing log file")
		logFile.Close()
	}()

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file present.")
		createEmptyEnv()
		log.Fatalf("Error loading .env file: %s", err)
	}

	URL := os.Getenv("URL")
	go scheduledWork(URL, 5)
	go initServer()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Received termination signal, shutting down...")

	/* try to save unsaved games */
	saveMap()
	closeLogger()
}
