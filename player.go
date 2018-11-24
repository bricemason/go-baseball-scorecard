package main

import (
	"strconv"
)

// FieldingPosition represents a player position in a game
type FieldingPosition struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Player represents a player in a lineup
type Player struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	BattingPosition  int              `json:"battingPosition"`
	FieldingPosition FieldingPosition `json:"fieldingPosition"`
}

// Lineup The starting lineup for teams in a game
type Lineup struct {
	Visitor []Player `json:"visitor"`
	Home    []Player `json:"home"`
}

// CreateLineup Splits up the raw "start" record types and creates starting lineups for each team
func CreateLineup(startRecords [][]string, useDH bool) Lineup {
	// by default, we'll assume there's no DH so we're working with 18 records (9 each team)
	splitPos := 8
	visitor := make([]Player, 0)
	home := make([]Player, 0)

	if useDH {
		// if the game uses a DH, then we have 20 records (10 each team)
		// in this case, the pitcher is listed as batting in position 0
		splitPos = 9
	}

	for i, entry := range startRecords {
		battingPosition, _ := strconv.Atoi(entry[4])
		fieldingPosition, _ := strconv.Atoi(entry[5])
		player := Player{
			ID:               entry[1],
			Name:             entry[2],
			BattingPosition:  battingPosition,
			FieldingPosition: FieldingPositions[fieldingPosition],
		}

		if i <= splitPos {
			// visitor team looks to be chunked at the top of the record set
			visitor = append(visitor, player)
		} else {
			home = append(home, player)
		}
	}

	return Lineup{
		Visitor: visitor,
		Home:    home,
	}
}
