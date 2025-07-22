package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type PlayerId shortUuid

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

func ParsePlayer(id PlayerId, team TeamId, name string, created time.Time) (Player, error) {
	if id == "" || team == "" || name == "" || created.IsZero() {
		return Player{}, ErrInvalidPlayer
	}
	return Player{
		id,
		team,
		name,
		created,
	}, nil
}

func MustParsePlayer(id PlayerId, team TeamId, name string, created time.Time) Player {
	g, err := ParsePlayer(id, team, name, created)
	if err != nil {
		panic(err)
	}
	return g
}

func NewPlayer(team TeamId, name string) (Player, error) {
	id := uuid.GenerateShort()
	created := time.Now()
	return ParsePlayer(PlayerId(id), team, name, created)
}

func MustNewPlayer(team TeamId, name string) Player {
	p, err := NewPlayer(team, name)
	if err != nil {
		panic(err)
	}
	return p
}
