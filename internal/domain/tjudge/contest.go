package tjudge

import (
	"context"
	"errors"
)

type ContestId uuid

var ErrContestNotExist = errors.New("contest doesn't exist")
var ErrInvalidContest = errors.New("invalid contest passed")

type Contest struct {
	Id       ContestId
	GameId   GameId
	RoundIds []RoundId
}

type ContestRepository interface {
	Create(context.Context, GameId) (Contest, error)
	ActiveContest(context.Context, GameId) (Contest, error)
	Contests(context.Context) ([]Contest, error)
	Update(context.Context, Contest) error
}
