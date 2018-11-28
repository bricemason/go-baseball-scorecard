package main

import (
	"fmt"
	"regexp"
	"strings"
)

// PlayMatchConfig provides for a callback
type PlayMatchConfig struct {
	shortDescription string
	longDescription  string
}

var playMatchers = map[*regexp.Regexp]PlayMatchConfig{
	regexp.MustCompile("^K$"): PlayMatchConfig{
		shortDescription: "K",
		longDescription:  "strike out",
	},
	regexp.MustCompile("^(\\d)$"): PlayMatchConfig{
		shortDescription: "F%s",
		longDescription:  "fly ball caught by %s",
	},
	regexp.MustCompile("^(\\d)(\\d)$"): PlayMatchConfig{
		shortDescription: "%s-%s",
		longDescription:  "ground out %s-%s",
	},
}

func test() {
	fmt.Println("testing")
}

func processPlays(plays []string) {
	playRecords := GetRecords(plays)
	fmt.Println("there are", len(playRecords), "plays")

	for _, play := range playRecords {
		playString := play[6]

		// split up basic play and runner advancements
		playPieces := strings.Split(playString, ".")
		modifiers := strings.Split(playPieces[0], "/")[1:]

		fmt.Println("basic:", playPieces[0])
		fmt.Println("modifiers:", modifiers)

		// run the basic plays through the matcher
		for r, v := range playMatchers {
			playResult := r.FindStringSubmatch(playPieces[0])

			if len(playResult) > 0 {
				args := ToInterface(playResult[1:])
				fmt.Println(fmt.Sprintf(v.shortDescription, args...))
				fmt.Println(fmt.Sprintf(v.longDescription, args...))
			}
		}

		if len(playPieces) == 2 {
			fmt.Println("runner advancement:", playPieces[1])
		}

		fmt.Println("-----")

	}
}
