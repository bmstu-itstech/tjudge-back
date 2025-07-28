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
var ErrNoLateJoin = errors.New("cannot join contest after it started")

type Contest struct {
	id         ContestId
	name       string
	team_limit uint
	starts     time.Time
	ends       time.Time
	game_ids   []GameId
}

func (c Contest) Id() ContestId {
	return c.id
}

func (c Contest) Name() string {
	return c.name
}

func (c Contest) TeamLimit() uint {
	return c.team_limit
}

func (c Contest) Starts() time.Time {
	return c.starts
}

func (c Contest) Ends() time.Time {
	return c.ends
}

func (c Contest) GameIds() []GameId {
	return c.game_ids
}

type ContestRepository interface {
	Contest(context.Context, ContestId) (Contest, error)
	Active(context.Context) ([]Contest, error)
	Upsert(context.Context, Contest) error
}

func ParseContest(id ContestId, name string, lim uint, starts time.Time, ends time.Time, g []GameId) (Contest, error) {
	if id == "" || name == "" || lim == 0 || starts.IsZero() || ends.IsZero() || ends.Before(starts) || g == nil {
		return Contest{}, ErrInvalidContest
	}
	return Contest{
		id,
		name,
		lim,
		starts,
		ends,
		g,
	}, nil
}

func MustParseContest(id ContestId, name string, lim uint, starts time.Time, ends time.Time, g []GameId) Contest {
	c, err := ParseContest(id, name, lim, starts, ends, g)
	if err != nil {
		panic(err)
	}
	return c
}

func NewContest(name string, lim uint, starts time.Time, ends time.Time, g []GameId) (Contest, error) {
	id := uuid.GenerateShort()
	return ParseContest(ContestId(id), name, lim, starts, ends, g)
}

func MustNewContest(name string, lim uint, starts time.Time, ends time.Time, g []GameId) Contest {
	c, err := NewContest(name, lim, starts, ends, g)
	if err != nil {
		panic(err)
	}
	return c
}
