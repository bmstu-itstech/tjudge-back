package tjudge

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type ResultId shortUuid

var ErrResultNotExist = errors.New("result doesn't exist")
var ErrInvalidResult = errors.New("invalid result passed")

type Result struct {
	id         ResultId
	program_id ProgramId
	score      int
}

func (r Result) Id() ResultId {
	return r.id
}

func (r Result) ProgramId() ProgramId {
	return r.program_id
}

func (r Result) Score() int {
	return r.score
}

type ResultRepository interface {
	Result(context.Context, ResultId) (Result, error)
	Upsert(context.Context, Result) error
}

func ParseResult(id ResultId, program ProgramId, score int) (Result, error) {
	// TODO: can a result be negative? can we *lose* real hard?
	if id == "" || program == "" || score < 0 {
		return Result{}, ErrInvalidResult
	}
	return Result{id, program, score}, nil
}

func MustParseResult(id ResultId, program ProgramId, score int) Result {
	r, err := ParseResult(id, program, score)
	if err != nil {
		panic(err)
	}
	return r
}

func NewResult(program ProgramId, score int) (Result, error) {
	id := uuid.GenerateShort()
	return ParseResult(ResultId(id), program, score)
}

func MustNewResult(program ProgramId, score int) Result {
	r, err := NewResult(program, score)
	if err != nil {
		panic(err)
	}
	return r
}
