package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type GetTeam struct {
	Code string
}

type GetTeamHandler struct {
	teamRep tjudge.TeamRepository
}

func NewGetTeamHandler(teamRep tjudge.TeamRepository) GetTeamHandler {
	return GetTeamHandler{teamRep}
}

func (h *GetTeamHandler) Execute(ctx context.Context, cmd GetTeam) (Team, error) {
	team, err := h.teamRep.Team(ctx, cmd.Code)
	if err != nil {
		return Team{}, err
	}
	leader := User{int(team.Leader().Id()), team.Leader().Username(), team.Leader().Fullname()}

	members := make([]*User, len(team.Members()))
	for i, member := range team.Members() {
		members[i] = &User{int(member.Id()), member.Username(), member.Fullname()}
	}

	dteam := Team{
		Code:    string(team.Code()),
		Name:    team.Name(),
		Leader:  &leader,
		Contest: int(team.Contest()),
		MaxSize: team.MaxSize(),
		Members: members,
	}
	return dteam, nil
}
