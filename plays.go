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

/*
	Records with no docs:
		6(1)
*/

var playMatchers = map[*regexp.Regexp]PlayMatchConfig{
	// I've seen some records with "K23" so we're keeping this matcher loose and just looking for a "K" at the beginning
	// I'm assuming the "23" in this case would be a ball in the dirt and a throw to first by the catcher for the out
	regexp.MustCompile("^K"): PlayMatchConfig{
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
	regexp.MustCompile("^S(\\d)?$"): PlayMatchConfig{
		shortDescription: "S%s",
		longDescription: func(matches []string) string {
			if matches[0] != "" {
				return fmt.Sprintf("single to %s", FieldingPositions[matches[0]].Code)
			}

			return "single"
		},
	},
	regexp.MustCompile("^D(\\d)?$"): PlayMatchConfig{
		shortDescription: "D%s",
		longDescription: func(matches []string) string {
			if matches[0] != "" {
				return fmt.Sprintf("double to %s", FieldingPositions[matches[0]].Code)
			}

			return "double"
		},
	},
	regexp.MustCompile("^T(\\d)?$"): PlayMatchConfig{
		shortDescription: "T%s",
		longDescription: func(matches []string) string {
			if matches[0] != "" {
				return fmt.Sprintf("triple to %s", FieldingPositions[matches[0]].Code)
			}

			return "triple"
		},
	},
	regexp.MustCompile("^FC(\\d)?$"): PlayMatchConfig{
		shortDescription: "FC%s",
		longDescription: func(matches []string) string {
			if matches[0] != "" {
				return fmt.Sprintf("fielder's choice to %s", FieldingPositions[matches[0]].Code)
			}

			return "fielder's choice"
		},
	},
	regexp.MustCompile("^E(\\d)?$"): PlayMatchConfig{
		shortDescription: "E%s",
		longDescription: func(matches []string) string {
			if matches[0] != "" {
				return fmt.Sprintf("error by %s", FieldingPositions[matches[0]].Code)
			}

			return "error"
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
	// The two caught stealing situations below have fielding scenarios after them  to note where the throw was from and to whom
	// but I'm not interested in logging that
	regexp.MustCompile("^CS(2|3|H)"): PlayMatchConfig{
		shortDescription: "CS%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("caught stealing %s", baseMap[matches[0]])
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

		hasPlayConfig := false

		// run the basic plays through the matcher
		for r, v := range playMatchers {
			playResult := r.FindStringSubmatch(basicPlay)

			if len(playResult) > 0 {
				args := ToInterface(playResult[1:])
				fmt.Println(fmt.Sprintf(v.shortDescription, args...))
				fmt.Println(v.longDescription(playResult[1:]))
				hasPlayConfig = true
				break
			}
		}

		if hasPlayConfig == false {
			panic(fmt.Sprintf("no play config found for %s", basicPlay))
		}

		if len(playPieces) == 2 {
			fmt.Println("runner advancement:", playPieces[1])
		}

		fmt.Println("-----")

	}
}
