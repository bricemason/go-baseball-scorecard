package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Info Represents all the game and administrative information tied to a game
type Info struct {
	Visteam    Team   `json:"visitingTeam"`
	Hometeam   Team   `json:"homeTeam"`
	Site       Park   `json:"site"`
	Date       string `json:"date"`
	GameType   string `json:"gameType"`
	Starttime  string `json:"startTime"`
	Daynight   string `json:"dayNight"`
	Usedh      bool   `json:"useDH"`
	Umphome    Person `json:"homeUmpire"`
	Ump1b      Person `json:"1bUmpire"`
	Ump2b      Person `json:"2bUmpire"`
	Ump3b      Person `json:"3bUmpire"`
	Scorer     string `json:"scorer"`
	Inputter   string `json:"inputter"`
	Howscored  string `json:"howScored"`
	Pitches    string `json:"pitches"`
	Temp       int    `json:"temp"`
	Winddir    string `json:"windDirection"`
	Windspeed  int    `json:"windSpeed"`
	Fieldcond  string `json:"fieldConditions"`
	Precip     string `json:"precipitation"`
	Sky        string `json:"sky"`
	Timeofgame int    `json:"timeOfGame"`
	Attendance int    `json:"attendance"`
	Wp         Person `json:"winningPitcher"`
	Lp         Person `json:"losingPitcher"`
	Save       Person `json:"save"`
	Gwrbi      Person `json:"gameWinningRBI"`
}

// CreateInfo Given a slice of raw comma-delimited strings, return an Info struct
func CreateInfo(source []string) Info {
	info := Info{}
	infoRecords := GetRecords(source)

	for _, r := range infoRecords {
		switch r[1] {
		case "visteam":
			info.Visteam = Teams[r[2]]
		case "hometeam":
			info.Hometeam = Teams[r[2]]
		case "site":
			info.Site = Parks[r[2]]
		case "date":
			info.Date = r[2]
		case "number":
			{
				gameType, _ := strconv.Atoi(r[2])
				info.GameType = GameTypes[gameType]
			}
		case "starttime":
			info.Starttime = r[2]
		case "daynight":
			info.Daynight = strings.Title(r[2])
		case "usedh":
			useDh, _ := strconv.ParseBool(r[2])
			info.Usedh = useDh
		case "umphome":
			info.Umphome = People[r[2]]
		case "ump1b":
			info.Ump1b = People[r[2]]
		case "ump2b":
			info.Ump2b = People[r[2]]
		case "ump3b":
			info.Ump3b = People[r[2]]
		case "scorer":
			info.Scorer = r[2]
		case "inputter":
			info.Inputter = r[2]
		case "howscored":
			info.Howscored = r[2]
		case "pitches":
			info.Pitches = r[2]
		case "temp":
			{
				temp, _ := strconv.Atoi(r[2])
				info.Temp = temp
			}
		case "winddir":
			info.Winddir = WindDirections[r[2]]
		case "windspeed":
			{
				speed, _ := strconv.Atoi(r[2])
				info.Windspeed = speed
			}
		case "fieldcond": // dry, soaked, wet, unknown
			info.Fieldcond = strings.Title(r[2])
		case "precip": // drizzle, none, rain, showers, snow, unknown
			info.Precip = strings.Title(r[2])
		case "sky": // cloudy, dome, night, overcast, sunny, unknown
			info.Sky = strings.Title(r[2])
		case "timeofgame":
			{
				time, _ := strconv.Atoi(r[2])
				info.Timeofgame = time
			}
		case "attendance":
			{
				attendance, _ := strconv.Atoi(r[2])
				info.Attendance = attendance
			}
		case "wp":
			info.Wp = People[r[2]]
		case "lp":
			info.Lp = People[r[2]]
		case "save":
			info.Save = People[r[2]]
		case "gwrbi":
			info.Gwrbi = People[r[2]]
		default:
			panic(fmt.Sprintf("The field '%s' does not have a handler", r[1]))
		}
	}

	return info
}
