package tjudge

import (
	"context"
	"errors"
)

type ResultId uuid

var ErrResultNotExist = errors.New("result doesn't exist")
var ErrInvalidResult = errors.New("invalid result passed")

type Result struct {
	Id        ResultId
	ProgramId ProgramId
	Score     int
}

type ResultRepository interface {
	Result(context.Context, ResultId) (Result, error)
	Add(context.Context, Result) error
}
