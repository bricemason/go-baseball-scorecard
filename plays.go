package main

import (
	"fmt"
	"regexp"
	"strings"
)

// RunnerAdvancement represents how baserunners are impacted from a play
type RunnerAdvancement struct {
	fromBase string
	toBase   string
	isOut    bool
}

// PlayConfig represents the chunks in an entry in a game event file
type PlayConfig struct {
	source             string
	basicPlay          string
	modifiers          []string
	runnerAdvancements []RunnerAdvancement
	code               string
	description        string
}

// PlayCreator is a factor for Plays
type PlayCreator func(playConfig PlayConfig, matches []string) PlayConfig

func createPlayConfig(playSource string) PlayConfig {
	// playSource is: S9/L9S.2-H;1-3
	play := strings.Split(playSource, ".")     // ["S9/L9S", "2-H;1-3"]
	batterEvent := strings.Split(play[0], "/") // ["S9", "L9S"]
	config := PlayConfig{
		source:    playSource,
		basicPlay: batterEvent[0],  // S9
		modifiers: batterEvent[1:], // ["L9S"]
	}

	if len(play) == 2 {
		runners := make([]RunnerAdvancement, 0)

		for _, v := range strings.Split(play[1], ";") { // 2-H;1-3
			r := regexp.MustCompile("^(1|2|3|H)(\\-|X)(1|2|3|H)?")
			runnerMatch := r.FindStringSubmatch(v)

			if len(runnerMatch) == 4 {
				runners = append(runners, RunnerAdvancement{
					fromBase: runnerMatch[1],
					toBase:   runnerMatch[3],
					isOut:    runnerMatch[2] == "X",
				})
			}
		}

		config.runnerAdvancements = runners
	}

	return config
}

/*
	Records with no docs:
		6(1)
		143
		366(1)
*/

var uniqueValues = make([]string, 0)

