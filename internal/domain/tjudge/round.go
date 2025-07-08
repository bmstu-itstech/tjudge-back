package tjudge

import (
	"context"
	"errors"
)

type RoundId uuid

var ErrRoundNotExist = errors.New("round doesn't exist")
var ErrInvalidRound = errors.New("invalid round passed")

type Round struct {
	Id        RoundId
	ResultIds []ResultId
}

type RoundRepository interface {
	Round(context.Context, RoundId) (Round, error)
	Add(context.Context, Round) error
}
