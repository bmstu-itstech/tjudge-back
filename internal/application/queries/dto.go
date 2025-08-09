package queries

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
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