var playMatchers = map[*regexp.Regexp]PlayCreator{
	regexp.MustCompile("^K$"): func(playConfig PlayConfig, matches []string) PlayConfig {
		playConfig.code = "K"
		playConfig.description = "strike out"

		return playConfig
	},
	/*
		basic plays for this group: 3 2 4 6 5 1
		modifiers for this play: G FL P L F BP SH F1 BG
		modifiers groups for this play: G FL P L F BP,FL SH BP F1 BG FL,P
	*/
	regexp.MustCompile("^(\\d)$"): func(playConfig PlayConfig, matches []string) PlayConfig {
		switch position := matches[0]; position {
		case "7", "8", "9":
			playConfig.code = "F" + position
			playConfig.description = fmt.Sprintf("fly ball out %s", FieldingPositions[position].Code)
		default:
			if len(playConfig.modifiers) > 0 {
				playConfig.code = position
			}
		}

		playConfig.description = fmt.Sprintf("%s %s", playConfig.description, TranslateModifiers(playConfig.modifiers))

		return playConfig
	},
	/*
		basic plays for this group: 41 43 53 13 63 31 23 34 14 24 54
	*/
	regexp.MustCompile("^(\\d)(\\d)$"): func(playConfig PlayConfig, matches []string) PlayConfig {
		playConfig.code = fmt.Sprintf("%s-%s", matches[0], matches[1])
		playConfig.description = fmt.Sprintf("ground ball %s to %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code)

		return playConfig
	},
	// 	// regexp.MustCompile("^NP$"): PlayMatchConfig{
	// 	// 	shortDescription: "NP",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "no play"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^WP$"): PlayMatchConfig{
	// 	// 	shortDescription: "WP",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "wild pitch"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^SB(2|3|H)$"): PlayMatchConfig{
	// 	// 	shortDescription: "SB%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("stole %s", baseMap[matches[0]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^SB(2|3|H);SB(2|3|H)$"): PlayMatchConfig{
	// 	// 	shortDescription: "DSB %s,%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("double steal %s, %s", baseMap[matches[0]], baseMap[matches[1]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^S(\\d)?$"): PlayMatchConfig{
	// 	// 	shortDescription: "S%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		if matches[0] != "" {
	// 	// 			return fmt.Sprintf("single to %s", FieldingPositions[matches[0]].Code)
	// 	// 		}

	// 	// 		return "single"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^D(\\d)?$"): PlayMatchConfig{
	// 	// 	shortDescription: "D%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		if matches[0] != "" {
	// 	// 			return fmt.Sprintf("double to %s", FieldingPositions[matches[0]].Code)
	// 	// 		}

	// 	// 		return "double"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^T(\\d)?$"): PlayMatchConfig{
	// 	// 	shortDescription: "T%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		if matches[0] != "" {
	// 	// 			return fmt.Sprintf("triple to %s", FieldingPositions[matches[0]].Code)
	// 	// 		}

	// 	// 		return "triple"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^FC(\\d)?$"): PlayMatchConfig{
	// 	// 	shortDescription: "FC%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		if matches[0] != "" {
	// 	// 			return fmt.Sprintf("fielder's choice to %s", FieldingPositions[matches[0]].Code)
	// 	// 		}

	// 	// 		return "fielder's choice"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^E(\\d)?$"): PlayMatchConfig{
	// 	// 	shortDescription: "E%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		if matches[0] != "" {
	// 	// 			return fmt.Sprintf("error by %s", FieldingPositions[matches[0]].Code)
	// 	// 		}

	// 	// 		return "error"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^W$"): PlayMatchConfig{
	// 	// 	shortDescription: "W",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "walk"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^HP$"): PlayMatchConfig{
	// 	// 	shortDescription: "HP",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "hit by pitch"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^I|IW$"): PlayMatchConfig{
	// 	// 	shortDescription: "IW",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "intentional walk"
	// 	// 	},
	// 	// },
	// 	// // The two caught stealing situations below have fielding scenarios after them  to note where the throw was from and to whom
	// 	// // but I'm not interested in logging that
	// 	// regexp.MustCompile("^CS(2|3|H)"): PlayMatchConfig{
	// 	// 	shortDescription: "CS%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("caught stealing %s", baseMap[matches[0]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^POCS(2|3|H)"): PlayMatchConfig{
	// 	// 	shortDescription: "POCS%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("picked off at %s (caught stealing)", baseMap[matches[0]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^PO(1|2|3)\\(\\d\\d\\)"): PlayMatchConfig{
	// 	// 	shortDescription: "PO%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("picked off at %s", baseMap[matches[0]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^PO(1|2|3)\\(E(\\d)\\)"): PlayMatchConfig{
	// 	// 	shortDescription: "PO%s-E%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("pick off attempt at %s, error by %s", baseMap[matches[0]], FieldingPositions[matches[1]].Code)
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^(\\d)(\\d)\\((\\d)\\)$"): PlayMatchConfig{
	// 	// 	shortDescription: "DP%s%s(%s)",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "supposedly some sort of double play"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^(\\d)(\\d)(\\d)$"): PlayMatchConfig{
	// 	// 	shortDescription: "DP-%s%s%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("double play? %s-%s-%s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code, FieldingPositions[matches[2]].Code)
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^(\\d)\\(B\\)(\\d)\\((\\d)\\)$"): PlayMatchConfig{
	// 	// 	shortDescription: "LDP-%s-%s(%s)",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("lined into double play %s to %s at %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code, baseMap[matches[2]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^(\\d)\\(\\d\\)(\\d)$"): PlayMatchConfig{
	// 	// 	shortDescription: "UGDP-%s-%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("ground ball double play, unassisted %s to %s", FieldingPositions[matches[0]].Code, FieldingPositions[matches[1]].Code)
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^(\\d)(\\d)\\((\\d)\\)(\\d)$"): PlayMatchConfig{
	// 	// 	shortDescription: "DP%s%s(%s)%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("%s%s%s DP, runner at %s was the initial out", matches[0], matches[1], matches[3], baseMap[matches[2]])
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^PB$"): PlayMatchConfig{
	// 	// 	shortDescription: "PB",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "passed ball"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^H|HR$"): PlayMatchConfig{
	// 	// 	shortDescription: "HR",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "home run"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^BK$"): PlayMatchConfig{
	// 	// 	shortDescription: "BK",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "balk"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^FLE(\\d)$"): PlayMatchConfig{
	// 	// 	shortDescription: "FLE%s",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return fmt.Sprintf("error by %s on foul fly ball", FieldingPositions[matches[0]].Code)
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^C$"): PlayMatchConfig{
	// 	// 	shortDescription: "C",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "interference"
	// 	// 	},
	// 	// },
	// 	// regexp.MustCompile("^OA$"): PlayMatchConfig{
	// 	// 	shortDescription: "A",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "base runner advance"
	// 	// 	},
	// 	// },
	// 	// // @TODO research all these matches below further
	// 	// regexp.MustCompile("^366\\(1\\)$"): PlayMatchConfig{
	// 	// 	shortDescription: "366(1)",
	// 	// 	longDescription: func(matches []string) string {
	// 	// 		return "UNKNOWN"
	// 	// 	},
	// 	// },
}

func processPlays(plays []string) {
	playRecords := GetRecords(plays)
	fmt.Println("there are", len(playRecords), "plays")

	for _, play := range playRecords {
		config := createPlayConfig(play[6])

		for r, v := range playMatchers {
			playResult := r.FindStringSubmatch(config.basicPlay)

			if len(playResult) > 0 {
				fmt.Println(v(config, playResult[1:]))
				break
			}
		}

		// if hasPlayConfig == false {
		// 	panic(fmt.Sprintf("no play config found for %s", basicPlay))
		// }

		fmt.Println("-----")
	}

	fmt.Println("unique values:", uniqueValues)
}
