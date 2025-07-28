package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type TourId shortUuid

var ErrNoActiveTour = errors.New("there is no active tour")
var ErrTourNotExist = errors.New("tour doesn't exist")
var ErrInvalidTour = errors.New("invalid tour passed")

type Tour struct {
	id         TourId
	game_id    GameId
	round_ids  []RoundId
	created_at time.Time
}

func (t Tour) Id() TourId {
	return t.id
}

func (t Tour) GameId() GameId {
	return t.game_id
}

func (t Tour) RoundIds() []RoundId {
	return t.round_ids
}

func (t Tour) CreatedAt() time.Time {
	return t.created_at
}

type TourRepository interface {
	Tour(context.Context, TourId) (Tour, error)
	Active(context.Context, GameId) (Tour, error)
	Tours(context.Context, GameId) ([]Tour, error)
	Upsert(context.Context, Tour) error
}

func ParseTour(id TourId, game GameId, rounds []RoundId, created_at time.Time) (Tour, error) {
	if id == "" || game == "" || rounds == nil || created_at.IsZero() {
		return Tour{}, ErrInvalidTour
	}
	return Tour{id, game, rounds, created_at}, nil
}

func MustParseTour(id TourId, game GameId, rounds []RoundId, created_at time.Time) Tour {
	t, err := ParseTour(id, game, rounds, created_at)
	if err != nil {
		panic(err)
	}
	return t
}

func NewTour(game GameId, rounds []RoundId) (Tour, error) {
	id := uuid.GenerateShort()
	created_at := time.Now()
	return ParseTour(TourId(id), game, rounds, created_at)
}

func MustNewTour(game GameId, rounds []RoundId) Tour {
	t, err := NewTour(game, rounds)
	if err != nil {
		panic(err)
	}
	return t
}
