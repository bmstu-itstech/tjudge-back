// get_team.go
package app

import (
	"context"
	"fmt"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type GetTeam struct {
	Code string
}

type GetTeamHandler struct {
	teamRep tjudge.TeamRepository
}

func NewGetTeamHandler(teamRep tjudge.TeamRepository) *GetTeamHandler {
	return &GetTeamHandler{teamRep: teamRep}
}

func (h *GetTeamHandler) Execute(ctx context.Context, request GetTeam) (*Team, error) {
	if h == nil {
		return nil, fmt.Errorf("GetTeamHandler is nil")
	}
	if h.teamRep == nil {
		return nil, fmt.Errorf("team repository is nil")
	}
	
	team, err := h.teamRep.Team(ctx, request.Code)
	if err != nil {
		return nil, err
	}

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

	dteam := &Team{
		Code:    string(team.Code()),
		Name:    team.Name(),
		Leader:  leader,
		Contest: int(team.Contest()),
		MaxSize: team.MaxSize(),
		Members: members,
	}
	return dteam, nil
}