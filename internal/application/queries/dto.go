package queries

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/google/uuid"
)

type Contest struct {
	Id      string
	Name    string
	Starts  time.Time
	Ends    time.Time
	GameIds []string
}

type Game struct {
	Id       string
	Name     string
	RulesUrl string
}

type Program struct {
	Id        string
	ContestId string
	GameId    string
	TeamId    string
	Contents  []byte
}

func mapToIds[E any](m map[shared.ID]*E) []string {
	out := make([]string, 0)
	for i := range m {
		out = append(out, uuid.UUID(i).String())
	}
	return out
}

func contestToDto(c *contest.Contest) Contest {
	return Contest{
		uuid.UUID(c.Id).String(),
		c.Name,
		c.Starts,
		c.Ends,
		mapToIds(c.Games),
	}
}

func batchContestsToDto(cs []*contest.Contest) []Contest {
	out := make([]Contest, 0, len(cs))
	for _, c := range cs {
		out = append(out, contestToDto(c))
	}
	return out
}

func gameToDto(g *contest.Game) Game {
	return Game{uuid.UUID(g.Id).String(), g.Name, g.RulesUrl}
}

func batchGamesToDto(gs []*contest.Game) []Game {
	out := make([]Game, 0, len(gs))
	for _, g := range gs {
		out = append(out, gameToDto(g))
	}
	return out
}

func standingsToDto(s map[shared.ID]contest.Score) map[string]int {
	out := make(map[string]int)
	for i, v := range s {
		out[uuid.UUID(i).String()] = int(v)
	}
	return out
}

func programToDto(p program.Program, file []byte) Program {
	return Program{
		uuid.UUID(p.Id).String(),
		uuid.UUID(p.ContestId).String(),
		uuid.UUID(p.GameId).String(),
		uuid.UUID(p.TeamId).String(),
		file,
	}
}
