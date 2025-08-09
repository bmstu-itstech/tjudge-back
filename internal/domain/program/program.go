package program

import (
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

var ErrNoActiveProgram = errors.New("there is no active program")
var ErrInvalidProgram = errors.New("invalid program passed")
var ErrProgramNotExist = errors.New("program doesn't exist")

type Path string

type Program struct {
	Id        shared.ID
	ContestId shared.ID
	GameId    shared.ID
	TeamId    shared.ID
	Path      Path
}

func ParseProgram(id shared.ID, teamId shared.ID, contestId shared.ID, gameId shared.ID, path Path) (Program, error) {
	if len(id) == 0 || len(teamId) == 0 || len(contestId) == 0 || len(gameId) == 0 || path == "" {
		return Program{}, ErrInvalidProgram
	}
	return Program{id, teamId, contestId, gameId, path}, nil
}

func MustParseProgram(id shared.ID, teamId shared.ID, contestId shared.ID, gameId shared.ID, path Path) Program {
	p, err := ParseProgram(id, teamId, contestId, gameId, path)
	if err != nil {
		panic(err)
	}
	return p
}

func NewProgram(teamId shared.ID, contestId shared.ID, gameId shared.ID, path Path) (Program, error) {
	return ParseProgram(shared.NewID(), teamId, contestId, gameId, path)
}

func MustNewProgram(teamId shared.ID, contestId shared.ID, gameId shared.ID, path Path) Program {
	p, err := NewProgram(teamId, contestId, gameId, path)
	if err != nil {
		panic(err)
	}
	return p
}

func New(teamId shared.ID, contestId shared.ID, gameId shared.ID, path Path) (*Program, UploadedEvent, error) {
	p, err := NewProgram(teamId, contestId, gameId, path)
	if err != nil {
		return nil, UploadedEvent{}, err
	}

	e := UploadedEvent{
		ContestId: contestId,
		GameId:    gameId,
		TeamId:    teamId,
		ProgramId: p.Id,
		Timestamp: time.Now(),
	}

	return &p, e, nil
}
