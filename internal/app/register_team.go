package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type RegisterTeam struct {
	Name string
	Code string
	Contest int
	MaxSize int
}

type RegisterTeamHandler struct {
	teamRep tjudge.TeamRepository
	factory tjudge.TeamFactory
}

func NewRegisterTeamHandler(teamRep tjudge.TeamRepository, factory tjudge.TeamFactory) RegisterTeamHandler {
	return RegisterTeamHandler {teamRep, factory}
}

func (h *RegisterTeamHandler) Execute(ctx context.Context, cmd RegisterTeam) error {
	newteam, err := h.factory.Create(cmd.Name, cmd.Code, tjudge.ContestID(cmd.Contest))
	if err == nil {
		_, err = h.teamRep.Save(ctx, newteam)
	}
	return err;
}
