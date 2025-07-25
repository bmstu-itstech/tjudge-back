package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type ContestId shortUuid

var ErrContestNotExist = errors.New("contest doesn't exist")
var ErrInvalidContest = errors.New("invalid contest passed")

type Contest struct {
	Id        ContestId
	Name      string
	TeamLimit uint
	Start     time.Time
	End       time.Time
	Games     []GameId
}

type ContestRepository interface {
	Contest(context.Context, ContestId) (Contest, error)
	Active(context.Context) ([]Contest, error)
	Upsert(context.Context, Contest) error
}

func ParseContest(id ContestId, name string, lim uint, start time.Time, end time.Time, g []GameId) (Contest, error) {
	if id == "" || name == "" || lim == 0 || start.IsZero() || end.IsZero() || start.After(end) || g == nil {
		return Contest{}, ErrInvalidContest
	}
	return Contest{
		id,
		name,
		lim,
		start,
		end,
		g,
	}, nil
}

func MustParseContest(id ContestId, name string, lim uint, start time.Time, end time.Time, g []GameId) Contest {
	c, err := ParseContest(id, name, lim, start, end, g)
	if err != nil {
		panic(err)
	}
	return c
}

func NewContest(name string, lim uint, start time.Time, end time.Time) (Contest, error) {
	id := uuid.GenerateShort()
	g := make([]GameId, 0)
	return ParseContest(ContestId(id), name, lim, start, end, g)
}

func MustNewContest(name string, lim uint, start time.Time, end time.Time) Contest {
	c, err := NewContest(name, lim, start, end)
	if err != nil {
		panic(err)
	}
	return c
}
