package tjudge

import "errors"

var ErrDifferentProgramType = errors.New("invalid program game type")
var ErrJudgeMissing = errors.New("testing program not found") // time to panic?
var ErrInvalidPlayerCount = errors.New("invalid amount of game players")

type ErrProgramException struct {
	msg string
}

func NewErrProgramException(msg string) ErrProgramException {
	return ErrProgramException{msg: msg}
}

func (e ErrProgramException) Error() string {
	return e.msg
}

// Result sans the id
type RunResult struct {
	ProgramId ProgramId
	Score     int
}

type Launcher interface {
	Run(Game, []Program) ([]RunResult, error)
}
