package tjudge

import (
	"context"
	"errors"
)

var ErrProgramTooLarge = errors.New("program is too large")
var ErrWrongProgramExt = errors.New("invalid program extension")
var ErrEmptyProgram = errors.New("empty program code")
var ErrInvalidProgramSrc = errors.New("invalid program source passed")

type ProgramSource struct {
	id   ProgramId
	code []byte
	ext  string
}

func (p ProgramSource) Id() ProgramId {
	return p.id
}

func (p ProgramSource) Code() []byte {
	return p.code
}

func (p ProgramSource) Ext() string {
	return p.ext
}

type ProgramSourceRepository interface {
	Upsert(context.Context, ProgramSource) error
	ProgramSource(context.Context, ProgramId) (ProgramSource, error)
	Delete(context.Context, ProgramId) error
}

func NewProgramSource(id ProgramId, code []byte, ext string) (ProgramSource, error) {
	if id == "" || code == nil || ext == "" {
		return ProgramSource{}, ErrInvalidProgramSrc
	}
	if len(code) == 0 {
		return ProgramSource{}, ErrEmptyProgram
	}
	return ProgramSource{id, code, ext}, nil
}

func MustNewProgramSource(id ProgramId, code []byte, ext string) ProgramSource {
	p, err := NewProgramSource(id, code, ext)
	if err != nil {
		panic(err)
	}
	return p
}

// reminder: you should get the ID from NewProgram
