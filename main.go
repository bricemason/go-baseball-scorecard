package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// League An MLB league
type League struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Team An MLB team
type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	League League `json:"league"`
}

// Game The top-level data structure to contain all game data
type Game struct {
	ID       string `json:"id"`
	Date     string `json:"date"`
	GameType string `json:"gameType"`
	Team     Team   `json:"team"`
}

// GameTypes Whether the game is a single game or part of a doubleheader
var GameTypes = map[int]string{
	0: "Single Game",
	1: "First Game of Doubleheader",
	2: "Second Game of Doubleheader",
}

// Leagues Reference data that contains all MLB leagues
var Leagues = map[string]League{
	"A": League{
		ID:   "A",
		Code: "AL",
		Name: "American League",
	},
	"N": League{
		ID:   "N",
		Code: "NL",
		Name: "National League",
	},
}

// Teams Reference data representing all MLB teams
// @TODO just use Mets for now, add more later
var Teams = map[string]Team{
	"NYN": Team{
		ID:     "NYN",
		Name:   "New York Mets",
		League: Leagues["N"],
	},
}

func main() {
	// hard-coding in 1986 Mets home games for now
	eventData, err := ioutil.ReadFile("data/in/1986NYN.EVN")
	if err != nil {
		fmt.Println("Error reading event file", err)
		return
	}

	// []string representing complete game data
	// note that we start by chopping off the id at the beginning of the event file to maintain consistency with subsequent game records in the file
	gamesData := strings.Split(string(eventData[2:]), "\nid")

	for _, gameData := range gamesData {
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
					recordData := fields[1]
					gameType, _ := strconv.Atoi(recordData[11:12])

					game.ID = recordData[:12]
					game.Team = Teams[recordData[:3]]
					game.Date = recordData[3:11]
					game.GameType = GameTypes[gameType]
				}
			}
		}

		// convert to json
		j, err := json.Marshal(game)
		if err != nil {
			fmt.Println("Error converting game to JSON:", err)
			return
		}

		// write out the resulting game file
		err = ioutil.WriteFile("data/out/"+game.ID+".json", j, 0644)
		if err != nil {
			fmt.Println("Error occurred attempting to write game file with ID", game.ID)
			return
		}
	}
}
