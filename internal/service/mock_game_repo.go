package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockGameRepository struct {
	sync.RWMutex
	m map[tjudge.GameId]tjudge.Game
}

func (r *MockGameRepository) Game(ctx context.Context, id tjudge.GameId) (tjudge.Game, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Game{}, fmt.Errorf("%w: %s", tjudge.ErrGameNotExist, id)
	}
	return contest, nil
}

func (r *MockGameRepository) Upsert(ctx context.Context, contest tjudge.Game) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}
