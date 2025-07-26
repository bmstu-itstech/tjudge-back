package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockTourRepository struct {
	sync.RWMutex
	m map[tjudge.TourId]tjudge.Tour
}

func (r *MockTourRepository) Tour(ctx context.Context, id tjudge.TourId) (tjudge.Tour, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Tour{}, fmt.Errorf("%w: %s", tjudge.ErrTourNotExist, id)
	}
	return contest, nil
}

func (r *MockTourRepository) Upsert(ctx context.Context, contest tjudge.Tour) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}

func (r *MockTourRepository) Active(ctx context.Context, id tjudge.GameId) (tjudge.Tour, error) {
	r.Lock()
	defer r.Unlock()
	// the crutch returns
	tour := tjudge.MustParseTour("id", id, make([]tjudge.RoundId, 0), time.Unix(0, 1))
	for _, v := range r.m {
		if v.GameId() == id && v.CreatedAt().After(tour.CreatedAt()) {
			tour = v
		}
	}
	if tour.CreatedAt().Equal(time.Unix(0, 1)) {
		return tjudge.Tour{}, fmt.Errorf("%w: game:%s", tjudge.ErrNoActiveTour, id)
	}
	return tour, nil
}

func (r *MockTourRepository) Tours(ctx context.Context, id tjudge.GameId) ([]tjudge.Tour, error) {
	r.Lock()
	defer r.Unlock()
	tours := make([]tjudge.Tour, 0)
	for _, v := range r.m {
		if v.GameId() == id {
			tours = append(tours, v)
		}
	}
	return tours, nil
}
