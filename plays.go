package main

import (
	"fmt"
	"regexp"
	"strings"
)

// PlayMatchConfig provides for a callback
type PlayMatchConfig struct {
	shortDescription string
	longDescription  func([]string) string
}

var baseMap = map[string]string{
	"1": "first",
	"2": "second",
	"3": "third",
	"H": "home",
}

var playMatchers = map[*regexp.Regexp]PlayMatchConfig{
	regexp.MustCompile("^K$"): PlayMatchConfig{
		shortDescription: "K",
		longDescription: func(matches []string) string {
			return "strike out"
		},
	},
	regexp.MustCompile("^(\\d)$"): PlayMatchConfig{
		shortDescription: "F%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("fly ball caught by %s", FieldingPositions[matches[0]].Code)
		},
	},
	regexp.MustCompile("^(\\d)(\\d)$"): PlayMatchConfig{
		shortDescription: "%s-%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("ground ball %s to %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code)
		},
	},
	regexp.MustCompile("^NP$"): PlayMatchConfig{
		shortDescription: "NP",
		longDescription: func(matches []string) string {
			return "no play"
		},
	},
	regexp.MustCompile("^WP$"): PlayMatchConfig{
		shortDescription: "WP",
		longDescription: func(matches []string) string {
			return "wild pitch"
		},
	},
	regexp.MustCompile("^SB(2|3|H)$"): PlayMatchConfig{
		shortDescription: "SB%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("stole %s", baseMap[matches[0]])
		},
	},
	regexp.MustCompile("^S$"): PlayMatchConfig{
		shortDescription: "S",
		longDescription: func(matches []string) string {
			return "single"
		},
	},
	regexp.MustCompile("^W$"): PlayMatchConfig{
		shortDescription: "W",
		longDescription: func(matches []string) string {
			return "walk"
		},
	},
	regexp.MustCompile("^HP$"): PlayMatchConfig{
		shortDescription: "HP",
		longDescription: func(matches []string) string {
			return "hit by pitch"
		},
	},
	regexp.MustCompile("^I|IW$"): PlayMatchConfig{
		shortDescription: "IW",
		longDescription: func(matches []string) string {
			return "intentional walk"
		},
	},
	regexp.MustCompile("^POCS(2|3|H)"): PlayMatchConfig{
		shortDescription: "POCS%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("picked off at %s (caught stealing)", baseMap[matches[0]])
		},
	},
	regexp.MustCompile("^(\\d)(\\d)\\((\\d)\\)$"): PlayMatchConfig{
		shortDescription: "DP%s%s(%s)",
		longDescription: func(matches []string) string {
			return "supposedly some sort of double play"
		},
	},
	regexp.MustCompile("^(\\d)(\\d)\\((\\d)\\)(\\d)$"): PlayMatchConfig{
		shortDescription: "DP%s%s(%s)%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("%s%s%s DP, runner at %s was the initial out", matches[0], matches[1], matches[3], baseMap[matches[2]])
		},
	},
	regexp.MustCompile("^PB$"): PlayMatchConfig{
		shortDescription: "PB",
		longDescription: func(matches []string) string {
			return "passed ball"
		},
	},
	regexp.MustCompile("^H|HR$"): PlayMatchConfig{
		shortDescription: "HR",
		longDescription: func(matches []string) string {
			return "home run"
		},
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
		basicPlayPieces := strings.Split(playPieces[0], "/")
		basicPlay := basicPlayPieces[0]
		modifiers := basicPlayPieces[1:]

		fmt.Println("basic:", basicPlay)
		fmt.Println("modifiers:", modifiers)

		// run the basic plays through the matcher
		for r, v := range playMatchers {
			playResult := r.FindStringSubmatch(basicPlay)

			if len(playResult) > 0 {
				args := ToInterface(playResult[1:])
				fmt.Println(fmt.Sprintf(v.shortDescription, args...))
				fmt.Println(v.longDescription(playResult[1:]))
				continue
			}
		}

		if len(playPieces) == 2 {
			fmt.Println("runner advancement:", playPieces[1])
		}

		fmt.Println("-----")

	}
}
