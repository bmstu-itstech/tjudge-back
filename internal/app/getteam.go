package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetTeam struct {
	TeamId string
}

type GetTeamHandler decorator.QueryHandler[GetTeam, Team]

type getTeamHandler struct {
	r tjudge.TeamRepository
}

func (h getTeamHandler) Handle(ctx context.Context, q GetTeam) (Team, error) {
	team, err := h.r.Team(ctx, tjudge.TeamId(q.TeamId))
	if err != nil {
		return Team{}, err
	}
	return teamToDto(team), nil
}

func NewGetTeamHandler(
	r tjudge.TeamRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetTeamHandler {
	return decorator.ApplyQueryDecorators(getTeamHandler{r}, l, mc)
}
