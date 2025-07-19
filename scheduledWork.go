package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const TARGET_DIRECTORY string = "games"

type PlayerData struct {
	Name   string
	Scrore int
}

type GameData struct {
	Id       string
	Room     string
	Arda     bool
	Single   bool
	Created  string
	Time     int64
	Duration int64
	Players  []PlayerData
}

var g_MapGames map[string][]byte = make(map[string][]byte)

func processGames(list []GameData) {

	if list == nil {
		return
	}

	var size = len(list)
	if size < 1 {
		return
	}

	for _, a := range list {
		sJson, err := json.Marshal(a)
		if err != nil {
			log.Println(err)
			continue
		}

		g_MapGames[a.Id] = sJson
	}

	log.Println("Games processed", size)
}

func fetchActiveGames(url string) []GameData {

	log.Println("Fetching games from", url)

	response, err := http.Get(url)
	if err != nil {
		log.Println("Could not fetch data from url:", err)
		return nil
	}

	if response.StatusCode != 200 {
		log.Println("Unexpected status code:", response.StatusCode)
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Could not read data from body:", err)
		return nil
	}

	response.Body.Close()

	var gameData []GameData
	err = json.Unmarshal([]byte(body), &gameData)
	if err != nil {
		log.Println("Could not read data from body:", err)
		return nil
	}

	return gameData
}

func saveGame(id string, data []byte) bool {

	if data == nil || len(id) == 0 {
		return false
	}

	var file = TARGET_DIRECTORY + "/" + id + ".json"

	f, err := os.Create(file)
	if err != nil {
		log.Println("Could not create file", file)
		log.Println(err)
		return false
	}

	defer f.Close()

	_, errW := f.Write(data)
	if errW != nil {
		log.Println(errW)
		return false
	}

	log.Println("Game saved:", id)
	return true
}

func saveMap() {

	if len(g_MapGames) == 0 {
		log.Printf("No active games")
		return
	}

	ids := []string{}
	for id, data := range g_MapGames {
		if saveGame(id, data) {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		delete(g_MapGames, id)
		log.Println("Removed game from map", id)
	}
}

func scheduledWork(url string, waitInMinutes int32) {
	if len(url) == 0 {
		log.Println("NO URL GIVEN")
		return
	}

	if waitInMinutes < 1 {
		log.Panicln("Invalid waiting time, defaulting to 5mins")
		waitInMinutes = 5
	}

	if !createDirectory(TARGET_DIRECTORY) {
		log.Println("Cannot obtain output directory")
		return
	}

	log.Printf("Fetching games every %dmins", waitInMinutes)

	var count = 0
	var list []GameData
	var wait = time.Duration(waitInMinutes) * time.Minute

	for {

		list = fetchActiveGames(url)
		processGames(list)

		count++
		if count == 5 {
			count = 0
			saveMap()
		}

		time.Sleep(wait)
	}
}
