package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

var g_MapGames map[string]GameData = make(map[string]GameData)

func processGames(list []GameData) {

	if list == nil {
		return
	}

	var size = len(list)
	if size < 1 {
		return
	}

	for _, a := range list {
		_, err := json.Marshal(a)
		if err != nil {
			log.Println(err)
			continue
		}

		g_MapGames[a.Id] = a
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

func getDateInfo(date string) string {

	var parts = strings.Split(date, " ")
	if len(parts) < 4 {
		return ""
	}

	return parts[3] + "-" + parts[2] + "-" + parts[1]
}

func saveGame(game GameData) bool {

	var filename = getDateInfo(game.Created) + "-" + game.Id + ".json"
	var file = TARGET_DIRECTORY + "/" + filename

	f, err := os.Create(file)
	if err != nil {
		log.Println("Could not create file", file)
		log.Println(err)
		return false
	}

	defer f.Close()

	data, err := json.Marshal(&game)
	if err != nil {
		log.Println(err)
		return false
	}

	_, errW := f.Write(data)
	if errW != nil {
		log.Println(errW)
		return false
	}

	log.Println("Game saved:", filename)
	return true
}

func saveMap() {

	if len(g_MapGames) == 0 {
		log.Printf("No active games")
		return
	}

	ids := []string{}
	for id, data := range g_MapGames {
		if saveGame(data) {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		delete(g_MapGames, id)
		log.Println("Removed game from map", id)
	}
}

func listGameFiles() []string {

	files, err := os.ReadDir(TARGET_DIRECTORY)
	if err != nil {
		log.Println(err)
		return nil
	}

	var list []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			list = append(list, file.Name())
		}
	}

	return list
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

	var list []GameData
	var wait = time.Duration(waitInMinutes) * time.Minute

	for {

		list = fetchActiveGames(url)
		processGames(list)
		saveMap()

		time.Sleep(wait)
	}
}
