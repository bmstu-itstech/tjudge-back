package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

var ErrNoActiveProgram = errors.New("active program doesn't exist")
var ErrEmptyProgram = errors.New("empty program code")
var ErrInvalidProgram = errors.New("invalid program passed")

type ProgramId shortUuid
type Path string

type Program struct {
	Id         ProgramId
	TeamId     TeamId
	GameId     GameId
	Path       Path
	UploadedAt time.Time
}

type ProgramRepository interface {
	ActiveProgram(context.Context, Game, Team) (Program, error)
	Program(context.Context, ProgramId) (Program, error)
	Programs(context.Context, Game, Team) ([]Program, error)
	Add(context.Context, Program) error
}

func ParseProgram(id ProgramId, team TeamId, game GameId, path Path, uploaded time.Time) (Program, error) {
	if id == "" || team == "" || game == "" || uploaded.IsZero() {
		return Program{}, ErrInvalidProgram
	}
	if path == "" {
		// TODO: Need a better way to check if program is empty.
		// Are they stored as paths to files on server?
		// Can we just try to read them?
		return Program{}, ErrEmptyProgram
	}

	return Program{
		id,
		team,
		game,
		path,
		uploaded,
	}, nil
}

func MustParseProgram(id ProgramId, team TeamId, game GameId, path Path, uploaded time.Time) Program {
	p, err := ParseProgram(id, team, game, path, uploaded)
	if err != nil {
		panic(err)
	}
	return p
}

func NewProgram(team TeamId, game GameId, path Path) (Program, error) {
	id := uuid.GenerateShort()
	uploaded := time.Now()
	return ParseProgram(ProgramId(id), team, game, path, uploaded)
}

func MustNewProgram(team TeamId, game GameId, path Path) Program {
	p, err := NewProgram(team, game, path)
	if err != nil {
		panic(err)
	}
	return p
}
