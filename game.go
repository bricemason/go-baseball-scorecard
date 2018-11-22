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
	records := strings.Split(gameData, "\n")

	for _, record := range records {
		// @NOTE I wonder how reliable this is. Surely there has to be a comma embedded somewhere in quotes in a field or something
		fields := strings.Split(record, ",")

		// a good record will have at least two fields, one for the type, the rest being the comma-delimited data
		if len(fields) >= 2 {
			recordType := fields[0]

			// this is an "id" record (we chopped this off earlier)
			if recordType == "" {
				game.ID = fields[1]
			}
		}
	}

	return game
}
