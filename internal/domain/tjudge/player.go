package tjudge

import (
	"context"
	"errors"
	"time"
)

type PlayerId uuid

var ErrPlayerNotExist = errors.New("player doesn't exist")
var ErrInvalidPlayer = errors.New("invalid player passed")

type Player struct {
	Id        PlayerId
	TeamId    TeamId
	Username  string
	CreatedAt time.Time
}

type PlayerRepository interface {
	Player(context.Context, PlayerId) (Player, error)
	Add(context.Context, Player) error
}
