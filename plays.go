package main

import (
	"fmt"
	"regexp"
	"strings"
)

var playMatchers = map[string]string{
	"^(\\d)(\\d)$": "ground out from %s to %s",
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
		for k, v := range playMatchers {
			r := regexp.MustCompile(k)
			result := r.FindStringSubmatch(playPieces[0])

			if len(result) > 0 {
				// @TODO we want to use the results slice to feed into Sprintf but we need []interface{} to spread into that
				// we'll probably need a simple util to run that conversion per the golang faqs
				fmt.Println("FOUND results on play", playPieces[0], "with", v)
			}
		}

		if len(playPieces) == 2 {
			fmt.Println("runner advancement:", playPieces[1])
		}

		fmt.Println("-----")

	}
}
