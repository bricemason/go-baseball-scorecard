package main

// Info Represents all the game and administrative information tied to a game
type Info struct {
	Visteam string
}

// CreateInfoRecords Given a slice of raw comma-delimited strings, produce a slice of Info structs
func CreateInfoRecords(source []string) [][]string {
	infoRecords := GetRecords(source)

	// @TODO process fields in the info records and dynamically set in Info structs
	return infoRecords
}
