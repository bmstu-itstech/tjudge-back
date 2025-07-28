package app

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type Contest struct {
	Id        string
	Name      string
	TeamLimit uint
	Starts    time.Time
	Ends      time.Time
	GameIds   []string
}

type Game struct {
	Id          string
	Name        string
	Players     uint
	RulesUrl    string
	AllowedExts []string
}

type Player struct {
	Id        string
	TeamId    string
	Username  string
	CreatedAt time.Time
}

type Program struct {
	Id         string
	TeamId     string
	GameId     string
	Path       string
	UploadedAt time.Time
}

type Result struct {
	Id        string
	ProgramId string
	Score     int
}

type Round struct {
	Id        string
	ResultIds []string
}

type Team struct {
	Id        string
	Name      string
	ContestId string
	CreatedAt time.Time
	JoinCode  string
}

type Tour struct {
	Id        string
	GameId    string
	RoundIds  []string
	CreatedAt time.Time
}

func gameIdsToDto(g []tjudge.GameId) []string {
	ids := make([]string, 0, len(g))
	for _, v := range g {
		ids = append(ids, string(v))
	}
	return ids
}

func gameIdsFromDto(g []string) []tjudge.GameId {
	ids := make([]tjudge.GameId, 0, len(g))
	for _, v := range g {
		ids = append(ids, tjudge.GameId(v))
	}
	return ids
}

func contestToDto(c tjudge.Contest) Contest {
	return Contest{
		string(c.Id()),
		c.Name(),
		c.TeamLimit(),
		c.Starts(),
		c.Ends(),
		gameIdsToDto(c.GameIds()),
	}
}

func batchContestsToDto(cs []tjudge.Contest) []Contest {
	out := make([]Contest, 0, len(cs))
	for _, c := range cs {
		out = append(out, contestToDto(c))
	}
	return out
}

func gameToDto(g tjudge.Game) Game {
	return Game{
		string(g.Id()),
		g.Name(),
		g.Players(),
		g.RulesUrl(),
		g.AllowedExts(),
	}
}

func roundIdsToDto(r []tjudge.RoundId) []string {
	ids := make([]string, 0, len(r))
	for _, v := range r {
		ids = append(ids, string(v))
	}
	return ids
}

func tourToDto(t tjudge.Tour) Tour {
	return Tour{
		string(t.Id()),
		string(t.GameId()),
		roundIdsToDto(t.RoundIds()),
		t.CreatedAt(),
	}
}

func batchToursToDto(ts []tjudge.Tour) []Tour {
	out := make([]Tour, 0, len(ts))
	for _, t := range ts {
		out = append(out, tourToDto(t))
	}
	return out
}

func resultToDto(r tjudge.Result) Result {
	return Result{
		string(r.Id()),
		string(r.ProgramId()),
		r.Score(),
	}
}

func resultIdsToDto(r []tjudge.ResultId) []string {
	ids := make([]string, 0, len(r))
	for _, v := range r {
		ids = append(ids, string(v))
	}
	return ids
}

func roundToDto(r tjudge.Round) Round {
	return Round{
		string(r.Id()),
		resultIdsToDto(r.ResultIds()),
	}
}
