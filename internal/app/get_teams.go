package app

import (
	"context"
	"fmt"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type GetTeamsHandler struct {
	teamRep tjudge.TeamRepository
}

func NewGetTeamsHandler(teamRep tjudge.TeamRepository) *GetTeamsHandler {
	return &GetTeamsHandler{teamRep: teamRep}
}

func (h *GetTeamsHandler) Execute(ctx context.Context) ([]*Team, error) {
	if h == nil {
		return nil, fmt.Errorf("GetTeamsHandler is nil")
	}
	if h.teamRep == nil {
		return nil, fmt.Errorf("team repository is nil")
	}
	
	teams, err := h.teamRep.Teams(ctx)
	if err != nil {
		return nil, err
	}
	var dteams = []*Team{}
	for _, team := range teams {
		// Добавляем проверки на nil для безопасности
		var leader *User
		if team.Leader() != nil {
			leader = &User{int(team.Leader().Id()), team.Leader().Username(), team.Leader().Fullname()}
		}

		members := make([]*User, 0, len(team.Members()))
		for _, member := range team.Members() {
			if member != nil {
				members = append(members, &User{int(member.Id()), member.Username(), member.Fullname()})
			}
		}

		dteam := Team{
			Code:    string(team.Code()),
			Name:    team.Name(),
			Leader:  leader,
			Contest: int(team.Contest()),
			MaxSize: team.MaxSize(),
			Members: members,
		}
		dteams = append(dteams, &dteam)
	}
	return dteams, nil
}