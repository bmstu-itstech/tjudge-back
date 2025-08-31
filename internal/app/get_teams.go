package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type GetTeamsHandler struct {
	teamRep tjudge.TeamRepository
}

func NewGetTeamsHandler(teamRep tjudge.TeamRepository) GetTeamsHandler {
	return GetTeamsHandler{teamRep}
}

func (h *GetTeamsHandler) Execute(ctx context.Context) ([]*Team, error) {
	teams, err := h.teamRep.Teams(ctx)
	if err != nil {
		return nil, err
	}
	var dteams = []*Team{}
	for _, team := range teams {
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
		dteams = append(dteams, &dteam)
	}
	return dteams, nil
}
