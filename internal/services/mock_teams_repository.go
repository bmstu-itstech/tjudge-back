package services

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockTeamRepository struct {
	m     map[tjudge.TeamCode]*tjudge.Team
	index int
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		m:     make(map[tjudge.TeamCode]*tjudge.Team),
		index: 1,
	}
}

func (r *MockTeamRepository) Save(ctx context.Context, proto *tjudge.Team) (*tjudge.Team, error) {
	for _, team := range r.m {
		if team.Code() == proto.Code() {
			return nil, tjudge.ErrUserAlreadyExists
		}
	}
	r.m[proto.Code()] = proto
	return proto, nil
}

func (r *MockTeamRepository) Team(ctx context.Context, code tjudge.TeamCode) (*tjudge.Team, error) {
	for _, team := range r.m {
		if team.Code() == code {
			return team, nil
		}
	}
	return nil, tjudge.ErrTeamNotFound
}

func (r *MockTeamRepository) Teams(ctx context.Context) ([]*tjudge.Team, error) {
	if len(r.m) == 0 {
		return nil, tjudge.ErrTeamNotFound
	}
	var teams []*tjudge.Team
	for _, team := range r.m {
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *MockTeamRepository) TeamsByContest(ctx context.Context, contest tjudge.ContestID) ([]*tjudge.Team, error) {
	var teams []*tjudge.Team
	for _, team := range r.m {
		if team.Contest() == contest {
			teams = append(teams, team)
		}
	}
	if len(teams) == 0 {
		return nil, tjudge.ErrTeamNotFound
	}
	return teams, nil
}

func (r *MockTeamRepository) Update(ctx context.Context, team *tjudge.Team) error {
	for i, tteam := range r.m {
		if tteam.Code() == team.Code() {
			r.m[i] = team
			return nil
		}
	}
	return tjudge.ErrTeamNotFound
}

func (r *MockTeamRepository) Delete(ctx context.Context, code tjudge.TeamCode) error {
	for i, tteam := range r.m {
		if tteam.Code() == code {
			delete(r.m, i)
			return nil
		}
	}
	return tjudge.ErrTeamNotFound
}
