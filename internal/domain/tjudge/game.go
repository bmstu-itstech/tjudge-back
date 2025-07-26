package tjudge

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type GameId shortUuid

var ErrGameNotExist = errors.New("game doesn't exist")
var ErrInvalidGame = errors.New("invalid game passed")

type Game struct {
	id           GameId
	name         string
	players      uint
	rules_url    string
	allowed_exts []string
}

func (g Game) Id() GameId {
	return g.id
}

func (g Game) Name() string {
	return g.name
}

func (g Game) Players() uint {
	return g.players
}

func (g Game) RulesUrl() string {
	return g.rules_url
}

func (g Game) AllowedExts() []string {
	return g.allowed_exts
}

type GameRepository interface {
	Game(context.Context, GameId) (Game, error)
	Upsert(context.Context, Game) error
}

func ParseGame(id GameId, name string, players uint, rules string, exts []string) (Game, error) {
	if id == "" || name == "" || players == 0 || rules == "" || exts == nil || len(exts) == 0 {
		return Game{}, ErrInvalidGame
	}
	return Game{
		id,
		name,
		players,
		rules,
		exts,
	}, nil
}

func MustParseGame(id GameId, name string, players uint, rules string, exts []string) Game {
	g, err := ParseGame(id, name, players, rules, exts)
	if err != nil {
		panic(err)
	}
	return g
}

func NewGame(name string, players uint, rules string, exts []string) (Game, error) {
	id := uuid.GenerateShort()
	return ParseGame(GameId(id), name, players, rules, exts)
}

func MustNewGame(name string, players uint, rules string, exts []string) Game {
	g, err := NewGame(name, players, rules, exts)
	if err != nil {
		panic(err)
	}
	return g
}
