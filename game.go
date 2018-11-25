package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Game The top-level data structure to contain all game data
type Game struct {
	ID     string `json:"id"`
	Info   Info   `json:"info"`
	Lineup Lineup `json:"lineup"`
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
	allRecords := strings.Split(gameData, "\r\n")
	idSource := make([]string, 0)
	infoSource := make([]string, 0)
	startSource := make([]string, 0)
	playSource := make([]string, 0)

	for _, record := range allRecords {
		recordType := strings.Split(record, ",")[0]

		switch recordType {
		case "id":
			idSource = append(idSource, record)
		case "info":
			infoSource = append(infoSource, record)
		case "start":
			startSource = append(startSource, record)
		case "play":
			playSource = append(playSource, record)
		}
	}

	idRecords := GetRecords(idSource)
	info := CreateInfo(infoSource)
	lineup := CreateLineup(GetRecords(startSource), info.Usedh)

	game := Game{
		ID:     idRecords[0][1],
		Info:   info,
		Lineup: lineup,
	}

	return game
}
