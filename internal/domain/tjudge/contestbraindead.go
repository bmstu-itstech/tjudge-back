package tjudge

import "time"

type ContestAlt struct {
	Id        ContestId
	Name      string
	TeamLimit uint
	Starts    time.Time
	Ends      time.Time
	Games     []ContestAlt
}

type GameAlt struct {
	Id          GameId
	Name        string
	Players     uint
	RulesUrl    string
	AllowedExts []string

	ActivePrograms map[TeamId]ProgramId
	Rounds         []Round2ElectricBoogaloo
}

type Round2ElectricBoogaloo struct {
	State    string // scheduled/finished/aborted
	Programs map[TeamId]ProgramId
	Results  map[TeamId]int
	Time     time.Time
}
