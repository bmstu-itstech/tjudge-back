package contest

import (
	"errors"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

var ErrInvalidGame = errors.New("invalid game passed")
var ErrGameNotExist = errors.New("game doesn't exist")

type Game struct {
	Id       shared.ID
	Name     string
	RulesUrl string

	Matches map[shared.ID]*Match
}

func ParseGame(id shared.ID, name string, rules_url string, matches map[shared.ID]*Match) (Game, error) {
	// Badly formatted urls? Not MY problem :D
	if len(id) == 0 || name == "" || rules_url == "" || matches == nil {
		return Game{}, ErrInvalidGame
	}
	return Game{id, name, rules_url, matches}, nil
}

func MustParseGame(id shared.ID, name string, rules_url string, matches map[shared.ID]*Match) Game {
	g, err := ParseGame(id, name, rules_url, matches)
	if err != nil {
		panic(err)
	}
	return g
}

func NewGame(name string, rules_url string) (Game, error) {
	return ParseGame(shared.NewID(), name, rules_url, make(map[shared.ID]*Match))
}

func MustNewGame(name string, rules_url string) Game {
	g, err := NewGame(name, rules_url)
	if err != nil {
		panic(err)
	}
	return g
}
