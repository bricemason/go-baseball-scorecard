package main

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
