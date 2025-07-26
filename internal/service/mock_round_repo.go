package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockRoundRepository struct {
	sync.RWMutex
	m map[tjudge.RoundId]tjudge.Round
}

func (r *MockRoundRepository) Round(ctx context.Context, id tjudge.RoundId) (tjudge.Round, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Round{}, fmt.Errorf("%w: %s", tjudge.ErrRoundNotExist, id)
	}
	return contest, nil
}

func (r *MockRoundRepository) Upsert(ctx context.Context, contest tjudge.Round) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}
