package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const outputPath string = "data/out"

func main() {
	// hard-coding in 1986 Mets home games for now
	eventData, err := ioutil.ReadFile("data/in/1986NYN.EVN")
	if err != nil {
		fmt.Println("Error reading event file", err)
		return
	}

	// []string representing complete game data
	// note that we start by chopping off the id at the beginning of the event file, but we'll put it back together when processing
	gamesData := strings.Split(string(eventData[2:]), "\nid")

	for _, gameData := range gamesData {
		// make the first "id" record whole again
		game := CreateGame("id" + gameData)

		// dump to terminal to check data as we go
		fmt.Println(game.toJSON())

		// dump to disk as well
		game.toDisk(outputPath)

	}
}
