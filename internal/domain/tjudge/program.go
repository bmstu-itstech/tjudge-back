package tjudge

import (
	"context"
	"errors"
	"time"
)

var ErrNoActiveProgram = errors.New("active program doesn't exist")
var ErrEmptyProgramList = errors.New("empty program list")
var ErrInvalidProgram = errors.New("invalid program passed")

type ProgramId int
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
