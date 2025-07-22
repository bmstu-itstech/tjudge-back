package tjudge

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type RoundId shortUuid

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

func ParseRound(id RoundId, results []ResultId) (Round, error) {
	if id == "" || results == nil {
		return Round{}, ErrInvalidRound
	}
	return Round{id, results}, nil
}

func MustParseRound(id RoundId, results []ResultId) Round {
	r, err := ParseRound(id, results)
	if err != nil {
		panic(err)
	}
	return r
}

func NewRound(results []ResultId) (Round, error) {
	// would we need a function which also initialises results?
	id := uuid.GenerateShort()
	return ParseRound(RoundId(id), results)
}

func MustNewRound(results []ResultId) Round {
	r, err := NewRound(results)
	if err != nil {
		panic(err)
	}
	return r
}
