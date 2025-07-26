package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockContestRepository struct {
	sync.RWMutex
	m map[tjudge.ContestId]tjudge.Contest
}

func (r *MockContestRepository) Contest(ctx context.Context, id tjudge.ContestId) (tjudge.Contest, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Contest{}, fmt.Errorf("%w: %s", tjudge.ErrContestNotExist, id)
	}
	return contest, nil
}

func (r *MockContestRepository) Active(ctx context.Context) ([]tjudge.Contest, error) {
	r.Lock()
	defer r.Unlock()
	active := make([]tjudge.Contest, 0)
	for _, v := range r.m {
		if v.Start().Before(time.Now()) && v.End().After(time.Now()) {
			active = append(active, v)
		}
	}
	return active, nil
}

func (r *MockContestRepository) Upsert(ctx context.Context, contest tjudge.Contest) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}
