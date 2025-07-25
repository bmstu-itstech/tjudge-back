package tjudge

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type TourId shortUuid

var ErrTourNotExist = errors.New("tour doesn't exist")
var ErrInvalidTour = errors.New("invalid tour passed")

type Tour struct {
	Id       TourId
	GameId   GameId
	RoundIds []RoundId
}

type TourRepository interface {
	Create(context.Context, GameId) (Tour, error) // initiate a tour by a game
	ActiveTour(context.Context, GameId) (Tour, error)
	Tours(context.Context) ([]Tour, error)
	Update(context.Context, Tour) error
}

func ParseTour(id TourId, game GameId, rounds []RoundId) (Tour, error) {
	if id == "" || game == "" || rounds == nil {
		return Tour{}, ErrInvalidTour
	}
	return Tour{id, game, rounds}, nil
}

func MustParseTour(id TourId, game GameId, rounds []RoundId) Tour {
	t, err := ParseTour(id, game, rounds)
	if err != nil {
		panic(err)
	}
	return t
}

func NewTour(game GameId) (Tour, error) {
	// here we assume a new tour doesn't have any rounds because... how?
	id := uuid.GenerateShort()
	rounds := make([]RoundId, 0)
	return ParseTour(TourId(id), game, rounds)
}

func MustNewTour(game GameId) Tour {
	t, err := NewTour(game)
	if err != nil {
		panic(err)
	}
	return t
}
