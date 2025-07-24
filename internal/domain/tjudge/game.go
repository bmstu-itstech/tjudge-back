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
	AllowedExts  []string
}

type GameRepository interface {
	Game(context.Context, GameId) (Game, error)
	Add(context.Context, Game) error
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
