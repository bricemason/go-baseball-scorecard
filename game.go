package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Game The top-level data structure to contain all game data
type Game struct {
	ID       string `json:"id"`
	Date     string `json:"date"`
	GameType string `json:"gameType"`
	Team     Team   `json:"team"`
}

// Game.toJSON() Converts a Game struct to a json string
func (game *Game) toJSON() string {
	j, err := json.Marshal(game)
	if err != nil {
		panic("Error converting game to JSON")
	}

	return string(j)
}

// Game.toDisk() Converts a Game struct to a json string and writes to the given output path
func (game *Game) toDisk(path string) {
	jsonData := []byte(game.toJSON())
	outputPath := fmt.Sprintf("%s/%s.json", path, game.ID)

	err := ioutil.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		panic("Error occurred attempting to write game file with ID")
	}
}

// CreateGame Given a string of game data, construct and return a Game struct
func CreateGame(gameData string) Game {
	game := Game{}
	allRecords := strings.Split(gameData, "\r\n")
	idSource := make([]string, 0)
	infoSource := make([]string, 0)

	for _, record := range allRecords {
		recordType := strings.Split(record, ",")[0]

		if recordType == "id" {
			idSource = append(idSource, record)
		}

		if recordType == "info" {
			infoSource = append(infoSource, record)
		}
	}

	idRecords := GetRecords(idSource)
	infoRecords := CreateInfoRecords(infoSource)

	game.ID = idRecords[0][1]

	fmt.Println(infoRecords)

	return game
}
