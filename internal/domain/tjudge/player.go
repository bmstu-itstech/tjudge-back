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
	id         PlayerId
	team_id    TeamId
	username   string
	created_at time.Time
}

func (p Player) Id() PlayerId {
	return p.id
}

func (p Player) TeamId() TeamId {
	return p.team_id
}

func (p Player) Username() string {
	return p.username
}

func (p Player) CreatedAt() time.Time {
	return p.created_at
}

type PlayerRepository interface {
	Player(context.Context, PlayerId) (Player, error)
	Upsert(context.Context, Player) error
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
