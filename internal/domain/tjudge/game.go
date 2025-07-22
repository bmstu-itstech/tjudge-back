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
	Id       GameId
	Name     string
	Players  uint
	RulesUrl string
}

type GameRepository interface {
	Game(context.Context, GameId) (Game, error)
	Add(context.Context, Game) error
}

func ParseGame(id GameId, name string, players uint, rules string) (Game, error) {
	if id == "" || name == "" || players == 0 || rules == "" {
		return Game{}, ErrInvalidGame
	}
	return Game{
		id,
		name,
		players,
		rules,
	}, nil
}

func MustParseGame(id GameId, name string, players uint, rules string) Game {
	g, err := ParseGame(id, name, players, rules)
	if err != nil {
		panic(err)
	}
	return g
}

func NewGame(name string, players uint, rules string) (Game, error) {
	id := uuid.GenerateShort()
	return ParseGame(GameId(id), name, players, rules)
}

func MustNewGame(name string, players uint, rules string) Game {
	g, err := NewGame(name, players, rules)
	if err != nil {
		panic(err)
	}
	return g
}
