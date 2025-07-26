package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockTeamRepository struct {
	sync.RWMutex
	m map[tjudge.TeamId]tjudge.Team
}

func (r *MockTeamRepository) Team(ctx context.Context, id tjudge.TeamId) (tjudge.Team, error) {
	r.RLock()
	defer r.RUnlock()
	contest, ok := r.m[id]
	if !ok {
		return tjudge.Team{}, fmt.Errorf("%w: %s", tjudge.ErrTeamNotExist, id)
	}
	return contest, nil
}

func (r *MockTeamRepository) Upsert(ctx context.Context, contest tjudge.Team) error {
	r.Lock()
	defer r.Unlock()
	r.m[contest.Id()] = contest
	return nil
}

func (r *MockTeamRepository) ByContest(ctx context.Context, id tjudge.ContestId) ([]tjudge.Team, error) {
	r.Lock()
	defer r.Unlock()
	teams := make([]tjudge.Team, 0)
	for _, v := range r.m {
		if v.ContestId() == id {
			teams = append(teams, v)
		}
	}
	return teams, nil
}

func (r *MockTeamRepository) ByJoinCode(ctx context.Context, code string) (tjudge.Team, error) {
	r.Lock()
	defer r.Unlock()
	for _, v := range r.m {
		if v.JoinCode() == code {
			return v, nil
		}
	}
	return tjudge.Team{}, fmt.Errorf("%w (code:%s", tjudge.ErrTeamNotExist, code)
}
