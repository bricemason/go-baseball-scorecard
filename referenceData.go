package main

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
