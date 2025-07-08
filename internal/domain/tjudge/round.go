package tjudge

import (
	"context"
	"errors"

	guuid "github.com/google/uuid"
)

type RoundId uuid

var ErrRoundNotExist = errors.New("round doesn't exist")
var ErrInvalidRound = errors.New("invalid round passed")

type Round struct {
	Id        RoundId
	ResultIds []ResultId
}

func NewRound(results []ResultId) Round {
	id := guuid.New()
	return Round{
		Id:        RoundId(id.String()),
		ResultIds: results,
	}
}

type RoundRepository interface {
	Round(context.Context, RoundId) (Round, error)
	Add(context.Context, Round) error
}
