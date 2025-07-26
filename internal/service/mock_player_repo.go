package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockPlayerRepository struct {
	sync.RWMutex
	m map[tjudge.PlayerId]tjudge.Player
}

func (r *MockPlayerRepository) Player(ctx context.Context, id tjudge.PlayerId) (tjudge.Player, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Player{}, fmt.Errorf("%w: %s", tjudge.ErrPlayerNotExist, id)
	}
	return contest, nil
}

func (r *MockPlayerRepository) Upsert(ctx context.Context, contest tjudge.Player) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}
