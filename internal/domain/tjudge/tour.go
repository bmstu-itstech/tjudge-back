package tjudge

import (
	"context"
	"errors"
)

type TourId uuid

var ErrTourNotExist = errors.New("tour doesn't exist")
var ErrInvalidTour = errors.New("invalid tour passed")

type Tour struct {
	Id       TourId
	GameId   GameId
	RoundIds []RoundId
}

type TourRepository interface {
	Create(context.Context, GameId) (Tour, error)
	ActiveTour(context.Context, GameId) (Tour, error)
	Tours(context.Context) ([]Tour, error)
	Update(context.Context, Tour) error
}
