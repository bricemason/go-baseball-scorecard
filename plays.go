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
		143
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
	// this is weird but going with it for now
	// originally the single number was going to represent a fly ball out but this can and is usually accompanied by
	// modifiers that can have it be things like an unassisted ground out or even the seemingly undocumented play of 6(1)
	// we may just end up having to pass in the modifiers here too
	regexp.MustCompile("^(\\d)(?:\\(\\d\\))?$"): PlayMatchConfig{
		shortDescription: "O%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("out %s", FieldingPositions[matches[0]].Code)
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
	regexp.MustCompile("^SB(2|3|H);SB(2|3|H)$"): PlayMatchConfig{
		shortDescription: "DSB %s,%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("double steal %s, %s", baseMap[matches[0]], baseMap[matches[1]])
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
	regexp.MustCompile("^PO(1|2|3)\\(\\d\\d\\)"): PlayMatchConfig{
		shortDescription: "PO%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("picked off at %s", baseMap[matches[0]])
		},
	},
	regexp.MustCompile("^PO(1|2|3)\\(E(\\d)\\)"): PlayMatchConfig{
		shortDescription: "PO%s-E%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("pick off attempt at %s, error by %s", baseMap[matches[0]], FieldingPositions[matches[1]].Code)
		},
	},
	regexp.MustCompile("^(\\d)(\\d)\\((\\d)\\)$"): PlayMatchConfig{
		shortDescription: "DP%s%s(%s)",
		longDescription: func(matches []string) string {
			return "supposedly some sort of double play"
		},
	},
	regexp.MustCompile("^(\\d)(\\d)(\\d)$"): PlayMatchConfig{
		shortDescription: "DP-%s%s%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("double play? %s-%s-%s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code, FieldingPositions[matches[2]].Code)
		},
	},
	regexp.MustCompile("^(\\d)\\(B\\)(\\d)\\((\\d)\\)$"): PlayMatchConfig{
		shortDescription: "LDP-%s-%s(%s)",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("lined into double play %s to %s at %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code, baseMap[matches[2]])
		},
	},
	regexp.MustCompile("^(\\d)\\(\\d\\)(\\d)$"): PlayMatchConfig{
		shortDescription: "UGDP-%s-%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("ground ball double play, unassisted %s to %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code)
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
	regexp.MustCompile("^BK$"): PlayMatchConfig{
		shortDescription: "BK",
		longDescription: func(matches []string) string {
			return "balk"
		},
	},
	regexp.MustCompile("^FLE(\\d)$"): PlayMatchConfig{
		shortDescription: "FLE%s",
		longDescription: func(matches []string) string {
			return fmt.Sprintf("error by %s on foul fly ball", FieldingPositions[matches[0]].Code)
		},
	},
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

		fmt.Println("full:", playString)
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
