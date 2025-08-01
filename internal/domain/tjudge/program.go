package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

var ErrNoActiveProgram = errors.New("there is no active program")
var ErrInvalidProgram = errors.New("invalid program passed")
var ErrProgramNotExist = errors.New("program doesn't passed")

type ProgramId shortUuid
type Path string

type Program struct {
	id          ProgramId
	team_id     TeamId
	game_id     GameId
	uploaded_at time.Time
}

func (p Program) Id() ProgramId {
	return p.id
}

func (p Program) TeamId() TeamId {
	return p.team_id
}

func (p Program) GameId() GameId {
	return p.game_id
}

func (p Program) UploadedAt() time.Time {
	return p.uploaded_at
}

type ProgramRepository interface {
	Active(context.Context, GameId, TeamId) (Program, error)
	Actives(context.Context, GameId) (map[TeamId]ProgramId, error)
	Program(context.Context, ProgramId) (Program, error)
	Programs(context.Context, GameId, TeamId) ([]Program, error)
	Upsert(context.Context, Program) error
}

func ParseProgram(id ProgramId, team TeamId, game GameId, uploaded time.Time) (Program, error) {
	if id == "" || team == "" || game == "" || uploaded.IsZero() {
		return Program{}, ErrInvalidProgram
	}

	return Program{
		id,
		team,
		game,
		uploaded,
	}, nil
}

func MustParseProgram(id ProgramId, team TeamId, game GameId, uploaded time.Time) Program {
	p, err := ParseProgram(id, team, game, uploaded)
	if err != nil {
		panic(err)
	}
	return p
}

func NewProgram(team TeamId, game GameId) (Program, error) {
	id := uuid.GenerateShort()
	uploaded := time.Now()
	return ParseProgram(ProgramId(id), team, game, uploaded)
}

func MustNewProgram(team TeamId, game GameId, path Path) Program {
	p, err := NewProgram(team, game)
	if err != nil {
		panic(err)
	}
	return p
}
