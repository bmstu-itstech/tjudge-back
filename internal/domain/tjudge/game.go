package tjudge

import (
	"context"
	"errors"
)

type GameId string

var ErrGameNotExist = errors.New("game doesn't exist")
var ErrInvalidGame = errors.New("invalid game passed")

type Game struct {
	Id      GameId
	Players uint
	Name    string
	Rules   string
}

type GameRepository interface {
	Game(context.Context, GameId) (Game, error)
	Add(context.Context, Game) error
}
