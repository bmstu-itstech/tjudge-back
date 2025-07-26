package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockResultRepository struct {
	sync.RWMutex
	m map[tjudge.ResultId]tjudge.Result
}

func (r *MockResultRepository) Result(ctx context.Context, id tjudge.ResultId) (tjudge.Result, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Result{}, fmt.Errorf("%w: %s", tjudge.ErrResultNotExist, id)
	}
	return contest, nil
}

func (r *MockResultRepository) Upsert(ctx context.Context, contest tjudge.Result) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}
