package app

import (
	"context"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type DeleteTeam struct {
	AuthorID int
	Code string
}

type DeleteTeamHandler struct {
	teamRep tjudge.TeamRepository
}

func NewDeleteTeamHandler(teamRep tjudge.TeamRepository) DeleteTeamHandler {
	return DeleteTeamHandler{teamRep}
}

func (h *DeleteTeamHandler) Execute(ctx context.Context, cmd DeleteTeam) error {
	err := h.teamRep.Delete(ctx, cmd.Code)
	return err
}
