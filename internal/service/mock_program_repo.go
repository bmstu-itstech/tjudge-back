package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockProgramRepository struct {
	sync.RWMutex
	m map[tjudge.ProgramId]tjudge.Program
}

func (r *MockProgramRepository) Program(ctx context.Context, id tjudge.ProgramId) (tjudge.Program, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Program{}, fmt.Errorf("%w: %s", tjudge.ErrProgramNotExist, id)
	}
	return contest, nil
}

func (r *MockProgramRepository) Upsert(ctx context.Context, contest tjudge.Program) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}

func (r *MockProgramRepository) Programs(ctx context.Context, game tjudge.GameId, team tjudge.TeamId) ([]tjudge.Program, error) {
	r.Lock()
	defer r.Unlock()
	programs := make([]tjudge.Program, 0)
	for _, v := range r.m {
		if v.GameId() == game && v.TeamId() == team {
			programs = append(programs, v)
		}
	}
	return programs, nil
}

func (r *MockProgramRepository) ActiveProgram(ctx context.Context, game tjudge.GameId, team tjudge.TeamId) (tjudge.Program, error) {
	r.Lock()
	defer r.Unlock()
	// can you smell it? the crutch?
	program := tjudge.MustParseProgram("id", team, game, "path", time.Unix(0, 1))
	for _, v := range r.m {
		if v.GameId() == game && v.TeamId() == team && v.UploadedAt().After(program.UploadedAt()) {
			program = v
		}
	}
	if program.UploadedAt().Equal(time.Unix(0, 1)) {
		return tjudge.Program{},  fmt.Errorf("%w: game:%s/team:%s", tjudge.ErrNoActiveProgram, game, team)
	}
	return program, nil
}
